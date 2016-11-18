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
		assert.IsType(t, (*StringItem)(nil), item, uut)
		strItem := item.(*StringItem)
		assert.Equal(t, results[i], strItem.Value, uut)
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
		assert.Error(t, err, uut)
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
		assert.IsType(t, (*StringItem)(nil), item, uut)
		assert.Equal(t, results[i], getString(item), uut)
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
		assert.Error(t, err, uut)
	}
}

func TestStringBuiltin(t *testing.T) {
	cases := []string{
		"string(1)",
		"string(1.1)",
		"string(boolean(0))",
		"string('hi there')",
		"string()",
	}
	ctx := MockDefaultContext()
	results := []string{
		"1", "1.1", "false", "hi there", ctx.ContextItem.ToString(),
	}
	for i, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.IsType(t, (*StringItem)(nil), item, uut)
		assert.Equal(t, results[i], getString(item), uut)
	}
}

func TestStringLengthBuiltin(t *testing.T) {
	cases := []string{
		"string-length(())",
		"string-length('hi there')",
		"string-length()",
	}
	ctx := MockDefaultContext()
	results := []int64{
		0, 8, int64(len(ctx.ContextItem.ToString())),
	}
	for i, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.IsType(t, (*IntegerItem)(nil), item, uut)
		assert.Equal(t, results[i], getInteger(item), uut)
	}
}

func TestStringLengthInvalid(t *testing.T) {
	cases := []string{
		"string-length(1)",
		"string-length(1.1)",
	}
	for _, uut := range cases {
		tree := assertParses(t, uut)
		ctx := MockDefaultContext()
		_, err := tree.Evaluate(ctx)
		assert.Error(t, err, uut)
	}
}

func TestEndsWith(t *testing.T) {
	cases := []string{
		"ends-with('abcdef', 'def')",
		"ends-with('fedcba', 'abc')",
		"ends-with((), ())",
		"ends-with('', '')",
		"ends-with('blah', '')",
		"ends-with('blah', ())",
		"ends-with('blah', 'blah')",
		"ends-with(., 'Mocked')",
	}
	results := []bool{true, false, true, true, true, true, true, false}
	for i, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.IsType(t, (*BooleanItem)(nil), item, uut)
		assert.Equal(t, results[i], getBool(item), uut)
	}
}

func TestEndsWithInvalid(t *testing.T) {
	cases := []string{
		"ends-with()",
		"ends-with('one')",
		"ends-with('one', 5)",
		"ends-with(3, 'one')",
		"ends-with('one', 'two', 'three')",
	}
	for _, uut := range cases {
		tree := assertParses(t, uut)
		ctx := MockDefaultContext()
		_, err := tree.Evaluate(ctx)
		assert.Error(t, err, uut)
	}
}

func TestStartswith(t *testing.T) {
	cases := []string{
		"starts-with('abcdef', 'abc')",
		"starts-with('fedcba', 'def')",
		"starts-with((), ())",
		"starts-with('', '')",
		"starts-with('blah', '')",
		"starts-with('blah', ())",
		"starts-with('blah', 'blah')",
		"starts-with(., 'Mocked')",
	}
	results := []bool{true, false, true, true, true, true, true, true}
	for i, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.IsType(t, (*BooleanItem)(nil), item, uut)
		assert.Equal(t, results[i], getBool(item), uut)
	}
}

func TestStartswithInvalid(t *testing.T) {
	cases := []string{
		"starts-with()",
		"starts-with('one')",
		"starts-with('one', 5)",
		"starts-with(3, 'one')",
		"starts-with('one', 'two', 'three')",
	}
	for _, uut := range cases {
		tree := assertParses(t, uut)
		ctx := MockDefaultContext()
		_, err := tree.Evaluate(ctx)
		assert.Error(t, err, uut)
	}
}

func TestContains(t *testing.T) {
	cases := []string{
		"contains('abcdef', 'bcd')",
		"contains('fedcba', 'def')",
		"contains((), ())",
		"contains('', '')",
		"contains('blah', '')",
		"contains('blah', ())",
		"contains('blah', 'blah')",
		"contains(., 'ked')",
	}
	results := []bool{true, false, true, true, true, true, true, true}
	for i, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.IsType(t, (*BooleanItem)(nil), item, uut)
		assert.Equal(t, results[i], getBool(item), uut)
	}
}

func TestContainsInvalid(t *testing.T) {
	cases := []string{
		"contains()",
		"contains('one')",
		"contains('one', 5)",
		"contains(3, 'one')",
		"contains('one', 'two', 'three')",
	}
	for _, uut := range cases {
		tree := assertParses(t, uut)
		ctx := MockDefaultContext()
		_, err := tree.Evaluate(ctx)
		assert.Error(t, err, uut)
	}
}

func TestMatches(t *testing.T) {
	cases := []string{
		"matches('abcdef', 'bcd')",
		"matches('fedcba', 'f[b-e]+a')",
		"matches('nope', ())",
		"matches('blah', 'blah')",
		"matches(., '.*Dir')",
		"matches('MockedDir', .)",
	}
	results := []bool{false, true, false, true, true, true}
	for i, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.IsType(t, (*BooleanItem)(nil), item, uut)
		assert.Equal(t, results[i], getBool(item), uut)
	}
}

func TestMatchesInvalid(t *testing.T) {
	cases := []string{
		"matches()",
		"matches('one')",
		"matches('one', 5)",
		"matches(3, 'one')",
		"matches('one', 'two', 'three')",
	}
	for _, uut := range cases {
		tree := assertParses(t, uut)
		ctx := MockDefaultContext()
		_, err := tree.Evaluate(ctx)
		assert.Error(t, err, uut)
	}
}

func TestEmptyExists(t *testing.T) {
	cases := []string{
		"()",
		"1",
		"'hi'",
		".",
		"1.0",
	}
	results := []bool{true, false, false, false, false}
	for i, uut := range cases {
		uutEmpty := "empty(" + uut + ")"
		seq, ctx := assertEvaluates(t, uutEmpty)
		item := assertSingleton(t, ctx, seq)
		assert.IsType(t, (*BooleanItem)(nil), item, uutEmpty)
		assert.Equal(t, results[i], getBool(item), uutEmpty)
		uutExists := "exists(" + uut + ")"
		seq, ctx = assertEvaluates(t, uutExists)
		item = assertSingleton(t, ctx, seq)
		assert.IsType(t, (*BooleanItem)(nil), item, uutExists)
		assert.Equal(t, !results[i], getBool(item), uutExists)
	}
}

func TestEmptyExistsInvalid(t *testing.T) {
	cases := []string{
		"empty()", "exists()",
		"empty(1, 2)", "exists(1, 2)",
	}
	for _, uut := range cases {
		tree := assertParses(t, uut)
		ctx := MockDefaultContext()
		_, err := tree.Evaluate(ctx)
		assert.Error(t, err, uut)
	}
}

func TestNamePath(t *testing.T) {
	cases := []string{
		"name(.)",
		"path(.)",
		"name()",
		"path()",
	}
	results := []string{
		"MockedDir",
		"/MockedDir",
		"MockedDir",
		"/MockedDir",
	}
	for i, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.IsType(t, (*StringItem)(nil), item, uut)
		assert.Equal(t, results[i], getString(item), uut)
	}
}

func TestNameInvalid(t *testing.T) {
	cases := []string{
		"name('one')", "path('one')",
		"name('one', 5)", "path('one', 5)",
		"name(3, 'one')", "path(3, 'one')",
		"name('one', 'two', 'three')", "path('one', 'two', 'three')",
	}
	for _, uut := range cases {
		tree := assertParses(t, uut)
		ctx := MockDefaultContext()
		_, err := tree.Evaluate(ctx)
		assert.Error(t, err, uut)
	}
}

func TestCount(t *testing.T) {
	cases := []string{
		"count(())",
		"count(20)",
		"count((1, 2, 3))",
		"count((1 to 100)[. mod 2 eq 0])",
	}
	results := []int64{0, 1, 3, 50}
	for i, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.IsType(t, (*IntegerItem)(nil), item)
		assert.Equal(t, results[i], getInteger(item), uut)
	}
}

func TestCountInvalid(t *testing.T) {
	cases := []string{
		"count()",
		"count(1, 2)",
	}
	for _, uut := range cases {
		tree := assertParses(t, uut)
		ctx := MockDefaultContext()
		_, err := tree.Evaluate(ctx)
		assert.Error(t, err, uut)
	}
}

func TestTrueFalseNot(t *testing.T) {
	cases := []string{
		"true()",
		"false()",
		"not(true())",
		"not(false())",
		"not(0)",
		"not(1)",
		"not(1 = 1)",
		"not('blah')",
		"not(())",
	}
	results := []bool{true, false, false, true, true, false, false, false, true}
	for i, uut := range cases {
		seq, ctx := assertEvaluates(t, uut)
		item := assertSingleton(t, ctx, seq)
		assert.Equal(t, results[i], getBool(item), uut)
	}
}

func TestTrueFalseNotInvalid(t *testing.T) {
	cases := []string{
		"true(1)", "false('abc')",
		"true(1, 2)", "false((1, 2, 3), .)",
		"not()",
		"not('blah', 'blah')",
	}
	for _, uut := range cases {
		tree := assertParses(t, uut)
		ctx := MockDefaultContext()
		_, err := tree.Evaluate(ctx)
		assert.Error(t, err, uut)
	}
}
