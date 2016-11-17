/*
item.go contains interfaces and structures related to Items, the atomic units of
the DPath data model. Since much of arithmetic is type-specific, most of the
arithmetic evaluation is implemented in this file.
*/

package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"strconv"
)

const (
	TYPE_INTEGER = "integer"
	TYPE_DOUBLE  = "double"
	TYPE_BOOLEAN = "boolean"
	TYPE_STRING  = "string"
	TYPE_FILE    = "file"
)

/*
The item interface is a "generic" interface for any type of item. It is
implemented for each type. The interface should implement any language
operations required of the types.

TypeName() should return, well, exactly that. Print() should print the item to
the Writer given. Compare() and the Eval functions should perform the evaluation
operations specified.
*/
type Item interface {
	TypeName() string
	RelativeCompare() bool
	Print(w io.Writer) error
	Compare(other Item) (int64, error)
	EvalPlus(other Item) (Sequence, error)
	EvalMinus(other Item) (Sequence, error)
	EvalMultiply(other Item) (Sequence, error)
	EvalDivide(other Item) (Sequence, error)
	EvalIntegerDivide(other Item) (Sequence, error)
	EvalModulus(other Item) (Sequence, error)
	EvalTo(other Item) (Sequence, error)
	EvalAnd(other Item) (Sequence, error)
	EvalOr(other Item) (Sequence, error)
}

/*
Contains some base implementations of type operations.
*/
type BaseItem struct{}

func unsupported(operator string, leftType string, right Item) (Sequence, error) {
	return nil, errors.New(fmt.Sprintf(
		"operator %s not supported on types %s, %s",
		operator, leftType, right.TypeName(),
	))
}

func (i *BaseItem) EvalPlus(other Item) (Sequence, error) {
	return unsupported("+", "", other)
}

func (i *BaseItem) EvalMinus(other Item) (Sequence, error) {
	return unsupported("-", "", other)
}

func (i *BaseItem) EvalMultiply(other Item) (Sequence, error) {
	return unsupported("+", "", other)
}

func (i *BaseItem) EvalDivide(other Item) (Sequence, error) {
	return unsupported("div", "", other)
}

func (i *BaseItem) EvalIntegerDivide(other Item) (Sequence, error) {
	return unsupported("idiv", "", other)
}

func (i *BaseItem) EvalModulus(other Item) (Sequence, error) {
	return unsupported("mod", "", other)
}

func (i *BaseItem) EvalTo(other Item) (Sequence, error) {
	return unsupported("to", "", other)
}

func (i *BaseItem) EvalAnd(other Item) (Sequence, error) {
	return unsupported("and", "", other)
}

func (i *BaseItem) EvalOr(other Item) (Sequence, error) {
	return unsupported("or", "", other)
}

/*
Compares a double and an int by casting the integer to a double.
*/
func compareDoubleAndInt(left int64, right float64) int64 {
	if float64(left) == right {
		return int64(0)
	} else if float64(left) < right {
		return int64(-1)
	} else {
		return int64(1)
	}
}

/*
Return an error that two items' types are not comparable.
*/
func incomparableError(left, right Item) error {
	return errors.New(fmt.Sprintf(
		"Not comparable types: %s, %s",
		left.TypeName(), right.TypeName(),
	))
}

/*
Return an error that two items may not be compared "relatively". Relatively means
le, lt, ge, gt.
*/
func noRelCmpError(left, right Item) error {
	return errors.New(fmt.Sprintf(
		"Illegal relative comparison between types: %s, %s",
		left.TypeName(), right.TypeName(),
	))
}

/*
An Item that can contain a 64-bit signed integer.
*/
type IntegerItem struct {
	*BaseItem
	Value int64
}

func (i *IntegerItem) TypeName() string { return TYPE_INTEGER }

func (i *IntegerItem) RelativeCompare() bool { return true }

func (i *IntegerItem) Print(w io.Writer) error {
	str := "integer:" + strconv.FormatInt(i.Value, 10) + "\n"
	_, err := io.WriteString(w, str)
	return err
}

func (i *IntegerItem) Compare(right Item) (int64, error) {
	switch right.TypeName() {
	case TYPE_INTEGER:
		return i.Value - getInteger(right), nil
	case TYPE_DOUBLE:
		return compareDoubleAndInt(i.Value, getDouble(right)), nil
	default:
		return int64(0), incomparableError(i, right)
	}
}

func (i *IntegerItem) EvalPlus(right Item) (Sequence, error) {
	switch right.TypeName() {
	case TYPE_INTEGER:
		return newSingletonSequence(newIntegerItem(i.Value + getInteger(right))), nil
	case TYPE_DOUBLE:
		return newSingletonSequence(newDoubleItem(float64(i.Value) + getDouble(right))), nil
	default:
		return nil, incomparableError(i, right)
	}
}

func (i *IntegerItem) EvalMinus(right Item) (Sequence, error) {
	switch right.TypeName() {
	case TYPE_INTEGER:
		return newSingletonSequence(newIntegerItem(i.Value - getInteger(right))), nil
	case TYPE_DOUBLE:
		return newSingletonSequence(newDoubleItem(float64(i.Value) - getDouble(right))), nil
	default:
		return nil, incomparableError(i, right)
	}
}

func (i *IntegerItem) EvalMultiply(right Item) (Sequence, error) {
	switch right.TypeName() {
	case TYPE_INTEGER:
		return newSingletonSequence(newIntegerItem(i.Value * getInteger(right))), nil
	case TYPE_DOUBLE:
		return newSingletonSequence(newDoubleItem(float64(i.Value) * getDouble(right))), nil
	default:
		return nil, incomparableError(i, right)
	}
}

func (i *IntegerItem) EvalDivide(right Item) (Sequence, error) {
	switch right.TypeName() {
	case TYPE_INTEGER:
		return newSingletonSequence(newDoubleItem(float64(i.Value) / getNumericAsFloat(right))), nil
	case TYPE_DOUBLE:
		return newSingletonSequence(newDoubleItem(float64(i.Value) / getDouble(right))), nil
	default:
		return nil, incomparableError(i, right)
	}
}

func (i *IntegerItem) EvalIntegerDivide(right Item) (Sequence, error) {
	switch right.TypeName() {
	case TYPE_INTEGER:
		return newSingletonSequence(newIntegerItem(i.Value / getInteger(right))), nil
	case TYPE_DOUBLE:
		return newSingletonSequence(newIntegerItem(i.Value / int64(getDouble(right)))), nil
	default:
		return nil, incomparableError(i, right)
	}
}

func (i *IntegerItem) EvalModulus(right Item) (Sequence, error) {
	switch right.TypeName() {
	case TYPE_INTEGER:
		return newSingletonSequence(newIntegerItem(i.Value % getInteger(right))), nil
	case TYPE_DOUBLE:
		return newSingletonSequence(newDoubleItem(math.Mod(float64(i.Value), getDouble(right)))), nil
	default:
		return nil, incomparableError(i, right)
	}
}

func (i *IntegerItem) EvalTo(right Item) (Sequence, error) {
	switch right.TypeName() {
	case TYPE_INTEGER:
		return newIntegerRange(i.Value, getInteger(right)), nil
	case TYPE_DOUBLE:
		return newDoubleRange(float64(i.Value), getDouble(right)), nil
	default:
		return nil, incomparableError(i, right)
	}
}

func newIntegerItem(v int64) *IntegerItem {
	return &IntegerItem{Value: v}
}

/*
An Item that can contain a 64-bit double precision floating point number.
*/
type DoubleItem struct {
	*BaseItem
	Value float64
}

func (i *DoubleItem) TypeName() string { return TYPE_DOUBLE }

func (i *DoubleItem) RelativeCompare() bool { return true }

func (i *DoubleItem) Print(w io.Writer) error {
	str := "double:" + strconv.FormatFloat(i.Value, 'f', -1, 64) + "\n"
	_, err := io.WriteString(w, str)
	return err
}

func (i *DoubleItem) Compare(right Item) (int64, error) {
	switch right.TypeName() {
	case TYPE_DOUBLE:
		rightDouble := getDouble(right)
		if i.Value == rightDouble {
			return int64(0), nil
		} else if i.Value < rightDouble {
			return int64(-1), nil
		} else {
			return int64(1), nil
		}
	case TYPE_INTEGER:
		return -compareDoubleAndInt(getInteger(right), i.Value), nil
	default:
		return int64(0), incomparableError(i, right)
	}
}

func (i *DoubleItem) EvalPlus(right Item) (Sequence, error) {
	if right.TypeName() != TYPE_INTEGER && right.TypeName() != TYPE_DOUBLE {
		return nil, incomparableError(i, right)
	}
	return newSingletonSequence(newDoubleItem(i.Value + getNumericAsFloat(right))), nil
}

func (i *DoubleItem) EvalMinus(right Item) (Sequence, error) {
	if right.TypeName() != TYPE_INTEGER && right.TypeName() != TYPE_DOUBLE {
		return nil, incomparableError(i, right)
	}
	return newSingletonSequence(newDoubleItem(i.Value - getNumericAsFloat(right))), nil
}

func (i *DoubleItem) EvalMultiply(right Item) (Sequence, error) {
	if right.TypeName() != TYPE_INTEGER && right.TypeName() != TYPE_DOUBLE {
		return nil, incomparableError(i, right)
	}
	return newSingletonSequence(newDoubleItem(i.Value * getNumericAsFloat(right))), nil
}

func (i *DoubleItem) EvalDivide(right Item) (Sequence, error) {
	if right.TypeName() != TYPE_INTEGER && right.TypeName() != TYPE_DOUBLE {
		return nil, incomparableError(i, right)
	}
	return newSingletonSequence(newDoubleItem(i.Value / getNumericAsFloat(right))), nil
}

func (i *DoubleItem) EvalIntegerDivide(right Item) (Sequence, error) {
	if right.TypeName() != TYPE_INTEGER && right.TypeName() != TYPE_DOUBLE {
		return nil, incomparableError(i, right)
	}
	return newSingletonSequence(newIntegerItem(int64(i.Value / getNumericAsFloat(right)))), nil
}

func (i *DoubleItem) EvalModulus(right Item) (Sequence, error) {
	if right.TypeName() != TYPE_INTEGER && right.TypeName() != TYPE_DOUBLE {
		return nil, incomparableError(i, right)
	}
	return newSingletonSequence(newDoubleItem(math.Mod(i.Value, getNumericAsFloat(right)))), nil
}

func (i *DoubleItem) EvalTo(right Item) (Sequence, error) {
	if right.TypeName() != TYPE_INTEGER && right.TypeName() != TYPE_DOUBLE {
		return nil, incomparableError(i, right)
	}
	return newDoubleRange(i.Value, getNumericAsFloat(right)), nil
}

func newDoubleItem(v float64) *DoubleItem {
	return &DoubleItem{Value: v}
}

/*
A string!
*/
type StringItem struct {
	*BaseItem
	Value string
}

func (i *StringItem) TypeName() string { return TYPE_STRING }

func (i *StringItem) RelativeCompare() bool { return true }

func (i *StringItem) Print(w io.Writer) error {
	_, err := io.WriteString(w, "string:\""+i.Value+"\"\n")
	return err
}

func (i *StringItem) Compare(right Item) (int64, error) {
	if right.TypeName() != TYPE_STRING {
		return int64(0), incomparableError(i, right)
	}
	rightString := getString(right)
	if i.Value == rightString {
		return int64(0), nil
	} else if i.Value < rightString {
		return int64(-1), nil
	} else {
		return int64(1), nil
	}
}

func newStringItem(v string) *StringItem {
	return &StringItem{Value: v}
}

/*
A boolean!
*/
type BooleanItem struct {
	*BaseItem
	Value bool
}

func (i *BooleanItem) TypeName() string { return TYPE_BOOLEAN }

func (i *BooleanItem) RelativeCompare() bool { return true }

func (i *BooleanItem) Print(w io.Writer) error {
	_, err := io.WriteString(w, "boolean:"+strconv.FormatBool(i.Value)+"\n")
	return err
}

func newBooleanItem(v bool) *BooleanItem {
	return &BooleanItem{Value: v}
}

func (i *BooleanItem) Compare(right Item) (int64, error) {
	if right.TypeName() != TYPE_BOOLEAN {
		return int64(0), incomparableError(i, right)
	}
	rightBool := getBool(right)
	if i.Value == rightBool {
		return int64(0), nil
	} else if i.Value && !rightBool {
		return int64(1), nil
	} else {
		return int64(-1), nil
	}
}

func (i *BooleanItem) EvalAnd(right Item) (Sequence, error) {
	if right.TypeName() != TYPE_BOOLEAN {
		return unsupported("and", "boolean", right)
	}
	return newSingletonSequence(newBooleanItem(i.Value && getBool(right))), nil
}

func (i *BooleanItem) EvalOr(right Item) (Sequence, error) {
	if right.TypeName() != TYPE_BOOLEAN {
		return unsupported("or", "boolean", right)
	}
	return newSingletonSequence(newBooleanItem(i.Value || getBool(right))), nil
}

/*
File item (could be a directory too)!
*/
type FileItem struct {
	*BaseItem
	Path string
	Info os.FileInfo
}

func (i *FileItem) TypeName() string { return TYPE_FILE }

func (i *FileItem) RelativeCompare() bool { return false }

func (i *FileItem) Print(w io.Writer) error {
	_, err := io.WriteString(w, "file:"+i.Path+"\n")
	return err
}

func (i *FileItem) Compare(right Item) (int64, error) {
	if right.TypeName() != TYPE_FILE {
		return int64(0), incomparableError(i, right)
	}
	if i.Path == getFile(right).Path {
		return int64(0), nil
	} else {
		return int64(1), nil
	}
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
Value comparison functions!
*/

func CmpEq(left, right Item) (bool, error) {
	res, err := left.Compare(right)
	return res == 0, err
}

func CmpNe(left, right Item) (bool, error) {
	res, err := left.Compare(right)
	return res != 0, err
}

func CmpLe(left, right Item) (bool, error) {
	if !left.RelativeCompare() || !right.RelativeCompare() {
		return false, noRelCmpError(left, right)
	}
	res, err := left.Compare(right)
	return res <= 0, err
}

func CmpLt(left, right Item) (bool, error) {
	if !left.RelativeCompare() || !right.RelativeCompare() {
		return false, noRelCmpError(left, right)
	}
	res, err := left.Compare(right)
	return res < 0, err
}

func CmpGe(left, right Item) (bool, error) {
	if !left.RelativeCompare() || !right.RelativeCompare() {
		return false, noRelCmpError(left, right)
	}
	res, err := left.Compare(right)
	return res >= 0, err
}

func CmpGt(left, right Item) (bool, error) {
	if !left.RelativeCompare() || !right.RelativeCompare() {
		return false, noRelCmpError(left, right)
	}
	res, err := left.Compare(right)
	return res > 0, err
}

/*
Do a "general" comparison between two sequences. A "general" comparison means
that it searches for any pair of items which satisfies the comparison. On the
other hand, a value comparison (see above) can only happen between two Items.
*/
func GeneralComparison(ctx *Context, left, right Sequence,
	comparator func(left, right Item) (bool, error)) (Sequence, error) {

	var e error = nil
	var b, cmp bool
	var l, r Item

	// need to load the whole left side before we can begin comparing
	leftSlice := make([]Item, 0)
	for b, e = left.Next(ctx); b && e == nil; b, e = left.Next(ctx) {
		leftSlice = append(leftSlice, left.Value())
	}
	if e != nil {
		return nil, e
	}

	// now look for any combination of left and right that satisfies
	for b, e = right.Next(ctx); b && e == nil; b, e = right.Next(ctx) {
		for _, l = range leftSlice {
			r = right.Value()
			cmp, e = comparator(l, r)
			if e != nil {
				return nil, e
			}
			if cmp {
				return newSingletonSequence(newBooleanItem(true)), nil
			}
		}
	}
	if e != nil {
		return nil, e
	}
	return newSingletonSequence(newBooleanItem(false)), nil
}
