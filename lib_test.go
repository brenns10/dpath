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

func TestRoundDoubles(t *testing.T) {
	cases := []string{
		"round(2.0)",
		"round(5.499999)",
		"round(5.5)",
		"round(5.9)",
		"round(-1.5)",
		"round(-1.50000001)",
	}
	results := []float64{2.0, 5.0, 6.0, 6.0, -1.0, -2.0}
	for i, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.IsType(t, (*DoubleItem)(nil), item)
		assert.Equal(t, results[i], getDouble(item))
	}
}

func TestRoundIntegers(t *testing.T) {
	cases := []string{
		"round(2)",
		"round(0)",
		"round(-2)",
	}
	results := []int64{2, 0, -2}
	for i, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.IsType(t, (*IntegerItem)(nil), item)
		assert.Equal(t, results[i], getInteger(item))
	}
}

func TestRoundInvalid(t *testing.T) {
	cases := []string{
		"round()",
		"round(1.0, 2)",
		"round('im doctor, not a string, jim')",
		"round(.)",
		"round(boolean(0))",
	}
	for _, uut := range cases {
		tree := assertParses(t, uut)
		ctx := MockDefaultContext()
		_, err := tree.Evaluate(ctx)
		assert.Error(t, err)
	}
}

func TestSubstring(t *testing.T) {
	// cases shamelessly stolen from XPath function spec:
	// https://www.w3.org/TR/xpath-functions/#func-substring
	cases := []string{
		"substring('motor car', 6)",
		"substring('metadata', 4, 3)",
		"substring('overrun', 8, 1)",
		"substring('12345', 1.5, 2.6)",
		"substring('12345', 0, 3)",
		"substring('12345', 5, -3)",
		"substring('12345', -3, 5)",
		"substring((), 1, 3)",
	}
	results := []string{
		" car",
		"ada",
		"",
		"234",
		"12",
		"",
		"1",
		"",
	}
	for i, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.IsType(t, (*StringItem)(nil), item)
		assert.Equal(t, results[i], getString(item))
	}
}

func TestSubstringInvalid(t *testing.T) {
	cases := []string{
		"substring()",
		"substring('am string with no numbers')",
		"substring(('string followed by numbers', 1, 2.3), 1)",
		"substring(1, 1)",
		"substring(boolean(1), 1)",
		"substring(1.0, 1)",
		"substring(., 1)",
		"substring('so many args', 1, 7, 9000)",
		"substring('you''re not numeric...', 'no i''m not')",
		"substring('neither is that guy', 1, 'nope')",
		"substring('sequences?', (), 1)",
		"substring('sequences?', (1, 2), 1)",
		"substring('sequences?', 1, ())",
		"substring('sequences?', 1, (1, 2))",
	}
	for _, uut := range cases {
		tree := assertParses(t, uut)
		ctx := MockDefaultContext()
		_, err := tree.Evaluate(ctx)
		assert.Error(t, err)
	}
}
