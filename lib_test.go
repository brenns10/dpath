package main

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestBooleanEmptySequence(t *testing.T) {
	ctx := MockDefaultContext()
	s, e1 := BuiltinBooleanInvoke(ctx, newEmptySequence())
	assert.Nil(t, e1)
	i, e2 := getSingleItem(ctx, s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.False(t, v.Value)
}

func TestBooleanEmptyString(t *testing.T) {
	ctx := MockDefaultContext()
	input := newSingletonSequence(newStringItem(""))

	s, e1 := BuiltinBooleanInvoke(ctx, input)
	assert.Nil(t, e1)
	i, e2 := getSingleItem(ctx, s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.False(t, v.Value)
}

func TestBooleanNonEmptyString(t *testing.T) {
	ctx := MockDefaultContext()
	input := newSingletonSequence(newStringItem("foo"))

	s, e1 := BuiltinBooleanInvoke(ctx, input)
	assert.Nil(t, e1)
	i, e2 := getSingleItem(ctx, s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.True(t, v.Value)
}

func TestBooleanIntegerZero(t *testing.T) {
	ctx := MockDefaultContext()
	input := newSingletonSequence(newIntegerItem(int64(0)))

	s, e1 := BuiltinBooleanInvoke(ctx, input)
	assert.Nil(t, e1)
	i, e2 := getSingleItem(ctx, s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.False(t, v.Value)
}

func TestBooleanIntegerNonzero(t *testing.T) {
	ctx := MockDefaultContext()
	input := newSingletonSequence(newIntegerItem(int64(1)))

	s, e1 := BuiltinBooleanInvoke(ctx, input)
	assert.Nil(t, e1)
	i, e2 := getSingleItem(ctx, s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.True(t, v.Value)
}

func TestBooleanDoubleZero(t *testing.T) {
	ctx := MockDefaultContext()
	input := newSingletonSequence(newDoubleItem(float64(0.0)))

	s, e1 := BuiltinBooleanInvoke(ctx, input)
	assert.Nil(t, e1)
	i, e2 := getSingleItem(ctx, s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.False(t, v.Value)
}

func TestBooleanDoubleNan(t *testing.T) {
	ctx := MockDefaultContext()
	input := newSingletonSequence(newDoubleItem(math.NaN()))

	s, e1 := BuiltinBooleanInvoke(ctx, input)
	assert.Nil(t, e1)
	i, e2 := getSingleItem(ctx, s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.False(t, v.Value)
}

func TestBooleanDoubleTrue(t *testing.T) {
	ctx := MockDefaultContext()
	input := newSingletonSequence(newDoubleItem(float64(1.1)))

	s, e1 := BuiltinBooleanInvoke(ctx, input)
	assert.Nil(t, e1)
	i, e2 := getSingleItem(ctx, s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.True(t, v.Value)
}

func TestBooleanFalseBool(t *testing.T) {
	ctx := MockDefaultContext()
	input := newSingletonSequence(newBooleanItem(false))

	s, e1 := BuiltinBooleanInvoke(ctx, input)
	assert.Nil(t, e1)
	i, e2 := getSingleItem(ctx, s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.False(t, v.Value)
}

func TestBooleanTrueBool(t *testing.T) {
	ctx := MockDefaultContext()
	input := newSingletonSequence(newBooleanItem(true))

	s, e1 := BuiltinBooleanInvoke(ctx, input)
	assert.Nil(t, e1)
	i, e2 := getSingleItem(ctx, s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.True(t, v.Value)
}

func TestConcat(t *testing.T) {
	ctx := MockDefaultContext()
	cases := []string{
		"concat(1)",
		"concat(1.2)",
		"concat(boolean(1))",
		"concat(boolean(0))",
		"concat(.)", // this will give us the mocked current file
		"concat('iAmString')",
		"concat('s1', 's2')",
		"concat('s1', 3)",
		"concat('s1', 3.14)",
	}
	results := []string{
		"1", "1.2", "true", "false", getFile(ctx.ContextItem).Info.Name(),
		"iAmString", "s1s2", "s13", "s13.14",
	}
	for i, uut := range cases {
		seq := assertEvaluatesCtx(t, uut, ctx)
		item := assertSingleton(t, ctx, seq)
		assert.IsType(t, (*StringItem)(nil), item)
		strItem := item.(*StringItem)
		assert.Equal(t, results[i], strItem.Value)
	}
}

func TestConcatNoArgsFails(t *testing.T) {
	uut := "concat()"
	tree, err := ParseString(uut)
	assert.Nil(t, err)
	ctx := MockDefaultContext()
	_, err = tree.Evaluate(ctx)
	assert.Error(t, err)
}
