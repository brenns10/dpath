package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
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
func getSingleItem(s Sequence) (Item, error) {
	if !s.Next() {
		return nil, errors.New("Expected one value, found none.")
	}
	item := s.Value()
	if s.Next() {
		return nil, errors.New("Too many values provided to expression.")
	}
	return item, nil
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
func typeCheck(types []string, args ...Item) bool {
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

/*
Utility function for binary arithmetic operators. Performs boilerplate
type checking. Takes two items and a function to call when both are integers,
as well as a function to call when one is a double.
*/
func evalArithmetic(left Item, right Item,
	bothInt func(l, r int64) Item,
	otherwise func(l, r float64) Item,
) (Sequence, error) {
	// Ensure that both arguments are numeric.
	types := []string{TYPE_INTEGER, TYPE_DOUBLE}
	if !typeCheck(types, left, right) {
		errMsg := "Expected arguments of type integer or double, got "
		errMsg += left.Type() + " and " + right.Type() + "."
		return nil, errors.New(errMsg)
	}
	// When both arguments are integers, do integer addition.
	if typeCheck([]string{TYPE_INTEGER}, left, right) {
		res := bothInt(getInteger(left), getInteger(right))
		return newSingletonSequence(res), nil
	}
	// Otherwise, up-cast to double.
	res := otherwise(getNumericAsFloat(left), getNumericAsFloat(right))
	return newSingletonSequence(res), nil
}

func evalPlus(left Item, right Item) (Sequence, error) {
	return evalArithmetic(
		left, right,
		func(l, r int64) Item { return newIntegerItem(l + r) },
		func(l, r float64) Item { return newDoubleItem(l + r) },
	)
}

func evalMinus(left, right Item) (Sequence, error) {
	return evalArithmetic(
		left, right,
		func(l, r int64) Item { return newIntegerItem(l - r) },
		func(l, r float64) Item { return newDoubleItem(l - r) },
	)
}

func evalMultiply(left, right Item) (Sequence, error) {
	return evalArithmetic(
		left, right,
		func(l, r int64) Item { return newIntegerItem(l * r) },
		func(l, r float64) Item { return newDoubleItem(l * r) },
	)
}

func evalDivide(left, right Item) (Sequence, error) {
	return evalArithmetic(
		left, right,
		func(l, r int64) Item { return newDoubleItem(float64(l) / float64(r)) },
		func(l, r float64) Item { return newDoubleItem(l / r) },
	)
}

func evalIntegerDivision(left, right Item) (Sequence, error) {
	return evalArithmetic(
		left, right,
		func(l, r int64) Item { return newIntegerItem(l / r) },
		func(l, r float64) Item { return newIntegerItem(int64(l / r)) },
	)
}

func evalModulus(left, right Item) (Sequence, error) {
	return evalArithmetic(
		left, right,
		func(l, r int64) Item { return newIntegerItem(l % r) },
		func(l, r float64) Item { return newDoubleItem(math.Mod(l, r)) },
	)
}

func (bt *BinopTree) Evaluate(ctx *Context) (Sequence, error) {
	var left, right Sequence
	var leftItem, rightItem Item
	var err error
	if left, err = bt.Left.Evaluate(ctx); err != nil {
		return nil, err
	}
	if right, err = bt.Right.Evaluate(ctx); err != nil {
		return nil, err
	}
	if leftItem, err = getSingleItem(left); err != nil {
		return nil, err
	}
	if rightItem, err = getSingleItem(right); err != nil {
		return nil, err
	}

	switch bt.Operator {
	case "+":
		return evalPlus(leftItem, rightItem)
	case "-":
		return evalMinus(leftItem, rightItem)
	case "*":
		return evalMultiply(leftItem, rightItem)
	case "div":
		return evalDivide(leftItem, rightItem)
	case "idiv":
		return evalIntegerDivision(leftItem, rightItem)
	case "mod":
		return evalModulus(leftItem, rightItem)
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
	return &DummySequence{}, nil
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
	Type          string
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
		output = lt.Type
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

func (t *FunccallTree) Evaluate(ctx *Context) (Sequence, error) {
	var err error
	builtin, ok := ctx.Namespace[t.Function]
	if !ok {
		return nil, errors.New("builtin function " + t.Function + " not found.")
	}

	if len(t.Arguments) != builtin.NumArgs {
		return nil, errors.New(fmt.Sprintf(
			"in call to %s, expected %d args, got %d",
			t.Function, builtin.NumArgs, len(t.Arguments),
		))
	}

	arguments := make([]Sequence, len(t.Arguments))
	for i, tree := range t.Arguments {
		arguments[i], err = tree.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
	}
	return builtin.Invoke(ctx, arguments...)
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

func (bt *ContextItemTree) Evaluate(ctx *Context) (Sequence, error) { return &DummySequence{}, nil }

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

func (bt *FilteredSequenceTree) Evaluate(ctx *Context) (Sequence, error) { return &DummySequence{}, nil }

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

func (bt *KindTree) Evaluate(ctx *Context) (Sequence, error) { return &DummySequence{}, nil }

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

func (bt *NameTree) Evaluate(ctx *Context) (Sequence, error) { return &DummySequence{}, nil }

func (t *NameTree) Print(r io.Writer, indent int) error {
	indentStr := getIndent(indent)
	_, e := io.WriteString(r, indentStr+"Name("+t.Name+")\n")
	return e
}

type AttrTree struct {
	Attr string
}

func newAttrTree(s string) *AttrTree {
	return &AttrTree{Attr: s}
}

func (bt *AttrTree) Evaluate(ctx *Context) (Sequence, error) { return &DummySequence{}, nil }

func (t *AttrTree) Print(r io.Writer, indent int) error {
	indentStr := getIndent(indent)
	_, e := io.WriteString(r, indentStr+"Attr("+t.Attr+")\n")
	return e
}

type AxisTree struct {
	Axis       string
	Expression ParseTree
}

func newAxisTree(a string, e ParseTree) *AxisTree {
	return &AxisTree{Axis: a, Expression: e}
}

func (bt *AxisTree) Evaluate(ctx *Context) (Sequence, error) { return &DummySequence{}, nil }

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

func (bt *PathTree) Evaluate(ctx *Context) (Sequence, error) { return &DummySequence{}, nil }

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
