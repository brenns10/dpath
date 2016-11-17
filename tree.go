/*
tree.go contains the ParseTree interface and assorted implementations. These
are used heavily within the parser itself (dpath.y) and also do some of the heavy
lifting of evaluation.
*/

package main

import (
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
ParseTree is an interface that allows us to easily evaluate and print out code.
It has an implementation per grammar production.
*/
type ParseTree interface {
	Evaluate(ctx *Context) (Sequence, error)
	Print(to io.Writer, indent int) error
}

func getIndent(indent int) string {
	return strings.Repeat("  ", indent)
}

/*
BinopTree holds binary operator productions, including arithmetic and
comparisons.
*/
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
		return leftItem.EvalPlus(rightItem)
	case "-":
		return leftItem.EvalMinus(rightItem)
	case "*":
		return leftItem.EvalMultiply(rightItem)
	case "div":
		return leftItem.EvalDivide(rightItem)
	case "idiv":
		return leftItem.EvalIntegerDivide(rightItem)
	case "mod":
		return leftItem.EvalModulus(rightItem)
	case "to":
		return leftItem.EvalTo(rightItem)
	case "and":
		return leftItem.EvalAnd(rightItem)
	case "or":
		return leftItem.EvalOr(rightItem)
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

/*
UnopTree holds unary operator productions, which is really just plus and minus.
*/
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
	if item.TypeName() != TYPE_INTEGER && item.TypeName() != TYPE_DOUBLE {
		return nil, errors.New("unary operator expects numeric type")
	}
	if ut.Operator == "+" {
		return newSingletonSequence(item), nil
	} else if item.TypeName() == TYPE_INTEGER {
		return newSingletonSequence(newIntegerItem(-getInteger(item))), nil
	} else {
		return newSingletonSequence(newDoubleItem(-getDouble(item))), nil
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

/*
LiteralTree holds a String, Integer, Double, or literal Sequence.
*/
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
	return &LiteralTree{Type: TYPE_STRING, StringValue: parseStringLiteral(str)}
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

/*
FunccallTree holds function call production.
*/
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

	if builtin.NumArgs >= 0 && len(args) != builtin.NumArgs {
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

/*
ContextItemTree represents the use of . in an expression.
*/
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

/*
EmptySequenceTree represents an empty sequence, which is () in the language.
*/
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

/*
FilteredSequenceTree represents a predicate after an expression.
*/
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

/*
KindTree represents the special types of kind tests you can have within a step
expression (such as .. and *)
*/
type KindTree struct {
	Kind string
}

func newKindTree(s string) *KindTree {
	return &KindTree{Kind: s}
}

/*
Returns a filtered sequence of just files or directories.
Set file to true for files, or false for directories.
*/
func fileDirFilter(ctx *Context, findFiles bool) (Sequence, error) {
	seq, err := ctx.CurrentAxis.Iterate(ctx)
	if err != nil {
		return nil, err
	}
	return newConditionFilter(seq, func(it Item) bool {
		file, ok := it.(*FileItem)
		if !ok {
			return false
		}
		return findFiles != file.Info.IsDir()
	}), nil
}

func (bt *KindTree) Evaluate(ctx *Context) (Sequence, error) {
	switch bt.Kind {
	case "..":
		return ctx.Axes["parent"].Iterate(ctx)
	case "*":
		return ctx.CurrentAxis.Iterate(ctx)
	case "file":
		return fileDirFilter(ctx, true)
	case "dir":
		return fileDirFilter(ctx, false)
	default:
		return nil, errors.New("Not implemented.")
	}
}

func (t *KindTree) Print(r io.Writer, indent int) error {
	indentStr := getIndent(indent)
	_, e := io.WriteString(r, indentStr+t.Kind+"\n")
	return e
}

/*
NameTree represents any time a name occurs inside a step expression.
*/
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

/*
AxisItem represents a step expression that is prefixed by an axis, such as:
child::filename
*/
type AxisTree struct {
	Axis       string
	Expression ParseTree
}

func newAxisTree(a string, e ParseTree) *AxisTree {
	return &AxisTree{Axis: a, Expression: e}
}

func (bt *AxisTree) Evaluate(ctx *Context) (Sequence, error) {
	var err error = nil
	var ret Sequence = nil
	newAxis, ok := ctx.Axes[bt.Axis]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Axis %s not implemented.", bt.Axis))
	}
	oldAxis := ctx.CurrentAxis
	ctx.CurrentAxis = newAxis
	ret, err = bt.Expression.Evaluate(ctx)
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

/*
PathTree holds an entire path expression. Of course, single-step expressions
are not held in a PathTree (they are typically just NameTrees).
*/
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
	var Source Sequence
	pathToIterate := bt.Path

	if bt.Rooted {
		// When the path is rooted, we behave as if the path started with a step
		// expression that returned the root directory.
		rootItem, err := newFileItem("/")
		if err != nil {
			panic("Falied to set root as context item!")
		}
		Source = newSingletonSequence(rootItem)
	} else {
		// Otherwise, we use the first step expression as the source
		Source, err = bt.Path[0].Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		pathToIterate = bt.Path[1:]
	}
	// This "reduces" the path to a chain of PathSequences
	for _, pathItem := range pathToIterate {
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

/*
SequenceTree represents a sequence of comma separated expressions to evaluate and
return in sequence.
*/
type SequenceTree struct {
	Expressions []ParseTree
}

func newSequenceTree(p []ParseTree) ParseTree {
	if len(p) == 1 {
		return p[0]
	} else {
		return &SequenceTree{Expressions: p}
	}
}

func (st *SequenceTree) Evaluate(ctx *Context) (Sequence, error) {
	results := make([]Sequence, 0, len(st.Expressions))
	for _, tree := range st.Expressions {
		seq, err := tree.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		results = append(results, seq)
	}
	return newConcatenateSequence(results...), nil
}

func (st *SequenceTree) Print(w io.Writer, indent int) error {
	var e error
	indentStr := getIndent(indent)
	if _, e = io.WriteString(w, indentStr+"SEQUENCE\n"); e != nil {
		return e
	}
	for _, t := range st.Expressions {
		if e = t.Print(w, indent+1); e != nil {
			return e
		}
		if _, e = io.WriteString(w, indentStr+",\n"); e != nil {
			return e
		}
	}
	return nil
}
