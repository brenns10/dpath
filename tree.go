package main

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
)

// The following comments instruct go's build system on how to generate
// the lexer and parser.
//go:generate nex dpath.nex
//go:generate go tool yacc dpath.y

type Context struct {
}
type Item interface {
	Type() string
}
type DummyItem struct{}

func (d *DummyItem) Type() string {
	return "dummy"
}

type Sequence interface {
	Next() Item
}
type DummySequence struct{}

func (d *DummySequence) Next() Item {
	return &DummyItem{}
}

type ParseTree interface {
	Evaluate(ctx Context) Sequence
	Print(to *bufio.Writer, indent int) error
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

func (bt *BinopTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

func (bt *BinopTree) Print(r *bufio.Writer, indent int) error {
	var e error
	indentStr := getIndent(indent)
	if e = bt.Left.Print(r, indent+1); e != nil {
		return e
	}
	if _, e = r.WriteString(indentStr + bt.Operator + "\n"); e != nil {
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

func (ut *UnopTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

func (ut *UnopTree) Print(r *bufio.Writer, indent int) error {
	var e error
	indentStr := getIndent(indent)
	if _, e = r.WriteString(indentStr + ut.Operator + "\n"); e != nil {
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
	return &LiteralTree{Type: "int", IntegerValue: integer}
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

	return &LiteralTree{Type: "string", StringValue: buffer.String()}
}

func newDoubleTree(num string) *LiteralTree {
	flt, _ := strconv.ParseFloat(num, 64)
	return &LiteralTree{Type: "double", DoubleValue: flt}
}

func (bt *LiteralTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

func (lt *LiteralTree) Print(r *bufio.Writer, indent int) error {
	var e error
	indentStr := getIndent(indent)
	var output string
	switch lt.Type {
	case "string":
		output = lt.StringValue
	case "int":
		output = strconv.FormatInt(lt.IntegerValue, 10)
	case "double":
		output = strconv.FormatFloat(lt.DoubleValue, 'f', -1, 64)
	default:
		output = lt.Type
	}
	if _, e = r.WriteString(indentStr + output + "\n"); e != nil {
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

func (bt *FunccallTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

func (ft *FunccallTree) Print(r *bufio.Writer, indent int) error {
	var e error
	indentStr := getIndent(indent)
	if _, e = r.WriteString(indentStr + ft.Function + "()\n"); e != nil {
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

func (bt *ContextItemTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

func (t *ContextItemTree) Print(r *bufio.Writer, indent int) error {
	indentStr := getIndent(indent)
	_, e := r.WriteString(indentStr + ".\n")
	return e
}

type EmptySequenceTree struct {
}

func newEmptySequenceTree() *EmptySequenceTree {
	return &EmptySequenceTree{}
}

func (bt *EmptySequenceTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

func (et *EmptySequenceTree) Print(r *bufio.Writer, indent int) error {
	indentStr := getIndent(indent)
	_, e := r.WriteString(indentStr + "()\n")
	return e
}

type FilteredSequenceTree struct {
	Source ParseTree
	Filter []ParseTree
}

func newFilteredSequenceTree(s ParseTree, f []ParseTree) *FilteredSequenceTree {
	return &FilteredSequenceTree{Source: s, Filter: f}
}

func (bt *FilteredSequenceTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

func (t *FilteredSequenceTree) Print(r *bufio.Writer, indent int) error {
	var e error
	indentStr := getIndent(indent)
	if _, e = r.WriteString(indentStr + "FILTER EXPRESSION:\n"); e != nil {
		return e
	}
	if e = t.Source.Print(r, indent+1); e != nil {
		return e
	}
	for _, t := range t.Filter {
		if _, e = r.WriteString(indentStr + "FILTER BY:\n"); e != nil {
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

func (bt *KindTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

func (t *KindTree) Print(r *bufio.Writer, indent int) error {
	indentStr := getIndent(indent)
	_, e := r.WriteString(indentStr + t.Kind + "\n")
	return e
}

type NameTree struct {
	Name string
}

func newNameTree(s string) *NameTree {
	return &NameTree{Name: s}
}

func (bt *NameTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

func (t *NameTree) Print(r *bufio.Writer, indent int) error {
	indentStr := getIndent(indent)
	_, e := r.WriteString(indentStr + "Name(" + t.Name + ")\n")
	return e
}

type AttrTree struct {
	Attr string
}

func newAttrTree(s string) *AttrTree {
	return &AttrTree{Attr: s}
}

func (bt *AttrTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

func (t *AttrTree) Print(r *bufio.Writer, indent int) error {
	indentStr := getIndent(indent)
	_, e := r.WriteString(indentStr + "Attr(" + t.Attr + ")\n")
	return e
}

type AxisTree struct {
	Axis       string
	Expression ParseTree
}

func newAxisTree(a string, e ParseTree) *AxisTree {
	return &AxisTree{Axis: a, Expression: e}
}

func (bt *AxisTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

func (t *AxisTree) Print(r *bufio.Writer, indent int) error {
	var e error
	indentStr := getIndent(indent)
	if _, e = r.WriteString(indentStr + "ON AXIS " + t.Axis + "\n"); e != nil {
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

func (bt *PathTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

func (pt *PathTree) Print(r *bufio.Writer, indent int) error {
	var e error
	indentStr := getIndent(indent)
	startStr := "\n"
	if pt.Rooted {
		startStr = "/\n"
	}
	if _, e = r.WriteString(indentStr + "PATH" + startStr); e != nil {
		return e
	}
	for _, t := range pt.Path {
		if t == nil {
			if _, e = r.WriteString(indentStr + "(ANY CHILD)\n"); e != nil {
				return e
			}
		} else if e = t.Print(r, indent+1); e != nil {
			return e
		}
	}
	return nil
}
