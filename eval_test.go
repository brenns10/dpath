package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertEvaluates(t *testing.T, s string) Sequence {
	tree, e := ParseString(s)
	assert.Nil(t, e)
	ctx := DefaultContext()
	res, e := tree.Evaluate(ctx)
	assert.Nil(t, e)
	return res
}

func assertSingleton(t *testing.T, seq Sequence) Item {
	assert.True(t, seq.Next())
	item := seq.Value()
	assert.False(t, seq.Next())
	return item
}

func assertEmptySequence(t *testing.T, seq Sequence) {
	_, ok := seq.(*DummySequence)
	assert.False(t, ok)
	assert.False(t, seq.Next())
}

func TestIntegerLiteral(t *testing.T) {
	seq := assertEvaluates(t, "1989")
	item := assertSingleton(t, seq)
	assert.Equal(t, item.Type(), TYPE_INTEGER)
	intItem := item.(*IntegerItem)
	assert.Equal(t, int64(1989), intItem.Value)
}

func TestDecimalLiteral(t *testing.T) {
	seq := assertEvaluates(t, "1.234")
	item := assertSingleton(t, seq)
	assert.Equal(t, item.Type(), TYPE_DOUBLE)
	doubleItem := item.(*DoubleItem)
	assert.Equal(t, 1.234, doubleItem.Value)
}

func TestFloatLiteral(t *testing.T) {
	seq := assertEvaluates(t, "1.0e-1")
	item := assertSingleton(t, seq)
	assert.Equal(t, item.Type(), TYPE_DOUBLE)
	doubleItem := item.(*DoubleItem)
	assert.Equal(t, 1.0e-1, doubleItem.Value)
}

func TestStringLiteral(t *testing.T) {
	seq := assertEvaluates(t, "'foo'")
	item := assertSingleton(t, seq)
	assert.Equal(t, item.Type(), TYPE_STRING)
	stringItem := item.(*StringItem)
	assert.Equal(t, "foo", stringItem.Value)
}

func TestStringEscapes(t *testing.T) {
	seq := assertEvaluates(t, "\"bar\"\"\"")
	item := assertSingleton(t, seq)
	assert.Equal(t, item.Type(), TYPE_STRING)
	stringItem := item.(*StringItem)
	assert.Equal(t, "bar\"", stringItem.Value)
}

func TestEmptySequence(t *testing.T) {
	seq := assertEvaluates(t, "()")
	assertEmptySequence(t, seq)
}

func TestAdditionIntegers(t *testing.T) {
	seq := assertEvaluates(t, "1 + 1")
	item := assertSingleton(t, seq)
	assert.Equal(t, item.Type(), TYPE_INTEGER)
	intItem := item.(*IntegerItem)
	assert.Equal(t, int64(2), intItem.Value)
}

func TestAdditionDoubles(t *testing.T) {
	cases := []string{"1.0 + 1", "1 + 1.0", "1.0 + 1.0"}
	for _, uut := range cases {
		seq := assertEvaluates(t, uut)
		item := assertSingleton(t, seq)
		assert.Equal(t, item.Type(), TYPE_DOUBLE)
		doubleItem := item.(*DoubleItem)
		assert.Equal(t, float64(2.0), doubleItem.Value)
	}
}

func TestSubtractionIntegers(t *testing.T) {
	seq := assertEvaluates(t, "2 - 1")
	item := assertSingleton(t, seq)
	assert.Equal(t, item.Type(), TYPE_INTEGER)
	intItem := item.(*IntegerItem)
	assert.Equal(t, int64(1), intItem.Value)
}

func TestSubtractionDoubles(t *testing.T) {
	cases := []string{"2.0 - 1", "2 - 1.0", "2.0 - 1.0"}
	for _, uut := range cases {
		seq := assertEvaluates(t, uut)
		item := assertSingleton(t, seq)
		assert.Equal(t, item.Type(), TYPE_DOUBLE)
		doubleItem := item.(*DoubleItem)
		assert.Equal(t, float64(1.0), doubleItem.Value)
	}
}

func TestMultiplicationIntegers(t *testing.T) {
	seq := assertEvaluates(t, "5 * 3")
	item := assertSingleton(t, seq)
	assert.Equal(t, item.Type(), TYPE_INTEGER)
	intItem := item.(*IntegerItem)
	assert.Equal(t, int64(15), intItem.Value)
}

func TestMultiplactionDoubles(t *testing.T) {
	cases := []string{"5.0 * 3", "5 * 3.0", "5.0 * 3.0"}
	for _, uut := range cases {
		seq := assertEvaluates(t, uut)
		item := assertSingleton(t, seq)
		assert.Equal(t, item.Type(), TYPE_DOUBLE)
		doubleItem := item.(*DoubleItem)
		assert.Equal(t, float64(15.0), doubleItem.Value)
	}
}

func TestDivision(t *testing.T) {
	cases := []string{"5 div 2", "5.0 div 2", "5 div 2.0", "5.0 div 2.0"}
	for _, uut := range cases {
		seq := assertEvaluates(t, uut)
		item := assertSingleton(t, seq)
		assert.Equal(t, item.Type(), TYPE_DOUBLE)
		doubleItem := item.(*DoubleItem)
		assert.Equal(t, float64(2.5), doubleItem.Value)
	}
}

func TestIntegerDivision(t *testing.T) {
	cases := []string{"5 idiv 2", "5.0 idiv 2", "5 idiv 2.0", "5.0 idiv 2.0"}
	for _, uut := range cases {
		seq := assertEvaluates(t, uut)
		item := assertSingleton(t, seq)
		assert.Equal(t, item.Type(), TYPE_INTEGER)
		intItem := item.(*IntegerItem)
		assert.Equal(t, int64(2), intItem.Value)
	}
}

func TestModulusInteger(t *testing.T) {
	cases := []string{"5 mod 2"}
	for _, uut := range cases {
		seq := assertEvaluates(t, uut)
		item := assertSingleton(t, seq)
		assert.Equal(t, item.Type(), TYPE_INTEGER)
		intItem := item.(*IntegerItem)
		assert.Equal(t, int64(1), intItem.Value)
	}
}

func TestModulusDouble(t *testing.T) {
	cases := []string{"5.0 mod 2", "5 mod 2.0", "5.0 mod 2.0"}
	for _, uut := range cases {
		seq := assertEvaluates(t, uut)
		item := assertSingleton(t, seq)
		assert.Equal(t, item.Type(), TYPE_DOUBLE)
		doubleItem := item.(*DoubleItem)
		assert.Equal(t, float64(1.0), doubleItem.Value)
	}
}

func TestIncorrectTypesFail(t *testing.T) {
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
