package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// The following comments instruct go's build system on how to generate
// the lexer and parser.
//go:generate nex dpath.nex
//go:generate go tool yacc dpath.y

/*
A utility function that returns a single item in a sequence, raising an error
if there are zero or >1 items in the sequence.
*/
func getSingleItem(ctx *Context, s Sequence) (Item, error) {
	r, e := s.Next(ctx)
	if !r {
		return nil, errors.New("Expected one value, found none.")
	} else if e != nil {
		return nil, e
	}
	item := s.Value()
	r, e = s.Next(ctx)
	if r {
		return nil, errors.New("Too many values provided to expression.")
	} else if e != nil {
		return nil, e
	}
	return item, nil
}

/*
A utility function that "asserts" at least one item is in a sequence, panicking
if that's not the case.
*/
func panicUnlessOne(ctx *Context, s Sequence) Item {
	r, e := s.Next(ctx)
	if !r {
		panic("There wasn't a value in the sequence.")
	} else if e != nil {
		panic("Error getting value from sequence!")
	}
	return s.Value()
}

/*
Return file value, if you're certain it's a bool.
Will panic if you're wrong.
*/
func getFile(i Item) *FileItem {
	it := i.(*FileItem)
	return it
}

/*
Return bool value, if you're certain it's a bool.
Will panic if you're wrong.
*/
func getBool(i Item) bool {
	it := i.(*BooleanItem)
	return it.Value
}

/*
Return string value, if you're certain it's a string.
Will panic if you're wrong.
*/
func getString(i Item) string {
	it := i.(*StringItem)
	return it.Value
}

/*
Return integer value, if you're certain it's an integer.
Will panic if you're wrong.
*/
func getInteger(i Item) int64 {
	it := i.(*IntegerItem)
	return it.Value
}

/*
Return float value, if you're certain it's a float.
Will panic if you're wrong.
*/
func getFloat(i Item) float64 {
	it := i.(*DoubleItem)
	return it.Value
}

/*
Return numeric value as float, if you're certain it's numeric (i.e. integer
or double).
Will panic if you're wrong.
*/
func getNumericAsFloat(i Item) float64 {
	if i.Type() == TYPE_INTEGER {
		return float64(getInteger(i))
	} else {
		return getFloat(i)
	}
}

/*
Utility function for checking that each argument is one of a list of types.
The first argument is a slice of string type names. The remaining arguments are
DPath items to be type checked. Returns false if any item has type not included
in the type list. Otherwise returns true.
*/
func typeCheck(types []*Type, args ...Item) bool {
OUTER:
	for _, arg := range args {
		for _, typ := range types {
			if arg.Type() == typ {
				continue OUTER
			}
		}
		return false
	}
	return true
}

/*
ParseTree is an interface that allows us to easily print and eval code.
*/
type ParseTree interface {
	Evaluate(ctx *Context) (Sequence, error)
	Print(to io.Writer, indent int) error
}

func getIndent(indent int) string {
	return strings.Repeat("  ", indent)
}

type BinopTree struct {
	Operator string
	Left     ParseTree
	Right    ParseTree
}

func newBinopTree(op string, left ParseTree, right ParseTree) *BinopTree {
	return &BinopTree{Operator: op, Left: left, Right: right}
}

func (bt *BinopTree) Evaluate(ctx *Context) (Sequence, error) {
	var left, right Sequence
	var leftItem, rightItem Item
	var err error

	// evaluate either side
	if left, err = bt.Left.Evaluate(ctx); err != nil {
		return nil, err
	}
	if right, err = bt.Right.Evaluate(ctx); err != nil {
		return nil, err
	}

	// dispatch sequence operators
	switch bt.Operator {
	case "=":
		return GeneralComparison(ctx, left, right, CmpEq)
	case "!=":
		return GeneralComparison(ctx, left, right, CmpNe)
	case "<=":
		return GeneralComparison(ctx, left, right, CmpLe)
	case "<":
		return GeneralComparison(ctx, left, right, CmpLt)
	case ">=":
		return GeneralComparison(ctx, left, right, CmpGe)
	case ">":
		return GeneralComparison(ctx, left, right, CmpGt)
	}

	// now we get singletons from sequences
	if leftItem, err = getSingleItem(ctx, left); err != nil {
		return nil, err
	}
	if rightItem, err = getSingleItem(ctx, right); err != nil {
		return nil, err
	}

	// dispatch item operators
	switch bt.Operator {
	case "+":
		return leftItem.Type().EvalPlus(leftItem, rightItem)
	case "-":
		return leftItem.Type().EvalMinus(leftItem, rightItem)
	case "*":
		return leftItem.Type().EvalMultiply(leftItem, rightItem)
	case "div":
		return leftItem.Type().EvalDivide(leftItem, rightItem)
	case "idiv":
		return leftItem.Type().EvalIntegerDivide(leftItem, rightItem)
	case "mod":
		return leftItem.Type().EvalModulus(leftItem, rightItem)
	case "to":
		return leftItem.Type().EvalTo(leftItem, rightItem)
	case "eq":
		v, err := CmpEq(leftItem, rightItem)
		return newSingletonSequence(newBooleanItem(v)), err
	case "ne":
		v, err := CmpNe(leftItem, rightItem)
		return newSingletonSequence(newBooleanItem(v)), err
	case "le":
		v, err := CmpLe(leftItem, rightItem)
		return newSingletonSequence(newBooleanItem(v)), err
	case "lt":
		v, err := CmpLt(leftItem, rightItem)
		return newSingletonSequence(newBooleanItem(v)), err
	case "ge":
		v, err := CmpGe(leftItem, rightItem)
		return newSingletonSequence(newBooleanItem(v)), err
	case "gt":
		v, err := CmpGt(leftItem, rightItem)
		return newSingletonSequence(newBooleanItem(v)), err
	default:
		return nil, errors.New("not implemented")
	}
}

func (bt *BinopTree) Print(r io.Writer, indent int) error {
	var e error
	indentStr := getIndent(indent)
	if e = bt.Left.Print(r, indent+1); e != nil {
		return e
	}
	if _, e = io.WriteString(r, indentStr+bt.Operator+"\n"); e != nil {
		return e
	}
	if e = bt.Right.Print(r, indent+1); e != nil {
		return e
	}
	return nil
}

type UnopTree struct {
	Operator string
	Left     ParseTree
}

func newUnopTree(op string, left ParseTree) *UnopTree {
	return &UnopTree{Operator: op, Left: left}
}

func (ut *UnopTree) Evaluate(ctx *Context) (Sequence, error) {
	seq, err := ut.Left.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	item, err := getSingleItem(ctx, seq)
	if err != nil {
		return nil, err
	}
	if !typeCheck([]*Type{TYPE_INTEGER, TYPE_DOUBLE}, item) {
		return nil, errors.New("unary operator expects numeric type")
	}
	if ut.Operator == "+" {
		return newSingletonSequence(item), nil
	} else if item.Type() == TYPE_INTEGER {
		return newSingletonSequence(newIntegerItem(-getInteger(item))), nil
	} else {
		return newSingletonSequence(newDoubleItem(-getFloat(item))), nil
	}
}

func (ut *UnopTree) Print(r io.Writer, indent int) error {
	var e error
	indentStr := getIndent(indent)
	if _, e = io.WriteString(r, indentStr+ut.Operator+"\n"); e != nil {
		return e
	}
	if e = ut.Left.Print(r, indent+1); e != nil {
		return e
	}
	return nil
}

type LiteralTree struct {
	Type          *Type
	StringValue   string
	IntegerValue  int64
	DoubleValue   float64
	SequenceValue *Sequence
}

func newIntegerTree(num string) *LiteralTree {
	integer, _ := strconv.ParseInt(num, 10, 64)
	return &LiteralTree{Type: TYPE_INTEGER, IntegerValue: integer}
}

func newStringTree(str string) *LiteralTree {
	var buffer bytes.Buffer
	last := false
	sub := str[1 : len(str)-1]
	delim := rune(str[0])

	for _, char := range sub {
		if last && char == delim {
			buffer.WriteRune(char)
			last = false
		} else if !last && char == delim {
			last = true
		} else {
			buffer.WriteRune(char)
			last = false
		}
	}

	return &LiteralTree{Type: TYPE_STRING, StringValue: buffer.String()}
}

func newDoubleTree(num string) *LiteralTree {
	flt, _ := strconv.ParseFloat(num, 64)
	return &LiteralTree{Type: TYPE_DOUBLE, DoubleValue: flt}
}

func (lt *LiteralTree) Evaluate(ctx *Context) (Sequence, error) {
	switch lt.Type {
	case TYPE_STRING:
		return newSingletonSequence(newStringItem(lt.StringValue)), nil
	case TYPE_INTEGER:
		return newSingletonSequence(newIntegerItem(lt.IntegerValue)), nil
	case TYPE_DOUBLE:
		return newSingletonSequence(newDoubleItem(lt.DoubleValue)), nil
	default:
		return newEmptySequence(), nil
	}
}

func (lt *LiteralTree) Print(r io.Writer, indent int) error {
	var e error
	indentStr := getIndent(indent)
	var output string
	switch lt.Type {
	case TYPE_STRING:
		output = lt.StringValue
	case TYPE_INTEGER:
		output = strconv.FormatInt(lt.IntegerValue, 10)
	case TYPE_DOUBLE:
		output = strconv.FormatFloat(lt.DoubleValue, 'f', -1, 64)
	default:
		output = lt.Type.Name
	}
	if _, e = io.WriteString(r, indentStr+output+"\n"); e != nil {
		return e
	}
	return nil
}

type FunccallTree struct {
	Function  string
	Arguments []ParseTree
}

func newFunccallTree(name string, args []ParseTree) *FunccallTree {
	return &FunccallTree{Function: name, Arguments: args}
}

func execBuiltin(ctx *Context, name string, args ...ParseTree) (Sequence, error) {
	var err error
	builtin, ok := ctx.Namespace[name]
	if !ok {
		return nil, errors.New("builtin function " + name + " not found.")
	}

	if len(args) != builtin.NumArgs {
		return nil, errors.New(fmt.Sprintf(
			"in call to %s, expected %d args, got %d",
			name, builtin.NumArgs, len(args),
		))
	}

	arguments := make([]Sequence, len(args))
	for i, tree := range args {
		arguments[i], err = tree.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
	}
	return builtin.Invoke(ctx, arguments...)
}

func (t *FunccallTree) Evaluate(ctx *Context) (Sequence, error) {
	return execBuiltin(ctx, t.Function, t.Arguments...)
}

func (ft *FunccallTree) Print(r io.Writer, indent int) error {
	var e error
	indentStr := getIndent(indent)
	if _, e = io.WriteString(r, indentStr+ft.Function+"()\n"); e != nil {
		return e
	}
	for _, t := range ft.Arguments {
		if e = t.Print(r, indent+1); e != nil {
			return e
		}
	}
	return nil
}

type ContextItemTree struct {
}

func newContextItemTree() *ContextItemTree {
	return &ContextItemTree{}
}

func (bt *ContextItemTree) Evaluate(ctx *Context) (Sequence, error) {
	return newSingletonSequence(ctx.ContextItem), nil
}

func (t *ContextItemTree) Print(r io.Writer, indent int) error {
	indentStr := getIndent(indent)
	_, e := io.WriteString(r, indentStr+".\n")
	return e
}

type EmptySequenceTree struct {
}

func newEmptySequenceTree() *EmptySequenceTree {
	return &EmptySequenceTree{}
}

func (bt *EmptySequenceTree) Evaluate(ctx *Context) (Sequence, error) {
	return newEmptySequence(), nil
}

func (et *EmptySequenceTree) Print(r io.Writer, indent int) error {
	indentStr := getIndent(indent)
	_, e := io.WriteString(r, indentStr+"()\n")
	return e
}

type FilteredSequenceTree struct {
	Source ParseTree
	Filter []ParseTree
}

func newFilteredSequenceTree(s ParseTree, f []ParseTree) *FilteredSequenceTree {
	return &FilteredSequenceTree{Source: s, Filter: f}
}

func (bt *FilteredSequenceTree) Evaluate(ctx *Context) (Sequence, error) {
	seq, err := bt.Source.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	// BUG(stephen): numeric expressions should index into a sequence
	return newExpressionFilter(seq, bt.Filter), nil
}

func (t *FilteredSequenceTree) Print(r io.Writer, indent int) error {
	var e error
	indentStr := getIndent(indent)
	if _, e = io.WriteString(r, indentStr+"FILTER EXPRESSION:\n"); e != nil {
		return e
	}
	if e = t.Source.Print(r, indent+1); e != nil {
		return e
	}
	for _, t := range t.Filter {
		if _, e = io.WriteString(r, indentStr+"FILTER BY:\n"); e != nil {
			return e
		}
		if e = t.Print(r, indent+1); e != nil {
			return e
		}
	}
	return nil
}

type KindTree struct {
	Kind string
}

func newKindTree(s string) *KindTree {
	return &KindTree{Kind: s}
}

func (bt *KindTree) Evaluate(ctx *Context) (Sequence, error) {
	switch bt.Kind {
	case "..":
		return AXIS_PARENT.Iterate(ctx)
	case "*":
		return ctx.CurrentAxis.Iterate(ctx)
	default:
		return nil, errors.New("Not implemented.")
	}
}

func (t *KindTree) Print(r io.Writer, indent int) error {
	indentStr := getIndent(indent)
	_, e := io.WriteString(r, indentStr+t.Kind+"\n")
	return e
}

type NameTree struct {
	Name string
}

func newNameTree(s string) *NameTree {
	return &NameTree{Name: s}
}

func (bt *NameTree) Evaluate(ctx *Context) (Sequence, error) {
	return ctx.CurrentAxis.GetByName(ctx, bt.Name)
}

func (t *NameTree) Print(r io.Writer, indent int) error {
	indentStr := getIndent(indent)
	_, e := io.WriteString(r, indentStr+"Name("+t.Name+")\n")
	return e
}

type AxisTree struct {
	Axis       string
	Expression ParseTree
}

func newAxisTree(a string, e ParseTree) *AxisTree {
	return &AxisTree{Axis: a, Expression: e}
}

func (bt *AxisTree) Evaluate(ctx *Context) (Sequence, error) {
	oldAxis := ctx.CurrentAxis
	var err error = nil
	var ret Sequence = nil
	switch bt.Axis {
	case "child":
		ctx.CurrentAxis = AXIS_CHILD
	case "parent":
		ctx.CurrentAxis = AXIS_PARENT
	case "descendant":
		ctx.CurrentAxis = AXIS_DESCENDANT
	case "descendant-or-self":
		ctx.CurrentAxis = AXIS_DESCENDANT_OR_SELF
	default:
		err = errors.New(fmt.Sprintf("Axis %s not implemented.", bt.Axis))
		goto cleanup
	}
	ret, err = bt.Expression.Evaluate(ctx)
cleanup:
	ctx.CurrentAxis = oldAxis
	return ret, err
}

func (t *AxisTree) Print(r io.Writer, indent int) error {
	var e error
	indentStr := getIndent(indent)
	if _, e = io.WriteString(r, indentStr+"ON AXIS "+t.Axis+"\n"); e != nil {
		return e
	}
	return t.Expression.Print(r, indent+1)
}

type PathTree struct {
	Path   []ParseTree
	Rooted bool
}

func newPathTree(p []ParseTree, r bool) *PathTree {
	return &PathTree{Path: p, Rooted: r}
}

func (bt *PathTree) Evaluate(ctx *Context) (Sequence, error) {
	// BUG(stephen) currently this does not support the descendant axis shorthand
	// We know we have at least two items in the path due to the grammar.
	var err error = nil
	oldContext := ctx.ContextItem
	if bt.Rooted {
		ctx.ContextItem, err = newFileItem("/")
		if err != nil {
			panic("Falied to set root as context item!")
		}
	}
	Source, err := bt.Path[0].Evaluate(ctx)
	ctx.ContextItem = oldContext
	if err != nil {
		return nil, err
	}
	for _, pathItem := range bt.Path[1:] {
		if pathItem == nil {
			pathItem = newAxisTree("descendant-or-self", newKindTree("*"))
		}
		Source = newPathSequence(Source, pathItem)
	}
	return Source, nil
}

func (pt *PathTree) Print(r io.Writer, indent int) error {
	var e error
	indentStr := getIndent(indent)
	startStr := "\n"
	if pt.Rooted {
		startStr = "/\n"
	}
	if _, e = io.WriteString(r, indentStr+"PATH"+startStr); e != nil {
		return e
	}
	for _, t := range pt.Path {
		if t == nil {
			if _, e = io.WriteString(r, indentStr+"(ANY CHILD)\n"); e != nil {
				return e
			}
		} else if e = t.Print(r, indent+1); e != nil {
			return e
		}
	}
	return nil
}
