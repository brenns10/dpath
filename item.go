/*
item.go contains interfaces and structures related to Items, the atomic units of
the DPath data model.
*/

package main

import (
	"io"
	"os"
	"path"
	"strconv"
)

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
An Item that can contain a 64-bit signed integer.
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
An Item that can contain a 64-bit double precision floating point number.
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
