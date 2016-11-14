package main

import (
	"errors"
	"fmt"
)

type Type struct {
	Name    string
	Compare func(left, right Item, rel bool) (int64, error)
}

/*
Declarations for types!
*/

var (
	TYPE_INTEGER = &Type{
		Name:    "integer",
		Compare: nil,
	}
	TYPE_DOUBLE = &Type{
		Name:    "double",
		Compare: nil,
	}
	TYPE_BOOLEAN = &Type{
		Name:    "boolean",
		Compare: nil,
	}
	TYPE_STRING = &Type{
		Name:    "string",
		Compare: nil,
	}
	TYPE_FILE = &Type{
		Name:    "file",
		Compare: nil,
	}
	TYPE_DUMMY = &Type{
		Name:    "dummy",
		Compare: nil,
	}
)

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
		left.Type(), right.Type(),
	))
}

/*
Return an error that two items may not be compared "relatively". Relatively means
le, lt, ge, gt.
*/
func noRelCmpError(left, right Item) error {
	return errors.New(fmt.Sprintf(
		"Illegal relative comparison between types: %s, %s",
		left.Type(), right.Type(),
	))
}

/*
Actual comparison functions, one for each type.
*/

func IntegerCompare(left, right Item, rel bool) (int64, error) {
	leftInt := getInteger(left)
	switch right.Type() {
	case TYPE_INTEGER:
		return leftInt - getInteger(right), nil
	case TYPE_DOUBLE:
		return compareDoubleAndInt(leftInt, getFloat(right)), nil
	default:
		return int64(0), incomparableError(left, right)
	}
}

func DoubleCompare(left, right Item, rel bool) (int64, error) {
	leftDouble := getFloat(left)
	switch right.Type() {
	case TYPE_INTEGER:
		return -compareDoubleAndInt(getInteger(right), leftDouble), nil
	case TYPE_DOUBLE:
		rightDouble := getFloat(right)
		if leftDouble == rightDouble {
			return int64(0), nil
		} else if leftDouble < rightDouble {
			return int64(-1), nil
		} else {
			return int64(1), nil
		}
	default:
		return int64(0), incomparableError(left, right)
	}
}

func BooleanCompare(left, right Item, rel bool) (int64, error) {
	leftBool := getBool(left)
	if right.Type() != TYPE_BOOLEAN {
		return int64(0), incomparableError(left, right)
	}
	rightBool := getBool(right)
	if leftBool == rightBool {
		return int64(0), nil
	} else if leftBool && !rightBool {
		return int64(1), nil
	} else {
		return int64(-1), nil
	}
}

func StringCompare(left, right Item, rel bool) (int64, error) {
	leftString := getString(left)
	if right.Type() != TYPE_STRING {
		return int64(0), incomparableError(left, right)
	}
	rightString := getString(right)
	if leftString == rightString {
		return int64(0), nil
	} else if leftString < rightString {
		return int64(-1), nil
	} else {
		return int64(1), nil
	}
}

func FileCompare(left, right Item, rel bool) (int64, error) {
	leftFile := getFile(left)
	if right.Type() != TYPE_FILE {
		return int64(0), incomparableError(left, right)
	}
	if leftFile.Path == getFile(right).Path {
		return int64(0), nil
	} else if rel {
		return int64(0), noRelCmpError(left, right)
	} else {
		return int64(1), nil
	}
}

func DummyCompare(left, right Item, rel bool) (int64, error) {
	return int64(0), incomparableError(left, right)
}

/*
We need an init function to set the Compare attribute of the structs, since
the compare functions refer to the structs they're defining.
*/
func init() {
	TYPE_INTEGER.Compare = IntegerCompare
	TYPE_DOUBLE.Compare = DoubleCompare
	TYPE_BOOLEAN.Compare = BooleanCompare
	TYPE_STRING.Compare = StringCompare
	TYPE_DUMMY.Compare = DummyCompare
	TYPE_FILE.Compare = FileCompare
}

/*
Try using the left type item to compare with the right. If that fails, try using
the right type item to compare with the left.
*/
func cmpTryBoth(left, right Item, rel bool) (int64, error) {
	res, err := left.Type().Compare(left, right, rel)
	if err != nil {
		res, err = right.Type().Compare(right, left, rel)
		res = -res
	}
	return res, err
}

/*
Value comparison functions!
*/

func CmpEq(left, right Item) (Sequence, error) {
	res, err := cmpTryBoth(left, right, false)
	return newSingletonSequence(newBooleanItem(res == 0)), err
}

func CmpNe(left, right Item) (Sequence, error) {
	res, err := cmpTryBoth(left, right, false)
	return newSingletonSequence(newBooleanItem(res != 0)), err
}

func CmpLe(left, right Item) (Sequence, error) {
	res, err := cmpTryBoth(left, right, true)
	return newSingletonSequence(newBooleanItem(res <= 0)), err
}

func CmpLt(left, right Item) (Sequence, error) {
	res, err := cmpTryBoth(left, right, true)
	return newSingletonSequence(newBooleanItem(res < 0)), err
}

func CmpGe(left, right Item) (Sequence, error) {
	res, err := cmpTryBoth(left, right, true)
	return newSingletonSequence(newBooleanItem(res >= 0)), err
}

func CmpGt(left, right Item) (Sequence, error) {
	res, err := cmpTryBoth(left, right, true)
	return newSingletonSequence(newBooleanItem(res > 0)), err
}
