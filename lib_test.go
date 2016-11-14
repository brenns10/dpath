package main

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestBooleanEmptySequence(t *testing.T) {
	ctx := DefaultContext()
	s, e1 := BuiltinBooleanInvoke(ctx, newEmptySequence())
	assert.Nil(t, e1)
	i, e2 := getSingleItem(s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.False(t, v.Value)
}

func TestBooleanEmptyString(t *testing.T) {
	ctx := DefaultContext()
	input := newSingletonSequence(newStringItem(""))

	s, e1 := BuiltinBooleanInvoke(ctx, input)
	assert.Nil(t, e1)
	i, e2 := getSingleItem(s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.False(t, v.Value)
}

func TestBooleanNonEmptyString(t *testing.T) {
	ctx := DefaultContext()
	input := newSingletonSequence(newStringItem("foo"))

	s, e1 := BuiltinBooleanInvoke(ctx, input)
	assert.Nil(t, e1)
	i, e2 := getSingleItem(s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.True(t, v.Value)
}

func TestBooleanIntegerZero(t *testing.T) {
	ctx := DefaultContext()
	input := newSingletonSequence(newIntegerItem(int64(0)))

	s, e1 := BuiltinBooleanInvoke(ctx, input)
	assert.Nil(t, e1)
	i, e2 := getSingleItem(s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.False(t, v.Value)
}

func TestBooleanIntegerNonzero(t *testing.T) {
	ctx := DefaultContext()
	input := newSingletonSequence(newIntegerItem(int64(1)))

	s, e1 := BuiltinBooleanInvoke(ctx, input)
	assert.Nil(t, e1)
	i, e2 := getSingleItem(s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.True(t, v.Value)
}

func TestBooleanDoubleZero(t *testing.T) {
	ctx := DefaultContext()
	input := newSingletonSequence(newDoubleItem(float64(0.0)))

	s, e1 := BuiltinBooleanInvoke(ctx, input)
	assert.Nil(t, e1)
	i, e2 := getSingleItem(s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.False(t, v.Value)
}

func TestBooleanDoubleNan(t *testing.T) {
	ctx := DefaultContext()
	input := newSingletonSequence(newDoubleItem(math.NaN()))

	s, e1 := BuiltinBooleanInvoke(ctx, input)
	assert.Nil(t, e1)
	i, e2 := getSingleItem(s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.False(t, v.Value)
}

func TestBooleanDoubleTrue(t *testing.T) {
	ctx := DefaultContext()
	input := newSingletonSequence(newDoubleItem(float64(1.1)))

	s, e1 := BuiltinBooleanInvoke(ctx, input)
	assert.Nil(t, e1)
	i, e2 := getSingleItem(s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.True(t, v.Value)
}

func TestBooleanFalseBool(t *testing.T) {
	ctx := DefaultContext()
	input := newSingletonSequence(newBooleanItem(false))

	s, e1 := BuiltinBooleanInvoke(ctx, input)
	assert.Nil(t, e1)
	i, e2 := getSingleItem(s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.False(t, v.Value)
}

func TestBooleanTrueBool(t *testing.T) {
	ctx := DefaultContext()
	input := newSingletonSequence(newBooleanItem(true))

	s, e1 := BuiltinBooleanInvoke(ctx, input)
	assert.Nil(t, e1)
	i, e2 := getSingleItem(s)
	assert.Nil(t, e2)
	v, ok := i.(*BooleanItem)
	assert.True(t, ok)

	assert.True(t, v.Value)
}
