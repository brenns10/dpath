package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
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
		f.Close()
		log.WithFields(log.Fields{
			"error": err,
			"axis":  "ChildAxis",
		}).Warn("Error encountered while calling Open().")
		return newEmptySequence(), nil
	}
	contents, err := f.Readdir(0)
	if err != nil {
		f.Close()
		log.WithFields(log.Fields{
			"error": err,
			"axis":  "ChildAxis",
		}).Warn("Error encountered while calling Readdir().")
		return newEmptySequence(), nil
	}

	children := make([]Item, 0, len(contents))
	for _, info := range contents {
		children = append(children, newFileItemFromInfo(info, ctxItem.Path))
	}
	f.Close()

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
