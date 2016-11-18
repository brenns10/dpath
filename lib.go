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
	"path"
	"regexp"
	"strings"
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
	BUILTIN_STRING = Builtin{
		Name: "string", NumArgs: -1, Invoke: BuiltinStringInvoke}
	BUILTIN_STRING_LENGTH = Builtin{
		Name: "string-length", NumArgs: -1, Invoke: BuiltinStringLengthInvoke}
	BUILTIN_ENDS_WITH = Builtin{
		Name: "ends-with", NumArgs: 2, Invoke: BuiltinEndsWithInvoke}
	BUILTIN_STARTS_WITH = Builtin{
		Name: "starts-with", NumArgs: 2, Invoke: BuiltinStartsWithInvoke}
	BUILTIN_CONTAINS = Builtin{
		Name: "contains", NumArgs: 2, Invoke: BuiltinContainsInvoke}
	BUILTIN_MATCHES = Builtin{
		Name: "matches", NumArgs: 2, Invoke: BuiltinMatchesInvoke}
	BUILTIN_EMPTY = Builtin{
		Name: "empty", NumArgs: 1, Invoke: BuiltinEmptyInvoke}
	BUILTIN_EXISTS = Builtin{
		Name: "exists", NumArgs: 1, Invoke: BuiltinExistsInvoke}
	BUILTIN_NAME = Builtin{
		Name: "name", NumArgs: -1, Invoke: BuiltinNameInvoke}
	BUILTIN_PATH = Builtin{
		Name: "path", NumArgs: -1, Invoke: BuiltinPathInvoke}
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

func BuiltinStringInvoke(ctx *Context, args ...Sequence) (Sequence, error) {
	var str string
	var err error
	if len(args) == 0 {
		str = ctx.ContextItem.ToString()
	} else if len(args) == 1 {
		str, err = coerceGetString(ctx, args[0])
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("string() takes zero or one argument")
	}

	return newSingletonSequence(newStringItem(str)), nil
}

func BuiltinStringLengthInvoke(ctx *Context, args ...Sequence) (Sequence, error) {
	var str string
	var err error
	if len(args) == 0 {
		str = ctx.ContextItem.ToString()
	} else if len(args) == 1 {
		str, err = funcGetString(ctx, args[0])
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("string-length() takes zero or one argument")
	}

	return newSingletonSequence(newIntegerItem(int64(len(str)))), nil
}

func BuiltinEndsWithInvoke(ctx *Context, args ...Sequence) (Sequence, error) {
	str1, err := funcGetString(ctx, args[0])
	if err != nil {
		return nil, err
	}
	str2, err := funcGetString(ctx, args[1])
	if err != nil {
		return nil, err
	}
	return newSingletonSequence(newBooleanItem(strings.HasSuffix(str1, str2))), nil
}

func BuiltinStartsWithInvoke(ctx *Context, args ...Sequence) (Sequence, error) {
	str1, err := funcGetString(ctx, args[0])
	if err != nil {
		return nil, err
	}
	str2, err := funcGetString(ctx, args[1])
	if err != nil {
		return nil, err
	}
	return newSingletonSequence(newBooleanItem(strings.HasPrefix(str1, str2))), nil
}

func BuiltinContainsInvoke(ctx *Context, args ...Sequence) (Sequence, error) {
	str1, err := funcGetString(ctx, args[0])
	if err != nil {
		return nil, err
	}
	str2, err := funcGetString(ctx, args[1])
	if err != nil {
		return nil, err
	}
	return newSingletonSequence(newBooleanItem(strings.Contains(str1, str2))), nil
}

func BuiltinMatchesInvoke(ctx *Context, args ...Sequence) (Sequence, error) {
	str1, err := funcGetString(ctx, args[0])
	if err != nil {
		return nil, err
	}
	str2, err := funcGetString(ctx, args[1])
	if err != nil {
		return nil, err
	}
	re, err := regexp.Compile(str2)
	if err != nil {
		return nil, err
	}
	found := re.FindString(str1)
	return newSingletonSequence(newBooleanItem(found == str1)), nil
}

func BuiltinEmptyInvoke(ctx *Context, args ...Sequence) (Sequence, error) {
	hasNext, err := args[0].Next(ctx)
	if err != nil {
		return nil, err
	} else {
		return newSingletonSequence(newBooleanItem(!hasNext)), nil
	}
}

func BuiltinExistsInvoke(ctx *Context, args ...Sequence) (Sequence, error) {
	hasNext, err := args[0].Next(ctx)
	if err != nil {
		return nil, err
	} else {
		return newSingletonSequence(newBooleanItem(hasNext)), nil
	}
}

func BuiltinNameInvoke(ctx *Context, args ...Sequence) (Sequence, error) {
	var item Item
	var err error
	if len(args) == 1 {
		item, err = getSingleItem(ctx, args[0])
		if err != nil {
			return nil, err
		}
	} else if len(args) == 0 {
		item = ctx.ContextItem
	} else {
		return nil, errors.New("wrong number of arguments to name()")
	}
	if item.TypeName() != TYPE_FILE {
		return nil, errors.New("name() expects argument of type file)")
	}
	return newSingletonSequence(newStringItem(getFile(item).Info.Name())), nil
}

func BuiltinPathInvoke(ctx *Context, args ...Sequence) (Sequence, error) {
	var item Item
	var err error
	if len(args) == 1 {
		item, err = getSingleItem(ctx, args[0])
		if err != nil {
			return nil, err
		}
	} else if len(args) == 0 {
		item = ctx.ContextItem
	} else {
		return nil, errors.New("wrong number of arguments to path()")
	}
	if item.TypeName() != TYPE_FILE {
		return nil, errors.New("path() expects argument of type file)")
	}
	file := getFile(item)
	p := path.Join(file.Path, file.Info.Name())
	return newSingletonSequence(newStringItem(p)), nil
}

/*
Return a map of each builtin's name to its struct.
*/
func DefaultNamespace() map[string]Builtin {
	return map[string]Builtin{
		"boolean":       BUILTIN_BOOLEAN,
		"concat":        BUILTIN_CONCAT,
		"round":         BUILTIN_ROUND,
		"substring":     BUILTIN_SUBSTRING,
		"string":        BUILTIN_STRING,
		"string-length": BUILTIN_STRING_LENGTH,
		"ends-with":     BUILTIN_ENDS_WITH,
		"starts-with":   BUILTIN_STARTS_WITH,
		"contains":      BUILTIN_CONTAINS,
		"matches":       BUILTIN_MATCHES,
		"empty":         BUILTIN_EMPTY,
		"exists":        BUILTIN_EXISTS,
		"name":          BUILTIN_NAME,
		"path":          BUILTIN_PATH,
	}
}
