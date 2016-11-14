package main

import (
	"errors"
	"fmt"
	"math"
)

type Type struct {
	Name              string
	Compare           func(left, right Item, rel bool) (int64, error)
	EvalPlus          func(left, right Item) (Sequence, error)
	EvalMinus         func(left, right Item) (Sequence, error)
	EvalMultiply      func(left, right Item) (Sequence, error)
	EvalDivide        func(left, right Item) (Sequence, error)
	EvalIntegerDivide func(left, right Item) (Sequence, error)
	EvalModulus       func(left, right Item) (Sequence, error)
	EvalTo            func(left, right Item) (Sequence, error)
}

func Unsupported(operator string) func(left, right Item) (Sequence, error) {
	return func(left, right Item) (Sequence, error) {
		return nil, errors.New(fmt.Sprintf(
			"operator %s is not supported on types %s, %s",
			operator, right.Type().Name, left.Type().Name,
		))
	}
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
		Name:              "boolean",
		EvalPlus:          Unsupported("+"),
		EvalMinus:         Unsupported("-"),
		EvalMultiply:      Unsupported("*"),
		EvalDivide:        Unsupported("div"),
		EvalIntegerDivide: Unsupported("idiv"),
		EvalModulus:       Unsupported("mod"),
		EvalTo:            Unsupported("to"),
	}
	TYPE_STRING = &Type{
		Name:              "string",
		Compare:           nil,
		EvalPlus:          Unsupported("+"),
		EvalMinus:         Unsupported("-"),
		EvalMultiply:      Unsupported("*"),
		EvalDivide:        Unsupported("div"),
		EvalIntegerDivide: Unsupported("idiv"),
		EvalModulus:       Unsupported("mod"),
		EvalTo:            Unsupported("to"),
	}
	TYPE_FILE = &Type{
		Name:              "file",
		Compare:           nil,
		EvalPlus:          Unsupported("+"),
		EvalMinus:         Unsupported("-"),
		EvalMultiply:      Unsupported("*"),
		EvalDivide:        Unsupported("div"),
		EvalIntegerDivide: Unsupported("idiv"),
		EvalModulus:       Unsupported("mod"),
		EvalTo:            Unsupported("to"),
	}
	TYPE_DUMMY = &Type{
		Name:              "dummy",
		Compare:           nil,
		EvalPlus:          Unsupported("+"),
		EvalMinus:         Unsupported("-"),
		EvalMultiply:      Unsupported("*"),
		EvalDivide:        Unsupported("div"),
		EvalIntegerDivide: Unsupported("idiv"),
		EvalModulus:       Unsupported("mod"),
		EvalTo:            Unsupported("to"),
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
Value comparison functions!
*/

func CmpEq(left, right Item) (Sequence, error) {
	res, err := left.Type().Compare(left, right, false)
	return newSingletonSequence(newBooleanItem(res == 0)), err
}

func CmpNe(left, right Item) (Sequence, error) {
	res, err := left.Type().Compare(left, right, false)
	return newSingletonSequence(newBooleanItem(res != 0)), err
}

func CmpLe(left, right Item) (Sequence, error) {
	res, err := left.Type().Compare(left, right, true)
	return newSingletonSequence(newBooleanItem(res <= 0)), err
}

func CmpLt(left, right Item) (Sequence, error) {
	res, err := left.Type().Compare(left, right, true)
	return newSingletonSequence(newBooleanItem(res < 0)), err
}

func CmpGe(left, right Item) (Sequence, error) {
	res, err := left.Type().Compare(left, right, true)
	return newSingletonSequence(newBooleanItem(res >= 0)), err
}

func CmpGt(left, right Item) (Sequence, error) {
	res, err := left.Type().Compare(left, right, true)
	return newSingletonSequence(newBooleanItem(res > 0)), err
}

/*
Utility function for binary arithmetic operators. Performs boilerplate
type checking. Takes two items and a function to call when both are integers,
as well as a function to call when one is a double.
*/
func evalArithmetic(left Item, right Item,
	bothInt func(l, r int64) Item,
	otherwise func(l, r float64) Item,
) (Sequence, error) {
	// Ensure that both arguments are numeric.
	types := []*Type{TYPE_INTEGER, TYPE_DOUBLE}
	if !typeCheck(types, left, right) {
		errMsg := "Expected arguments of type integer or double, got "
		errMsg += left.Type().Name + " and " + right.Type().Name + "."
		return nil, errors.New(errMsg)
	}
	// When both arguments are integers, do integer addition.
	if typeCheck([]*Type{TYPE_INTEGER}, left, right) {
		res := bothInt(getInteger(left), getInteger(right))
		return newSingletonSequence(res), nil
	}
	// Otherwise, up-cast to double.
	res := otherwise(getNumericAsFloat(left), getNumericAsFloat(right))
	return newSingletonSequence(res), nil
}

func EvalPlusID(left Item, right Item) (Sequence, error) {
	return evalArithmetic(
		left, right,
		func(l, r int64) Item { return newIntegerItem(l + r) },
		func(l, r float64) Item { return newDoubleItem(l + r) },
	)
}

func EvalMinusID(left, right Item) (Sequence, error) {
	return evalArithmetic(
		left, right,
		func(l, r int64) Item { return newIntegerItem(l - r) },
		func(l, r float64) Item { return newDoubleItem(l - r) },
	)
}

func EvalMultiplyID(left, right Item) (Sequence, error) {
	return evalArithmetic(
		left, right,
		func(l, r int64) Item { return newIntegerItem(l * r) },
		func(l, r float64) Item { return newDoubleItem(l * r) },
	)
}

func EvalDivideID(left, right Item) (Sequence, error) {
	return evalArithmetic(
		left, right,
		func(l, r int64) Item { return newDoubleItem(float64(l) / float64(r)) },
		func(l, r float64) Item { return newDoubleItem(l / r) },
	)
}

func EvalIntegerDivideID(left, right Item) (Sequence, error) {
	return evalArithmetic(
		left, right,
		func(l, r int64) Item { return newIntegerItem(l / r) },
		func(l, r float64) Item { return newIntegerItem(int64(l / r)) },
	)
}

func EvalModulusID(left, right Item) (Sequence, error) {
	return evalArithmetic(
		left, right,
		func(l, r int64) Item { return newIntegerItem(l % r) },
		func(l, r float64) Item { return newDoubleItem(math.Mod(l, r)) },
	)
}

func EvalToID(left, right Item) (Sequence, error) {
	if left.Type() == TYPE_INTEGER && right.Type() == TYPE_INTEGER {
		return newIntegerRange(getInteger(left), getInteger(right)), nil
	} else if left.Type() == TYPE_DOUBLE && right.Type() == TYPE_DOUBLE {
		return newDoubleRange(getFloat(left), getFloat(right)), nil
	} else {
		return nil, errors.New("mismatched or undefined types in range expression")
	}
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

	TYPE_INTEGER.EvalPlus = EvalPlusID
	TYPE_INTEGER.EvalMinus = EvalMinusID
	TYPE_INTEGER.EvalMultiply = EvalMultiplyID
	TYPE_INTEGER.EvalDivide = EvalDivideID
	TYPE_INTEGER.EvalIntegerDivide = EvalIntegerDivideID
	TYPE_INTEGER.EvalModulus = EvalModulusID
	TYPE_INTEGER.EvalTo = EvalToID

	TYPE_DOUBLE.EvalPlus = EvalPlusID
	TYPE_DOUBLE.EvalMinus = EvalMinusID
	TYPE_DOUBLE.EvalMultiply = EvalMultiplyID
	TYPE_DOUBLE.EvalDivide = EvalDivideID
	TYPE_DOUBLE.EvalIntegerDivide = EvalIntegerDivideID
	TYPE_DOUBLE.EvalModulus = EvalModulusID
	TYPE_DOUBLE.EvalTo = EvalToID
}
