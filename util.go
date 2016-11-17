/*
util.go contains several private utility functions that are useful for making
code a bit simpler.
*/

package main

import (
	"errors"
)

/*
A utility function that returns a single item in a sequence, raising an error
if there are zero or >1 items in the sequence.
*/
func getSingleItem(ctx *Context, s Sequence) (Item, error) {
	r, e := s.Next(ctx)
	if !r {
		return nil, errors.New("Expected one value, found none.")
	} else if e != nil {
		return nil, e
	}
	item := s.Value()
	r, e = s.Next(ctx)
	if r {
		return nil, errors.New("Too many values provided to expression.")
	} else if e != nil {
		return nil, e
	}
	return item, nil
}

/*
A utility function that "asserts" at least one item is in a sequence, panicking
if that's not the case.
*/
func panicUnlessOne(ctx *Context, s Sequence) Item {
	r, e := s.Next(ctx)
	if !r {
		panic("There wasn't a value in the sequence.")
	} else if e != nil {
		panic("Error getting value from sequence!")
	}
	return s.Value()
}

/*
Return file value, if you're certain it's a bool.
Will panic if you're wrong.
*/
func getFile(i Item) *FileItem {
	it := i.(*FileItem)
	return it
}

/*
Return bool value, if you're certain it's a bool.
Will panic if you're wrong.
*/
func getBool(i Item) bool {
	it := i.(*BooleanItem)
	return it.Value
}

/*
Return string value, if you're certain it's a string.
Will panic if you're wrong.
*/
func getString(i Item) string {
	it := i.(*StringItem)
	return it.Value
}

/*
Return integer value, if you're certain it's an integer.
Will panic if you're wrong.
*/
func getInteger(i Item) int64 {
	it := i.(*IntegerItem)
	return it.Value
}

/*
Return double value, if you're certain it's a double.
Will panic if you're wrong.
*/
func getDouble(i Item) float64 {
	it := i.(*DoubleItem)
	return it.Value
}

/*
Return numeric value as float, if you're certain it's numeric (i.e. integer
or double).
Will panic if you're wrong.
*/
func getNumericAsFloat(i Item) float64 {
	if i.TypeName() == TYPE_INTEGER {
		return float64(getInteger(i))
	} else {
		return getDouble(i)
	}
}