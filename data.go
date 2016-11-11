package main

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
	Type() string
}

/*
The sequence interface is a way to provide multiple implementations of a sequence
depending on the need. For instance, we may want to use it to implement a
generator for files in a directory.
*/
type Sequence interface {
	Value() Item
	Next() bool
}

/*
Every DPath expression is evaluated within a context. The context contains
information such as the current context item (usually the current directory)
and the current axis (by default, files+subdirs).
*/
type Context struct {
	ContextItem Item
	CurrentAxis Axis
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

func (d *DummyItem) Type() string {
	return "dummy"
}

/*
A "dummy" implementation of a sequence for stub implementations.
*/
type DummySequence struct{}

func (d *DummySequence) Value() Item {
	return &DummyItem{}
}

func (d *DummySequence) Next() bool {
	return false
}
