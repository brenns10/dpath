/*
lib.go contains built-in functions of the DPath language, along with any
utilities necessary to run them.
*/

package main

import (
	"bytes"
	"errors"
	"math"
)

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
	BUILTIN_CONCAT = Builtin{
		Name: "concat", NumArgs: -1, Invoke: BuiltinConcatInvoke}
)

/*
Run the builtin function boolean(), which converts an atomic value to a boolean.
The steps are summarized by this logic, found in the specification.

https://www.w3.org/TR/xpath20/#id-ebv

1. Empty sequence -> false.
2. First item of sequence is node (file) -> true.
3. Singleton of type boolean -> value.
4. Singleton of string type -> false if zero length
5. Singleton of numeric type -> false if zero or NaN
6. Otherwise, type error.
*/
func BuiltinBooleanInvoke(ctx *Context, args ...Sequence) (Sequence, error) {
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
		if item.TypeName() == TYPE_FILE {
			// Case 2, true!
			value = true
		} else {
			// Case 6
			return nil, errors.New("type error in boolean(): sequence of non-file")
		}
	} else {
		// SINGLETON
		switch item.TypeName() {
		case TYPE_BOOLEAN:
			// Case 3
			return newSingletonSequence(item), nil
		case TYPE_INTEGER:
			// Case 5
			value = getInteger(item) != int64(0)
		case TYPE_DOUBLE:
			// Case 5
			value = getDouble(item) != float64(0.0) && !math.IsNaN(getDouble(item))
		case TYPE_STRING:
			// Case 4
			value = len(getString(item)) > 0
		default:
			errorMsg := "type error in boolean(): unexpected singleton type "
			errorMsg += item.TypeName()
			return nil, errors.New(errorMsg)
		}
	}

	return newSingletonSequence(newBooleanItem(value)), nil
}

/*
Invoke the builtin "concat" function, which takes args, converts them to strings,
and concatenates them into a single string.
*/
func BuiltinConcatInvoke(ctx *Context, args ...Sequence) (Sequence, error) {
	var buffer bytes.Buffer

	if len(args) <= 0 {
		return nil, errors.New("concat() requires at least one argument")
	}
	for _, seq := range args {
		item, err := getSingleItem(ctx, seq)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(item.ToString())
	}

	return newSingletonSequence(newStringItem(buffer.String())), nil
}

/*
Return a map of each builtin's name to its struct.
*/
func DefaultNamespace() map[string]Builtin {
	return map[string]Builtin{
		"boolean": BUILTIN_BOOLEAN,
		"concat":  BUILTIN_CONCAT,
	}
}
