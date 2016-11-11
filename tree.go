package main

import (
	"bytes"
	//	"io"
	"strconv"
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
	//	Print(to io.Writer) error
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

type UnopTree struct {
	Operator rune
	Left     ParseTree
}

func newUnopTree(op rune, left ParseTree) *UnopTree {
	return &UnopTree{Operator: op, Left: left}
}

func (bt *UnopTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

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
	return &LiteralTree{Type: "int", DoubleValue: flt}
}

func (bt *LiteralTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

type FunccallTree struct {
	Function  string
	Arguments []ParseTree
}

func newFunccallTree(name string, args []ParseTree) *FunccallTree {
	return &FunccallTree{Function: name, Arguments: args}
}

func (bt *FunccallTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

type ContextItemTree struct {
}

func newContextItemTree() *ContextItemTree {
	return &ContextItemTree{}
}

func (bt *ContextItemTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

type EmptySequenceTree struct {
}

func newEmptySequenceTree() *EmptySequenceTree {
	return &EmptySequenceTree{}
}

func (bt *EmptySequenceTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

type FilteredSequenceTree struct {
	Source ParseTree
	Filter []ParseTree
}

func newFilteredSequenceTree(s ParseTree, f []ParseTree) *FilteredSequenceTree {
	return &FilteredSequenceTree{Source: s, Filter: f}
}

func (bt *FilteredSequenceTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

type KindTree struct {
	Kind string
}

func newKindTree(s string) *KindTree {
	return &KindTree{Kind: s}
}

func (bt *KindTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

type NameTree struct {
	Name string
}

func newNameTree(s string) *NameTree {
	return &NameTree{Name: s}
}

func (bt *NameTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

type AttrTree struct {
	Attr string
}

func newAttrTree(s string) *AttrTree {
	return &AttrTree{Attr: s}
}

func (bt *AttrTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

type AxisTree struct {
	Axis       string
	Expression ParseTree
}

func newAxisTree(a string, e ParseTree) *AxisTree {
	return &AxisTree{Axis: a, Expression: e}
}

func (bt *AxisTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }

type PathTree struct {
	Path   []ParseTree
	Rooted bool
}

func newPathTree(p []ParseTree, r bool) *PathTree {
	return &PathTree{Path: p, Rooted: r}
}

func (bt *PathTree) Evaluate(ctx Context) Sequence { return &DummySequence{} }
