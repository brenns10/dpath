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
