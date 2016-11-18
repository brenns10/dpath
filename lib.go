/*
lib.go contains built-in functions of the DPath language, along with any
utilities necessary to run them.
*/

package main

import (
	"bytes"
	"errors"
	"fmt"
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
	BUILTIN_ROUND = Builtin{
		Name: "round", NumArgs: 1, Invoke: BuiltinRoundInvoke}
	BUILTIN_SUBSTRING = Builtin{
		Name: "substring", NumArgs: -1, Invoke: BuiltinSubstringInvoke}
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
Invoke the builtin "round" function, which takes a numeric and rounds it.
*/
func BuiltinRoundInvoke(ctx *Context, args ...Sequence) (Sequence, error) {
	item, err := getSingleItem(ctx, args[0])
	if err != nil {
		return nil, err
	}
	switch item.TypeName() {
	case TYPE_INTEGER:
		return newSingletonSequence(item), nil
	case TYPE_DOUBLE:
		return newSingletonSequence(newDoubleItem(math.Floor(
			getDouble(item) + 0.5,
		))), nil
	default:
		return nil, errors.New(fmt.Sprintf(
			"Type %s not supported for round() function.", item.TypeName(),
		))
	}
}

/*
Invoke the builtin "substring" function, which takes a string, a start index,
and an optional length, and returns the substring starting at the start index
with the given length.

This is based on the specification:
https://www.w3.org/TR/xpath-functions/#func-substring

XPath string and item indices are 1 based, not 0 based (yuck). Also, the spec's
interpretation of the semantics of the substring operation makes for the
absolute STUPIDEST substring function in the world. But such is life, I'm a spec
implementer, not writer ¯\_(ツ)_/¯
*/
func BuiltinSubstringInvoke(ctx *Context, args ...Sequence) (Sequence, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, errors.New("substring requires 2 or 3 arguments")
	}
	// Get and check first argument. The spec has some pretty stupid semantics.
	item1, err := getZeroOrOne(ctx, args[0])
	var str string
	if err != nil {
		return nil, err
	} else if item1 == nil {
		str = ""
	} else if item1.TypeName() != TYPE_STRING {
		return nil, errors.New(
			"first arg to substring must be empty sequence or string",
		)
	} else {
		str = getString(item1)
	}
	// Get and check second argument
	item2, err := getSingleItem(ctx, args[1])
	if err != nil {
		return nil, err
	}
	if item2.TypeName() != TYPE_INTEGER && item2.TypeName() != TYPE_DOUBLE {
		return nil, errors.New("second arg to substring must be numeric")
	}
	var start, end, strlen int64
	strlen = int64(len(str))
	start = getNumericAsInteger(item2) - 1 // (positions are 1 based)
	end = strlen
	// Get and apply the optional third argument.
	if len(args) == 3 {
		item3, err := getSingleItem(ctx, args[2])
		if err != nil {
			return nil, err
		}
		if item3.TypeName() != TYPE_INTEGER && item3.TypeName() != TYPE_DOUBLE {
			return nil, errors.New("third arg to substring must be numeric")
		}
		end = start + getNumericAsInteger(item3)
	}
	// Normalize start.
	if start < 0 {
		start = 0
	} else if start > strlen {
		start = int64(len(str))
	}
	// Normalize end
	if end < start {
		end = start
	} else if end > strlen {
		end = strlen
	}
	return newSingletonSequence(newStringItem(str[start:end])), nil
}

/*
Return a map of each builtin's name to its struct.
*/
func DefaultNamespace() map[string]Builtin {
	return map[string]Builtin{
		"boolean":   BUILTIN_BOOLEAN,
		"concat":    BUILTIN_CONCAT,
		"round":     BUILTIN_ROUND,
		"substring": BUILTIN_SUBSTRING,
	}
}
