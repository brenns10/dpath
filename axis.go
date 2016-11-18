/*
axis.go contains data structures related to the context and axes.
*/

package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"os"
	"path"
)

/*
An axis is the source for data in path expressions. You can think of an axis as
a "direction" that you can travel from an item. It should be able to take a
context item and give us additional items in that direction from the node. For
instance, in a file system, the child axis returns contents of a directory.

The interface for Axis specifies two functions. You could get away with just an
Iterate() function, but GetByName() can make finding files within an axis much
more efficient... why list a directory and then search through the listing when
you could just stat() the file and handle the error if it doesn't exist?

Anyway, GetByName() returns a sequence of items from the Axis matching a name.
Iterate() returns all the items in the axis (from the context item) in a
sequence.
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
	// don't bother trying with files
	if !ctxItem.Info.IsDir() {
		return newEmptySequence(), nil
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
AncestorAxis contains only the parent of a file.
*/
type AncestorAxis struct {
}

func newAncestorAxis() *AncestorAxis {
	return &AncestorAxis{}
}

func (a *AncestorAxis) GetByName(ctx *Context, name string) (Sequence, error) {
	seq, err := a.Iterate(ctx)
	if err != nil {
		return nil, err
	}
	return newConditionFilter(seq, func(it Item) bool {
		return getFile(it).Info.Name() == name
	}), nil
}

func (a *AncestorAxis) Iterate(ctx *Context) (Sequence, error) {
	ctxItem, ok := ctx.ContextItem.(*FileItem)
	if !ok {
		return nil, errors.New(
			"Attempting to use AncestorAxis when context item is not a file.",
		)
	}

	ancestors := make([]Item, 0, 5)
	p := ctxItem.Path
	for path.Join(p, "..") != p {
		p = path.Join(p, "..")
		newItem, err := newFileItem(p)
		if err != nil {
			panic("error finding parent of file node")
		}
		ancestors = append(ancestors, newItem)
	}
	return newWrapperSequence(ancestors), nil
}

/*
DescendantOrSelfAxis is just the DescendantAxis but with self added in.
*/
type AncestorOrSelfAxis struct {
	*AncestorAxis
}

func (a *AncestorOrSelfAxis) GetByName(ctx *Context, name string) (Sequence, error) {
	seq, err := a.Iterate(ctx)
	if err != nil {
		return nil, err
	}
	return newConditionFilter(seq, func(it Item) bool {
		return getFile(it).Info.Name() == name
	}), nil
}

func (a *AncestorOrSelfAxis) Iterate(ctx *Context) (Sequence, error) {
	seq, err := a.AncestorAxis.Iterate(ctx)
	if err != nil {
		return nil, err
	}
	return newConcatenateSequence(
		newSingletonSequence(ctx.ContextItem),
		seq,
	), nil
}

/*
DescendantAxis returns children and children of children. Its implementation is
mostly found within the DescendantSequence.
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
DescendantOrSelfAxis is just the DescendantAxis but with self added in.
*/
type DescendantOrSelfAxis struct {
	*DescendantAxis
}

func (a *DescendantOrSelfAxis) GetByName(ctx *Context, name string) (Sequence, error) {
	seq, err := a.Iterate(ctx)
	if err != nil {
		return nil, err
	}
	return newConditionFilter(seq, func(it Item) bool {
		return getFile(it).Info.Name() == name
	}), nil
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

/*
AttributeAxis gives file metadata. Unfortunately it's difficult to get this
in a cross platform way. So the only current attribute I have is size.
*/
type AttributeAxis struct{}

func (a *AttributeAxis) GetByName(ctx *Context, name string) (Sequence, error) {
	source, ok := ctx.ContextItem.(*FileItem)
	if !ok {
		return nil, errors.New(
			"Attempting to use AttributeAxis when context item is not a file.",
		)
	}
	switch name {
	case "size":
		return newSingletonSequence(newIntegerItem(source.Info.Size())), nil
	default:
		return newEmptySequence(), nil
	}
}

func (a *AttributeAxis) Iterate(ctx *Context) (Sequence, error) {
	source, ok := ctx.ContextItem.(*FileItem)
	if !ok {
		return nil, errors.New(
			"Attempting to use AttributeAxis when context item is not a file.",
		)
	}
	return newSingletonSequence(newIntegerItem(source.Info.Size())), nil
}

/*
Every DPath expression is evaluated within a context. The context contains
information such as the current context item (usually the current directory)
and the current axis (by default, children).
*/
type Context struct {
	ContextItem Item
	CurrentAxis Axis
	Namespace   map[string]Builtin
	Axes        map[string]Axis
}

/*
DefaultContext returns a Context object where the current item is the current
directory, the axis is the child axis, and the namespace is filled with all the
builtin functions. You need to call this to get a context before evaluating
a parsed expression.
*/
func DefaultContext() *Context {
	wd, err := os.Getwd()
	if err != nil {
		panic("Getwd() failed!")
	}
	item, err := newFileItem(wd)
	if err != nil {
		panic("Lstat() failed!")
	}
	axes := map[string]Axis{
		"child":              &ChildAxis{},
		"parent":             &ParentAxis{},
		"descendant":         &DescendantAxis{},
		"descendant-or-self": &DescendantOrSelfAxis{},
		"ancestor":           &AncestorAxis{},
		"ancestor-or-self":   &AncestorOrSelfAxis{},
		"attribute":          &AttributeAxis{},
	}
	return &Context{
		ContextItem: item,
		CurrentAxis: axes["child"],
		Namespace:   DefaultNamespace(),
		Axes:        axes,
	}
}
