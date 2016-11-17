package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleExpressionsParse(t *testing.T) {
	assertLiteral(t, "1")
	assertLiteral(t, "1.23")
	assertLiteral(t, "1.0E-1")
	assertLiteral(t, "'yo'")
	assertQName(t, "identifier")
	assertEmptySequenceTree(t, "()")
}

func TestRangeExpressions(t *testing.T) {
	bt := assertBinop(t, "1 + 1 to 2 + 2")
	assert.Equal(t, bt.Operator, "to")

	assert.IsType(t, (*BinopTree)(nil), bt.Left)
	left := bt.Left.(*BinopTree)
	assert.IsType(t, (*LiteralTree)(nil), left.Left)
	assert.IsType(t, (*LiteralTree)(nil), left.Right)

	assert.IsType(t, (*BinopTree)(nil), bt.Right)
	right := bt.Left.(*BinopTree)
	assert.IsType(t, (*LiteralTree)(nil), right.Left)
	assert.IsType(t, (*LiteralTree)(nil), right.Right)
}

func TestLogicExpressions(t *testing.T) {
	bt := assertBinop(t, "x = 5 and 3 + 3 eq 6")
	assert.Equal(t, bt.Operator, "and")

	assert.IsType(t, (*BinopTree)(nil), bt.Left)
	l := bt.Left.(*BinopTree)
	assert.Equal(t, "=", l.Operator)
	assert.IsType(t, (*NameTree)(nil), l.Left)
	assert.IsType(t, (*LiteralTree)(nil), l.Right)

	assert.IsType(t, (*BinopTree)(nil), bt.Right)
	r := bt.Right.(*BinopTree)
	assert.Equal(t, "eq", r.Operator)
	assert.IsType(t, (*LiteralTree)(nil), r.Right)
	assert.IsType(t, (*BinopTree)(nil), r.Left)
	rl := r.Left.(*BinopTree)
	assert.Equal(t, "+", rl.Operator)
	assert.IsType(t, (*LiteralTree)(nil), rl.Left)
	assert.IsType(t, (*LiteralTree)(nil), rl.Right)

}

func TestComparisons(t *testing.T) {
	assertParses(t, "1+1 eq 1")
	assertParses(t, "1 = 1")
	assertParses(t, "1 ne 1* /*")
	assertParses(t, "1 != 1")
	assertParses(t, "1 gt 1")
	assertParses(t, ". > /name[@blah]")
	assertParses(t, "1 ge 1")
	assertParses(t, "1 >= 1")
	assertParses(t, "1 lt 1")
	assertParses(t, "1 < 1")
	assertParses(t, "1 le 1")
	assertParses(t, "-1 <= 1")
}

func TestSimplePath(t *testing.T) {
	root := assertParses(t, "/simple/path/here")
	assert.IsType(t, (*PathTree)(nil), root)
	pt := root.(*PathTree)
	assert.True(t, pt.Rooted)
	assert.Len(t, pt.Path, 3)
	assert.IsType(t, (*NameTree)(nil), pt.Path[0])
	assert.IsType(t, (*NameTree)(nil), pt.Path[1])
	assert.IsType(t, (*NameTree)(nil), pt.Path[2])
}

func TestAnyChildPath(t *testing.T) {
	root := assertParses(t, "simple//path")
	assert.IsType(t, (*PathTree)(nil), root)
	pt := root.(*PathTree)
	assert.False(t, pt.Rooted)
	assert.Len(t, pt.Path, 3)
	assert.IsType(t, (*NameTree)(nil), pt.Path[0])
	assert.Nil(t, pt.Path[1])
	assert.IsType(t, (*NameTree)(nil), pt.Path[2])
}

func TestFilteredPath(t *testing.T) {
	root := assertParses(t, "simple[1 = 1][3 <= 3]/path")
	assert.IsType(t, (*PathTree)(nil), root)
	pt := root.(*PathTree)
	assert.False(t, pt.Rooted)
	assert.Len(t, pt.Path, 2)
	assert.IsType(t, (*FilteredSequenceTree)(nil), pt.Path[0])
	fst := pt.Path[0].(*FilteredSequenceTree)
	assert.IsType(t, (*NameTree)(nil), fst.Source)
	assert.Len(t, fst.Filter, 2)
	assert.IsType(t, (*BinopTree)(nil), fst.Filter[0])
	filter1 := fst.Filter[0].(*BinopTree)
	assert.Equal(t, "=", filter1.Operator)
	assert.IsType(t, (*BinopTree)(nil), fst.Filter[0])
	filter2 := fst.Filter[1].(*BinopTree)
	assert.Equal(t, "<=", filter2.Operator)
}

func TestLiteralName(t *testing.T) {
	uut := "#'My very long file name.docx'"
	root := assertParses(t, uut)
	assert.IsType(t, (*NameTree)(nil), root)
	pt := root.(*NameTree)
	assert.Equal(t, pt.Name, "My very long file name.docx")
}
