package main

import (
	"errors"
	"math"
)

/*
This file contains built-in functions of the DPath language, along with any
utilities necessary to run them.
*/

/*
Builtin specifies the interface that all builtin functions must satisfy.
*/
type Builtin struct {
	Name    string
	NumArgs int
	Invoke  func(ctx *Context, args ...Sequence) (Sequence, error)
}

var (
	BUILTIN_BOOLEAN = Builtin{
		Name: "boolean", NumArgs: 1, Invoke: BuiltinBooleanInvoke}
)

func BuiltinBooleanInvoke(ctx *Context, args ...Sequence) (Sequence, error) {
	// https://www.w3.org/TR/xpath20/#id-ebv
	// 1. Empty sequence -> false.
	// 2. First item of sequence is node (file) -> true.
	// 3. Singleton of type boolean -> value.
	// 4. Singleton of string type -> false if zero length
	// 5. Singleton of numeric type -> false if zero or NaN
	// 6. Otherwise, type error.
	var value bool

	// We can assume one argument, because either we're calling this from Go
	// code, and so we better only pass one argument, or else we're calling this
	// from DPath, which checks that the number of arguments equals the declared
	// number for the builtin
	arg := args[0]
	hasNext, err := arg.Next(ctx)
	if err != nil {
		return nil, err
	}
	if !hasNext {
		// Case 1, false!
		return newSingletonSequence(newBooleanItem(false)), nil
	}

	item := arg.Value()

	hasNext, err = arg.Next(ctx)
	if err != nil {
		return nil, err
	}
	if hasNext {
		// NOT A SINGLETON
		if item.Type() == TYPE_FILE {
			// Case 2, true!
			value = true
		} else {
			// Case 6
			return nil, errors.New("type error in boolean(): sequence of non-file")
		}
	} else {
		// SINGLETON
		switch item.Type() {
		case TYPE_BOOL:
			// Case 3
			return newSingletonSequence(item), nil
		case TYPE_INTEGER:
			// Case 5
			value = getInteger(item) != int64(0)
		case TYPE_DOUBLE:
			// Case 5
			value = getFloat(item) != float64(0.0) && !math.IsNaN(getFloat(item))
		case TYPE_STRING:
			// Case 4
			value = len(getString(item)) > 0
		default:
			errorMsg := "type error in boolean(): unexpected singleton type "
			errorMsg += item.Type()
			return nil, errors.New(errorMsg)
		}
	}

	return newSingletonSequence(newBooleanItem(value)), nil
}

func DefaultNamespace() map[string]Builtin {
	return map[string]Builtin{
		"boolean": BUILTIN_BOOLEAN,
	}
}
