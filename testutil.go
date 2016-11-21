/*
testutil.go contains utility functions and mocks for testing DPath
*/

package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"strings"
	"testing"
	"time"
)

/*
MockAxis implements the Axis interface, but it will do no file I/O and it will
record calls, so that we can run tests on Axis operations.
*/
type MockAxis struct {
	mock.Mock
	AxisName string
}

func (a *MockAxis) GetByName(ctx *Context, name string) (Sequence, error) {
	a.Called(ctx, name)
	return newEmptySequence(), nil
}

func (a *MockAxis) Iterate(ctx *Context) (Sequence, error) {
	a.Called(ctx)
	return newEmptySequence(), nil
}

/*
MockFileInfo mocks the os.FileInfo struct so that we can safely test things
related to files.
*/
type MockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (f *MockFileInfo) Name() string       { return f.name }
func (f *MockFileInfo) Size() int64        { return f.size }
func (f *MockFileInfo) Mode() os.FileMode  { return f.mode }
func (f *MockFileInfo) ModTime() time.Time { return f.modTime }
func (f *MockFileInfo) IsDir() bool        { return f.isDir }
func (f *MockFileInfo) Sys() interface{}   { return nil }

/*
Return a mocked file.
*/
func MockFile(path string, name string, isDir bool) *FileItem {
	return &FileItem{
		Path: path,
		Info: &MockFileInfo{
			name:    name,
			size:    1024,
			mode:    os.FileMode(0),
			modTime: time.Time{},
			isDir:   isDir,
		},
	}
}

/*
Return a context object with everything mocked. This should be used only in tests,
and tests should use only this (NOT the DefaultContext()).
*/
func MockDefaultContext() *Context {
	axes := map[string]Axis{
		"child":              &MockAxis{AxisName: "child"},
		"parent":             &MockAxis{AxisName: "parent"},
		"descendant":         &MockAxis{AxisName: "descendant"},
		"descendant-or-self": &MockAxis{AxisName: "descendant-or-self"},
		"ancestor":           &MockAxis{AxisName: "ancestor"},
		"ancestor-or-self":   &MockAxis{AxisName: "ancestor-or-self"},
		"attribute":          &MockAxis{AxisName: "attribute"},
	}
	return &Context{
		Axes:        axes,
		ContextItem: MockFile("/MockedDir", "MockedDir", true),
		CurrentAxis: axes["child"],
		Namespace:   DefaultNamespace(),
	}
}

/*
 * Assert that a string either does not lex as the token type, or it panics.
 */
func assertDoesNotLex(t *testing.T, s string, tok int) {
	var sym yySymType
	defer recover()
	l := NewLexer(strings.NewReader(s))
	assert.NotEqual(t, l.Lex(&sym), tok)
}

/*
Assert that a string parses, and return the parse tree.
*/
func assertParses(t *testing.T, s string) ParseTree {
	tree, e := ParseString(s)
	assert.Nil(t, e, s)
	return tree
}

func assertEvaluatesCtx(t *testing.T, s string, ctx *Context) Sequence {
	tree, e := ParseString(s)
	assert.Nil(t, e, s)
	res, e := tree.Evaluate(ctx)
	assert.Nil(t, e, s)
	return res
}

/*
Assert that a string evaluates without error, and return the Sequence generated
along with the Context used.
*/
func assertEvaluates(t *testing.T, s string) (Sequence, *Context) {
	ctx := MockDefaultContext()
	return assertEvaluatesCtx(t, s, ctx), ctx
}

/*
Assert that a sequence is a singleton.
*/
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

/*
Assert that a sequence is empty.
*/
func assertEmptySequence(t *testing.T, ctx *Context, seq Sequence) {
	hasNext, err := seq.Next(ctx)
	assert.Nil(t, err)
	assert.False(t, hasNext)
}

/*
Assert that a string parses to a LiteralTree and return it.
*/
func assertLiteral(t *testing.T, s string) *LiteralTree {
	tree := assertParses(t, s)
	lt, ok := tree.(*LiteralTree)
	assert.True(t, ok)
	return lt
}

/*
Assert that a string parses to a NameTree and return it.
*/
func assertQName(t *testing.T, s string) *NameTree {
	tree := assertParses(t, s)
	nt, ok := tree.(*NameTree)
	assert.True(t, ok)
	return nt
}

/*
Assert that a string parses to a BinopTree and return it.
*/
func assertBinop(t *testing.T, s string) *BinopTree {
	tree := assertParses(t, s)
	bt, ok := tree.(*BinopTree)
	assert.True(t, ok)
	return bt
}

/*
Assert that a string parses to an EmptySequenceTree and return it.
*/
func assertEmptySequenceTree(t *testing.T, s string) *EmptySequenceTree {
	tree := assertParses(t, s)
	bt, ok := tree.(*EmptySequenceTree)
	assert.True(t, ok)
	return bt
}
