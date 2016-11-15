package main

import (
	"errors"
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
		CurrentAxis: AXIS_CHILD,
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
	GetByName(ctx *Context, name string) (Sequence, error)
	Iterate(ctx *Context) (Sequence, error)
}

/*
Some axes
*/
var (
	AXIS_CHILD              = &ChildAxis{}
	AXIS_PARENT             = &ParentAxis{}
	AXIS_DESCENDANT         = &DescendantAxis{}
	AXIS_DESCENDANT_OR_SELF = &DescendantOrSelfAxis{}
)

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

/*
ConditionFilter is much simpler than ExpressionFilter. It takes a source sequence
and yields items from it that satisfy a condition function supplied to it.
*/
type ConditionFilter struct {
	Source  Sequence
	Current Item
	Filter  func(Item) bool
}

func newConditionFilter(src Sequence, cond func(Item) bool) *ConditionFilter {
	return &ConditionFilter{Source: src, Filter: cond}
}

func (f *ConditionFilter) Value() Item {
	return f.Current
}

func (f *ConditionFilter) Next(ctx *Context) (bool, error) {
	var e error = nil
	for r, e := f.Source.Next(ctx); r && e == nil; r, e = f.Source.Next(ctx) {
		f.Current = f.Source.Value()
		if f.Filter(f.Current) {
			return true, nil
		}
	}
	return false, e
}

/*
PathSequence takes items from the source, makes them the context item, and
then yields from the sequence returned by the path.
*/
type PathSequence struct {
	CtxSource  Sequence
	Expression ParseTree
	Source     Sequence
}

func newPathSequence(src Sequence, expr ParseTree) *PathSequence {
	return &PathSequence{CtxSource: src, Expression: expr, Source: nil}
}

func (s *PathSequence) Next(ctx *Context) (b bool, e error) {
	var err error = nil
	var hasNext bool
	for {
		if s.Source != nil {
			// Replace the context item with the one from the previous point
			// in the path. Then attempt to advance the source sequence.
			oldCtx := ctx.ContextItem
			ctx.ContextItem = s.CtxSource.Value()
			hasNext, err = s.Source.Next(ctx)
			ctx.ContextItem = oldCtx
			if err != nil {
				return false, err
			} else if hasNext {
				return true, nil
			}
			// Continue on if no error and the source sequence is empty.
		}

		// Get a new context item.
		hasNext, err = s.CtxSource.Next(ctx)
		if !hasNext || err != nil {
			// Return if we're out of context items, or if we've got an error.
			return hasNext, err
		}

		// Evaluate the path expression to get a source
		oldCtx := ctx.ContextItem
		ctx.ContextItem = s.CtxSource.Value()
		s.Source, err = s.Expression.Evaluate(ctx)
		ctx.ContextItem = oldCtx
		if err != nil {
			return false, err
		}
		// Fall through back to the top of the loop to try to get stuff from
		// the source again.
	}
}

func (s *PathSequence) Value() Item {
	if s.Source != nil {
		return s.Source.Value()
	} else {
		return nil
	}
}

/*
ConcatenateSequence is a sequence that takes a slice of sequences and yields from
them in order.
*/
type ConcatenateSequence struct {
	Sources []Sequence
	Current int
}

func newConcatenateSequence(sequences ...Sequence) *ConcatenateSequence {
	return &ConcatenateSequence{Sources: sequences, Current: 0}
}

func (s *ConcatenateSequence) Next(ctx *Context) (bool, error) {
	var err error = nil
	var hasNext bool
	for {
		if s.Current >= len(s.Sources) {
			return false, nil
		}
		hasNext, err = s.Sources[s.Current].Next(ctx)
		if hasNext || err != nil {
			return hasNext, err
		}
		s.Current++
	}
}

func (s *ConcatenateSequence) Value() Item {
	if s.Current <= len(s.Sources) {
		// Return current source sequence's current value.
		return s.Sources[s.Current].Value()
	} else {
		// If we're finished with all our sources, keep outputting last sequence's
		// last value.
		return s.Sources[len(s.Sources)-1].Value()
	}
}

/*
DescendentSequence implements the entire DescendantAxis, since there is not a
really nicer way to do GetByName on it. It implements a depth-first search of
the descendant tree by using the ToVisit slice as a stack.
*/
type DescendantSequence struct {
	Source  Sequence
	ToVisit []*FileItem
}

func newDescendantSequence(start *FileItem) *DescendantSequence {
	return &DescendantSequence{Source: nil, ToVisit: []*FileItem{start}}
}

func (s *DescendantSequence) Next(ctx *Context) (bool, error) {
	var err error = nil
	var hasNext bool
	for {
		if s.Source != nil {
			// Try to the source sequence
			hasNext, err = s.Source.Next(ctx)
			if err != nil {
				return false, err
			} else if hasNext {
				// If there is a next item, get it and add it to the visit
				// stack when it's a directory.
				it := s.Source.Value().(*FileItem)
				if it.Info.IsDir() {
					s.ToVisit = append(s.ToVisit, it)
				}
				return true, nil
			}
			// Continue on if no error and the source sequence is empty.
		}

		// If the visit stack is empty, we're done
		if len(s.ToVisit) <= 0 {
			return false, nil
		}

		// Get a new source.
		oldCtx := ctx.ContextItem
		ctx.ContextItem = s.ToVisit[len(s.ToVisit)-1]
		s.ToVisit = s.ToVisit[:len(s.ToVisit)-1]
		s.Source, err = AXIS_CHILD.Iterate(ctx)
		ctx.ContextItem = oldCtx

		if err != nil {
			return false, err
		}

		// Fall through back to the top of the loop to try to get stuff from
		// the source again.
	}
}

func (s *DescendantSequence) Value() Item {
	if s.Source != nil {
		// Make sure to add directories to the visit queue as we see them.
		return s.Source.Value().(*FileItem)
	} else {
		return nil
	}
}

/*
ChildAxis is the default axis for normal operation.
*/
type ChildAxis struct {
}

func newChildAxis() *ChildAxis {
	return &ChildAxis{}
}

func (a *ChildAxis) GetByName(ctx *Context, name string) (Sequence, error) {
	ctxItem, ok := ctx.ContextItem.(*FileItem)
	if !ok {
		return nil, errors.New(
			"Attempting to use ChildAxis when context item is not a file.",
		)
	}
	path := path.Join(ctxItem.Path, name)
	newItem, err := newFileItem(path)
	if err != nil {
		// assume file not found, and return empty sequence
		return newEmptySequence(), nil
	} else {
		return newSingletonSequence(newItem), nil
	}
}

func (a *ChildAxis) Iterate(ctx *Context) (Sequence, error) {
	ctxItem, ok := ctx.ContextItem.(*FileItem)
	if !ok {
		return nil, errors.New(
			"Attempting to use ChildAxis when context item is not a file.",
		)
	}
	f, err := os.Open(ctxItem.Path)
	if err != nil {
		return nil, errors.New(
			"Error while attempting to Open() context item.",
		)
	}
	contents, err := f.Readdir(0)
	if err != nil {
		f.Close()
		return nil, errors.New(
			"Error while attempting to Readdir() context item.",
		)
	}

	children := make([]Item, 0, len(contents))
	for _, info := range contents {
		children = append(children, newFileItemFromInfo(info, ctxItem.Path))
	}

	return newWrapperSequence(children), nil
}

/*
ParentAxis contains only the parent of a file.
*/
type ParentAxis struct {
}

func newParentAxis() *ParentAxis {
	return &ParentAxis{}
}

func (a *ParentAxis) GetByName(ctx *Context, name string) (Sequence, error) {
	ctxItem, ok := ctx.ContextItem.(*FileItem)
	if !ok {
		return nil, errors.New(
			"Attempting to use ParentAxis when context item is not a file.",
		)
	}
	path := path.Join(ctxItem.Path, "..")
	if path == ctxItem.Path {
		// tried to access parent of root! sneaky...
		return newEmptySequence(), nil
	}
	newItem, err := newFileItem(path)
	if err != nil {
		panic("error finding parent of file node")
	}

	// Since this is GetByName
	if newItem.Info.Name() == name {
		return newSingletonSequence(newItem), nil
	} else {
		return newEmptySequence(), nil
	}
}

func (a *ParentAxis) Iterate(ctx *Context) (Sequence, error) {
	ctxItem, ok := ctx.ContextItem.(*FileItem)
	if !ok {
		return nil, errors.New(
			"Attempting to use ParentAxis when context item is not a file.",
		)
	}
	path := path.Join(ctxItem.Path, "..")
	if path == ctxItem.Path {
		// tried to access parent of root! sneaky...
		return newEmptySequence(), nil
	}
	newItem, err := newFileItem(path)
	if err != nil {
		panic("error finding parent of file node")
	}

	return newSingletonSequence(newItem), nil
}

/*
DescendantAxis is the most tricky and inefficient axis, as it includes all sub
nodes of the context node.
*/
type DescendantAxis struct {
}

func (a *DescendantAxis) Iterate(ctx *Context) (Sequence, error) {
	source, ok := ctx.ContextItem.(*FileItem)
	if !ok {
		return nil, errors.New(
			"Attempting to use DescendantAxis when context item is not a file.",
		)
	}
	return newDescendantSequence(source), nil
}

func (a *DescendantAxis) GetByName(ctx *Context, name string) (Sequence, error) {
	seq, err := a.Iterate(ctx)
	if err != nil {
		return nil, err
	}
	return newConditionFilter(
		seq,
		func(i Item) bool {
			fi := i.(*FileItem)
			return fi.Info.Name() == name
		},
	), nil
}

/*
DescendantOrSelfAxis is just the DescendantAxis but with self concatenated.
*/
type DescendantOrSelfAxis struct {
	*DescendantAxis
}

func (a *DescendantOrSelfAxis) Iterate(ctx *Context) (Sequence, error) {
	seq, err := a.DescendantAxis.Iterate(ctx)
	if err != nil {
		return nil, err
	}
	return newConcatenateSequence(
		newSingletonSequence(ctx.ContextItem),
		seq,
	), nil
}
