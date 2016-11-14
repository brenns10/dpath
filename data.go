package main

import (
	"io"
	"os"
	"path"
	"strconv"
)

/*
This sub-package contains the interfaces and structures of the "data model".

The data model for DPath closely mirrors XPath. Everything is a "sequence" of
"items", and sequences are flat. Atomic values are sequences of length one.
*/

/*
The item interface is a "generic" interface for any type of item. The only
required operation is that it can tell us its type, so we can type check and
 cast it to its actual type.
*/
type Item interface {
	Type() *Type
	Print(w io.Writer) error
}

/*
The sequence interface is a way to provide multiple implementations of a sequence
depending on the need. For instance, we may want to use it to implement a
generator for files in a directory.
*/
type Sequence interface {
	Value() Item
	Next(ctx *Context) (bool, error)
}

/*
Every DPath expression is evaluated within a context. The context contains
information such as the current context item (usually the current directory)
and the current axis (by default, files+subdirs).
*/
type Context struct {
	ContextItem Item
	CurrentAxis Axis
	Namespace   map[string]Builtin
}

func DefaultContext() *Context {
	wd, err := os.Getwd()
	if err != nil {
		panic("Getwd() failed!")
	}
	item, err := newFileItem(wd)
	if err != nil {
		panic("Lstat() failed!")
	}
	return &Context{
		ContextItem: item,
		CurrentAxis: nil,
		Namespace:   DefaultNamespace(),
	}
}

/*
The axis is the source of file sequences. It needs to be able to get items by
their names, and also iterate over a list of items. The items could be the files
within the context directory, or it could be the files within all subdirectories
of the the context, etc.
*/
type Axis interface {
	GetByName(ctx *Context, name string) Item
	Iterate(ctx *Context) Sequence
}

/*
A "dummy" implementation of an item so that we can have some stub implementations
of interfaces.
*/
type DummyItem struct{}

func (d *DummyItem) Type() *Type { return TYPE_DUMMY }

func (d *DummyItem) Print(w io.Writer) error {
	_, err := io.WriteString(w, "dummy\n")
	return err
}

/*
An integer!
*/
type IntegerItem struct {
	Value int64
}

func (i *IntegerItem) Type() *Type { return TYPE_INTEGER }

func (i *IntegerItem) Print(w io.Writer) error {
	str := "integer:" + strconv.FormatInt(i.Value, 10) + "\n"
	_, err := io.WriteString(w, str)
	return err
}

func newIntegerItem(v int64) *IntegerItem {
	return &IntegerItem{Value: v}
}

/*
A double!
*/
type DoubleItem struct {
	Value float64
}

func (i *DoubleItem) Type() *Type { return TYPE_DOUBLE }

func (i *DoubleItem) Print(w io.Writer) error {
	str := "double:" + strconv.FormatFloat(i.Value, 'f', -1, 64) + "\n"
	_, err := io.WriteString(w, str)
	return err
}

func newDoubleItem(v float64) *DoubleItem {
	return &DoubleItem{Value: v}
}

/*
A string!
*/
type StringItem struct {
	Value string
}

func (i *StringItem) Type() *Type { return TYPE_STRING }

func (i *StringItem) Print(w io.Writer) error {
	_, err := io.WriteString(w, "string:\""+i.Value+"\"\n")
	return err
}

func newStringItem(v string) *StringItem {
	return &StringItem{Value: v}
}

/*
A boolean!
*/
type BooleanItem struct {
	Value bool
}

func (i *BooleanItem) Type() *Type { return TYPE_BOOLEAN }

func (i *BooleanItem) Print(w io.Writer) error {
	_, err := io.WriteString(w, "boolean:"+strconv.FormatBool(i.Value)+"\n")
	return err
}

func newBooleanItem(v bool) *BooleanItem {
	return &BooleanItem{Value: v}
}

/*
File item (could be a directory too)!
*/
type FileItem struct {
	Path string
	Info os.FileInfo
}

func (i *FileItem) Type() *Type { return TYPE_FILE }

func (i *FileItem) Print(w io.Writer) error {
	_, err := io.WriteString(w, "file:"+i.Path+"\n")
	return err
}

func newFileItem(absPath string) (*FileItem, error) {
	info, err := os.Lstat(absPath)
	if err != nil {
		return nil, err
	}
	return &FileItem{Path: absPath, Info: info}, nil
}

func newFileItemFromInfo(info os.FileInfo, parent string) *FileItem {
	absPath := path.Join(parent, info.Name())
	return &FileItem{Path: absPath, Info: info}
}

/*
A "dummy" implementation of a sequence for stub implementations.
It looks like it will yield a dummy item, but really it's empty since Next()
always returns false.
*/
type DummySequence struct{}

func (d *DummySequence) Value() Item {
	return &DummyItem{}
}

func (d *DummySequence) Next(ctx *Context) (bool, error) {
	return false, nil
}

/*
WrapperSequence wraps a slice of Items, potentially a slice of one.
*/
type WrapperSequence struct {
	Wrapped []Item
	Index   int
}

func newWrapperSequence(seq []Item) *WrapperSequence {
	return &WrapperSequence{Wrapped: seq, Index: -1}
}

func newSingletonSequence(item Item) *WrapperSequence {
	return &WrapperSequence{Wrapped: []Item{item}, Index: -1}
}

func newEmptySequence() *WrapperSequence {
	return &WrapperSequence{Wrapped: []Item{}, Index: -1}
}

func (s *WrapperSequence) Value() Item {
	if s.Index < 0 || s.Index >= len(s.Wrapped) {
		panic("Accessing sequence out-of-bounds.")
	}
	return s.Wrapped[s.Index]
}

func (s *WrapperSequence) Next(ctx *Context) (bool, error) {
	s.Index++
	return s.Index < len(s.Wrapped), nil
}

/*
RangeSequence creates a range of numbers, either integer or double
*/
type RangeSequence struct {
	IntCurrent int64
	IntStop    int64
	DblCurrent float64
	DblStop    float64
	IsInt      bool
}

func newIntegerRange(start, stop int64) *RangeSequence {
	return &RangeSequence{IsInt: true, IntCurrent: start - 1, IntStop: stop}
}

func newDoubleRange(start, stop float64) *RangeSequence {
	return &RangeSequence{IsInt: false, DblCurrent: start - 1, DblStop: stop}
}

func (s *RangeSequence) Value() Item {
	if s.IsInt {
		return newIntegerItem(s.IntCurrent)
	} else {
		return newDoubleItem(s.DblCurrent)
	}
}

func (s *RangeSequence) Next(ctx *Context) (bool, error) {
	if s.IsInt {
		s.IntCurrent++
		return s.IntCurrent <= s.IntStop, nil
	} else {
		s.DblCurrent++
		return s.DblCurrent <= s.DblStop, nil
	}
}

/*
ExpressionFilteredSequence is a sequence that filters by a list of predicates
which are simply expressions.
*/
type ExpressionFilter struct {
	Source  Sequence
	Current Item
	Filters []ParseTree
}

func newExpressionFilter(src Sequence, f []ParseTree) *ExpressionFilter {
	return &ExpressionFilter{Source: src, Current: nil, Filters: f}
}

func (f *ExpressionFilter) Value() Item {
	return f.Current
}

func (f *ExpressionFilter) Next(ctx *Context) (bool, error) {
	var e error = nil
OUTER:
	for r, e := f.Source.Next(ctx); r && e == nil; r, e = f.Source.Next(ctx) {
		f.Current = f.Source.Value()
		oldCtxItem := ctx.ContextItem
		ctx.ContextItem = f.Current
		for _, filter := range f.Filters {
			res, err := execBuiltin(ctx, "boolean", filter)
			if err != nil {
				ctx.ContextItem = oldCtxItem
				return false, err
			}
			if !getBool(panicUnlessOne(ctx, res)) {
				ctx.ContextItem = oldCtxItem
				continue OUTER
			}
		}
		ctx.ContextItem = oldCtxItem
		return true, nil
	}
	return false, e
}
