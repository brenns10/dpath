package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertEvaluates(t *testing.T, s string) (Sequence, *Context) {
	tree, e := ParseString(s)
	assert.Nil(t, e)
	ctx := DefaultContext()
	res, e := tree.Evaluate(ctx)
	assert.Nil(t, e)
	return res, ctx
}

func seqToSlice(seq Sequence, ctx *Context) ([]Item, error) {
	var err error
	var next bool
	items := make([]Item, 0, 5)
	for next, err = seq.Next(ctx); next && err == nil; next, err = seq.Next(ctx) {
		items = append(items, seq.Value())
	}
	return items, err
}

func assertSingleton(t *testing.T, ctx *Context, seq Sequence) Item {
	hasNext, err := seq.Next(ctx)
	assert.Nil(t, err)
	assert.True(t, hasNext)
	item := seq.Value()
	hasNext, err = seq.Next(ctx)
	assert.Nil(t, err)
	assert.False(t, hasNext)
	return item
}

func assertEmptySequence(t *testing.T, ctx *Context, seq Sequence) {
	hasNext, err := seq.Next(ctx)
	assert.Nil(t, err)
	assert.False(t, hasNext)
}

func TestIntegerLiteral(t *testing.T) {
	seq, ctx := assertEvaluates(t, "1989")
	item := assertSingleton(t, ctx, seq)
	assert.Equal(t, item.TypeName(), TYPE_INTEGER)
	intItem := item.(*IntegerItem)
	assert.Equal(t, int64(1989), intItem.Value)
}

func TestDecimalLiteral(t *testing.T) {
	seq, ctx := assertEvaluates(t, "1.234")
	item := assertSingleton(t, ctx, seq)
	assert.Equal(t, item.TypeName(), TYPE_DOUBLE)
	doubleItem := item.(*DoubleItem)
	assert.Equal(t, 1.234, doubleItem.Value)
}

func TestFloatLiteral(t *testing.T) {
	seq, ctx := assertEvaluates(t, "1.0e-1")
	item := assertSingleton(t, ctx, seq)
	assert.Equal(t, item.TypeName(), TYPE_DOUBLE)
	doubleItem := item.(*DoubleItem)
	assert.Equal(t, 1.0e-1, doubleItem.Value)
}

func TestStringLiteral(t *testing.T) {
	seq, ctx := assertEvaluates(t, "'foo'")
	item := assertSingleton(t, ctx, seq)
	assert.Equal(t, item.TypeName(), TYPE_STRING)
	stringItem := item.(*StringItem)
	assert.Equal(t, "foo", stringItem.Value)
}

func TestStringEscapes(t *testing.T) {
	seq, ctx := assertEvaluates(t, "\"bar\"\"\"")
	item := assertSingleton(t, ctx, seq)
	assert.Equal(t, item.TypeName(), TYPE_STRING)
	stringItem := item.(*StringItem)
	assert.Equal(t, "bar\"", stringItem.Value)
}

func TestEmptySequence(t *testing.T) {
	seq, ctx := assertEvaluates(t, "()")
	assertEmptySequence(t, ctx, seq)
}

func TestAdditionIntegers(t *testing.T) {
	seq, ctx := assertEvaluates(t, "1 + 1")
	item := assertSingleton(t, ctx, seq)
	assert.Equal(t, item.TypeName(), TYPE_INTEGER)
	intItem := item.(*IntegerItem)
	assert.Equal(t, int64(2), intItem.Value)
}

func TestAdditionDoubles(t *testing.T) {
	cases := []string{"1.0 + 1", "1 + 1.0", "1.0 + 1.0"}
	for _, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.Equal(t, item.TypeName(), TYPE_DOUBLE)
		doubleItem := item.(*DoubleItem)
		assert.Equal(t, float64(2.0), doubleItem.Value)
	}
}

func TestSubtractionIntegers(t *testing.T) {
	seq, ctx := assertEvaluates(t, "2 - 1")
	item := assertSingleton(t, ctx, seq)
	assert.Equal(t, item.TypeName(), TYPE_INTEGER)
	intItem := item.(*IntegerItem)
	assert.Equal(t, int64(1), intItem.Value)
}

func TestSubtractionDoubles(t *testing.T) {
	cases := []string{"2.0 - 1", "2 - 1.0", "2.0 - 1.0"}
	for _, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.Equal(t, item.TypeName(), TYPE_DOUBLE)
		doubleItem := item.(*DoubleItem)
		assert.Equal(t, float64(1.0), doubleItem.Value)
	}
}

func TestMultiplicationIntegers(t *testing.T) {
	seq, ctx := assertEvaluates(t, "5 * 3")
	item := assertSingleton(t, ctx, seq)
	assert.Equal(t, item.TypeName(), TYPE_INTEGER)
	intItem := item.(*IntegerItem)
	assert.Equal(t, int64(15), intItem.Value)
}

func TestMultiplactionDoubles(t *testing.T) {
	cases := []string{"5.0 * 3", "5 * 3.0", "5.0 * 3.0"}
	for _, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.Equal(t, item.TypeName(), TYPE_DOUBLE)
		doubleItem := item.(*DoubleItem)
		assert.Equal(t, float64(15.0), doubleItem.Value)
	}
}

func TestDivision(t *testing.T) {
	cases := []string{"5 div 2", "5.0 div 2", "5 div 2.0", "5.0 div 2.0"}
	for _, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.Equal(t, item.TypeName(), TYPE_DOUBLE)
		doubleItem := item.(*DoubleItem)
		assert.Equal(t, float64(2.5), doubleItem.Value)
	}
}

func TestIntegerDivision(t *testing.T) {
	cases := []string{"5 idiv 2", "5.0 idiv 2", "5 idiv 2.0", "5.0 idiv 2.0"}
	for _, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.Equal(t, item.TypeName(), TYPE_INTEGER)
		intItem := item.(*IntegerItem)
		assert.Equal(t, int64(2), intItem.Value)
	}
}

func TestModulusInteger(t *testing.T) {
	cases := []string{"5 mod 2"}
	for _, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.Equal(t, item.TypeName(), TYPE_INTEGER)
		intItem := item.(*IntegerItem)
		assert.Equal(t, int64(1), intItem.Value)
	}
}

func TestModulusDouble(t *testing.T) {
	cases := []string{"5.0 mod 2", "5 mod 2.0", "5.0 mod 2.0"}
	for _, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.Equal(t, item.TypeName(), TYPE_DOUBLE)
		doubleItem := item.(*DoubleItem)
		assert.Equal(t, float64(1.0), doubleItem.Value)
	}
}

func TestBinopIncorrectTypesFail(t *testing.T) {
	cases := []string{
		"1 + 'foo'",
		"1 - ()",
		"'str' div 7.3",
		"'blah' * 2",
		"'hello' idiv 7",
		"'bye' mod 3",
	}
	for _, uut := range cases {
		tree := assertParses(t, uut)
		ctx := DefaultContext()
		_, err := tree.Evaluate(ctx)
		assert.Error(t, err)
	}
}

func TestUnopPlusInteger(t *testing.T) {
	seq, ctx := assertEvaluates(t, "+5")
	item := assertSingleton(t, ctx, seq)
	assert.Equal(t, item.TypeName(), TYPE_INTEGER)
	intItem := item.(*IntegerItem)
	assert.Equal(t, int64(5), intItem.Value)
}

func TestUnopPlusDouble(t *testing.T) {
	seq, ctx := assertEvaluates(t, "+5.0")
	item := assertSingleton(t, ctx, seq)
	assert.Equal(t, item.TypeName(), TYPE_DOUBLE)
	doubleItem := item.(*DoubleItem)
	assert.Equal(t, float64(5.0), doubleItem.Value)
}

func TestUnopMinusInteger(t *testing.T) {
	seq, ctx := assertEvaluates(t, "-5")
	item := assertSingleton(t, ctx, seq)
	assert.Equal(t, item.TypeName(), TYPE_INTEGER)
	intItem := item.(*IntegerItem)
	assert.Equal(t, int64(-5), intItem.Value)
}

func TestUnopMinusDouble(t *testing.T) {
	seq, ctx := assertEvaluates(t, "-5.0")
	item := assertSingleton(t, ctx, seq)
	assert.Equal(t, item.TypeName(), TYPE_DOUBLE)
	doubleItem := item.(*DoubleItem)
	assert.Equal(t, float64(-5.0), doubleItem.Value)
}

func TestUnopIncorrectTypesFail(t *testing.T) {
	cases := []string{
		"+'foo'",
		"- ()",
	}
	for _, uut := range cases {
		tree := assertParses(t, uut)
		ctx := DefaultContext()
		_, err := tree.Evaluate(ctx)
		assert.Error(t, err)
	}
}

func TestBooleanOperators(t *testing.T) {
	cases := []string{
		"boolean(0) and boolean(0)",
		"boolean(0) and boolean(1)",
		"boolean(1) and boolean(0)",
		"boolean(1) and boolean(1)",
		"boolean(0) or boolean(0)",
		"boolean(0) or boolean(1)",
		"boolean(1) or boolean(0)",
		"boolean(1) or boolean(1)",
	}
	results := []bool{false, false, false, true, false, true, true, true}
	for i, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.Equal(t, getBool(item), results[i])
	}
}

func TestBooleanOperatorsIncorrectTypes(t *testing.T) {
	cases := []string{
		"1 or boolean(0)",
		"2.0 and boolean(0)",
		"(1, 2, 3) or boolean(0)",
		"(boolean(0), 1) and boolean(0)",
	}
	for _, uut := range cases {
		tree := assertParses(t, uut)
		ctx := DefaultContext()
		_, err := tree.Evaluate(ctx)
		assert.Error(t, err)
	}
}

func TestCommaOperator(t *testing.T) {
	uut := "1 + 1, boolean(0), 'hello', 3.14159, (5, 6)"
	seq, ctx := assertEvaluates(t, uut)
	items, err := seqToSlice(seq, ctx)
	assert.Nil(t, err)
	assert.Len(t, items, 6)
	assert.Equal(t, int64(2), getInteger(items[0]))
	assert.Equal(t, false, getBool(items[1]))
	assert.Equal(t, "hello", getString(items[2]))
	assert.Equal(t, float64(3.14159), getDouble(items[3]))
	assert.Equal(t, int64(5), getInteger(items[4]))
	assert.Equal(t, int64(6), getInteger(items[5]))
}
