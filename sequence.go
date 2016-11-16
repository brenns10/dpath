/*
sequence.go contains the Sequence interface and implementations thereof.
*/

package main

import (
	log "github.com/Sirupsen/logrus"
)

/*
Sequence is a general-purpose "iterator" or "generator". Everything in [XD]Path
is a sequence, so every expression evaluates to a Sequence.

Sequences may contain zero or more items. They are always "pointing" at an item,
which can be retrieved by the Value() function. The sequence can be advanced by
the Next() function, which will tell whether there is anything left. They always
start pointing BEFORE the first item, so you MUST call Next() before Value() for
the first time. After the Sequence is exhausted, Sequences continue returning
their last value from Value().

This is the standard way of iterating over a Sequence:

	var next bool
	var err error = nil
	for next, err = seq.Next(ctx); next && err == nil; next, err = seq.Next(ctx) {
		// do stuff with seq.Value()
	}
	if err != nil {
		// handle errors from Next()
	}
*/
type Sequence interface {
	Value() Item
	Next(ctx *Context) (bool, error)
}

/*
WrapperSequence is the most basic implementation of a Sequence. It wraps a slice
of Items, which could have length 0 or 1.
*/
type WrapperSequence struct {
	Wrapped []Item
	Index   int
}

/*
Return a new wrapper sequence wrapping a given slice of items. You should favor
newSingletonSequence() and newEmptySequence() if you need to create singletons
or empty sequences.
*/
func newWrapperSequence(seq []Item) *WrapperSequence {
	return &WrapperSequence{Wrapped: seq, Index: -1}
}

/*
Return a new sequence that will yield a single item. This is implemented with a
WrapperSequence, but there is no guarantee of that. Simply that it will return a
sequence that will work!
*/
func newSingletonSequence(item Item) Sequence {
	return &WrapperSequence{Wrapped: []Item{item}, Index: -1}
}

/*
Return a new sequence that is empty. Implemented with WrapperSequence but that
is not guaranteed.
*/
func newEmptySequence() Sequence {
	return &WrapperSequence{Wrapped: []Item{}, Index: -1}
}

func (s *WrapperSequence) Value() Item {
	if s.Index < 0 || s.Index >= len(s.Wrapped) {
		panic("Accessing sequence out-of-bounds.")
	}
	return s.Wrapped[s.Index]
}

func (s *WrapperSequence) Next(ctx *Context) (bool, error) {
	s.Index++
	return s.Index < len(s.Wrapped), nil
}

/*
RangeSequence creates a range of numbers, either integer or double. Although it
could be implemented with a WrapperSequence, this saves memory in the case of
large ranges. This is used to implement the "to" operator, a little known part
of XPath!
*/
type RangeSequence struct {
	IntCurrent int64
	IntStop    int64
	DblCurrent float64
	DblStop    float64
	IsInt      bool
}

/*
Return a Sequence that will yield integers from start to stop, inclusive.
*/
func newIntegerRange(start, stop int64) *RangeSequence {
	return &RangeSequence{IsInt: true, IntCurrent: start - 1, IntStop: stop}
}

/*
Return a Sequence that will yield doubles from start to stop, in increments of
1.0, inclusive.
*/
func newDoubleRange(start, stop float64) *RangeSequence {
	return &RangeSequence{IsInt: false, DblCurrent: start - 1, DblStop: stop}
}

func (s *RangeSequence) Value() Item {
	if s.IsInt {
		return newIntegerItem(s.IntCurrent)
	} else {
		return newDoubleItem(s.DblCurrent)
	}
}

func (s *RangeSequence) Next(ctx *Context) (bool, error) {
	if s.IsInt {
		s.IntCurrent++
		return s.IntCurrent <= s.IntStop, nil
	} else {
		s.DblCurrent++
		return s.DblCurrent <= s.DblStop, nil
	}
}

/*
ExpressionFilter is a "filtering" sequence, meaning that it takes a source
sequence and filters items by a condition. In this case, the condition is a list
of expressions from predicates. These are evaluated in order, and an item is
yielded only if it satisfies every expression. The expression is converted to
boolean by the built-in boolean() function -- see BuiltinBooleanInvoke().
*/
type ExpressionFilter struct {
	Source  Sequence
	Current Item
	Filters []ParseTree
}

/*
Return a new instance of ExpressionFilter.
*/
func newExpressionFilter(src Sequence, f []ParseTree) *ExpressionFilter {
	return &ExpressionFilter{Source: src, Current: nil, Filters: f}
}

func (f *ExpressionFilter) Value() Item {
	return f.Current
}

func (f *ExpressionFilter) Next(ctx *Context) (bool, error) {
	var e error = nil
	// Outer loop iterates over items from the source sequence. It terminates when
	// an item that satisfies all conditions, or when the source is exhausted.
OUTER:
	for r, e := f.Source.Next(ctx); r && e == nil; r, e = f.Source.Next(ctx) {
		// The context item needs to be set to the current item when evaluating the
		// conditions.
		f.Current = f.Source.Value()
		oldCtxItem := ctx.ContextItem
		ctx.ContextItem = f.Current
		// Inner loop goes over each expression.
		for _, filter := range f.Filters {
			res, err := execBuiltin(ctx, "boolean", filter)
			if err != nil {
				ctx.ContextItem = oldCtxItem
				return false, err
			}
			if !getBool(panicUnlessOne(ctx, res)) {
				ctx.ContextItem = oldCtxItem
				continue OUTER
			}
		}
		ctx.ContextItem = oldCtxItem
		return true, nil
	}
	return false, e
}

/*
ConditionFilter is much simpler than ExpressionFilter. Instead of using a DPath
expression as its condition for filtering, ConditionFilter uses a simple
function.
*/
type ConditionFilter struct {
	Source  Sequence
	Current Item
	Filter  func(Item) bool
}

/*
Return a ConditionFilter. cond is a function that returns true when the given
item should be yielded, or false if it should be discarded.
*/
func newConditionFilter(src Sequence, cond func(Item) bool) *ConditionFilter {
	return &ConditionFilter{Source: src, Filter: cond}
}

func (f *ConditionFilter) Value() Item {
	return f.Current
}

func (f *ConditionFilter) Next(ctx *Context) (bool, error) {
	var e error = nil
	for r, e := f.Source.Next(ctx); r && e == nil; r, e = f.Source.Next(ctx) {
		f.Current = f.Source.Value()
		if f.Filter(f.Current) {
			return true, nil
		}
	}
	return false, e
}

/*
PathSequence is used to implement each step of a path expression. It takes two
things. First, a sequence of input, generally from the previous step along the
path expression. Second, an expression to be evaluated for each Item from the
input sequence, thus producing the output of this Sequence.

The total output of a PathSequence is the concatenated output of the sequences
produced by expression when evaluated on each Item from the input Sequence.

Larger path expressions are built by simply chaining these PathSequences
together.
*/
type PathSequence struct {
	CtxSource  Sequence
	Expression ParseTree
	Source     Sequence
}

/*
Return a new PathSequence given the source sequence and step expression.
*/
func newPathSequence(src Sequence, expr ParseTree) *PathSequence {
	return &PathSequence{CtxSource: src, Expression: expr, Source: nil}
}

func (s *PathSequence) Next(ctx *Context) (b bool, e error) {
	var err error = nil
	var hasNext bool
	for {
		if s.Source != nil {
			// Replace the context item with the one from the previous point
			// in the path. Then attempt to advance the source sequence.
			oldCtx := ctx.ContextItem
			ctx.ContextItem = s.CtxSource.Value()
			hasNext, err = s.Source.Next(ctx)
			ctx.ContextItem = oldCtx
			if err != nil {
				return false, err
			} else if hasNext {
				return true, nil
			}
			// Continue on if no error and the source sequence is empty.
		}

		// Get the next input from our input sequence.
		hasNext, err = s.CtxSource.Next(ctx)
		if !hasNext || err != nil {
			// Return if we're out of context items, or if we've got an error.
			return hasNext, err
		}

		// Use the current item from our input sequence as the context for
		// evaluating the step expression, producing a new sequence of output.
		oldCtx := ctx.ContextItem
		ctx.ContextItem = s.CtxSource.Value()
		s.Source, err = s.Expression.Evaluate(ctx)
		ctx.ContextItem = oldCtx
		if err != nil {
			return false, err
		}
		// Fall through back to the top of the loop to try to get items from the
		// source again.
	}
}

func (s *PathSequence) Value() Item {
	if s.Source != nil {
		return s.Source.Value()
	} else {
		return nil
	}
}

/*
ConcatenateSequence is a sequence that takes a slice of sequences and yields
from each of them, one at a time, in order.
*/
type ConcatenateSequence struct {
	Sources []Sequence
	Current int
}

/*
Return a sequence concatenated from the given sequences.
*/
func newConcatenateSequence(sequences ...Sequence) *ConcatenateSequence {
	return &ConcatenateSequence{Sources: sequences, Current: 0}
}

func (s *ConcatenateSequence) Next(ctx *Context) (bool, error) {
	var err error = nil
	var hasNext bool
	for {
		if s.Current >= len(s.Sources) {
			return false, nil
		}
		hasNext, err = s.Sources[s.Current].Next(ctx)
		if hasNext || err != nil {
			return hasNext, err
		}
		s.Current++
	}
}

func (s *ConcatenateSequence) Value() Item {
	if s.Current <= len(s.Sources) {
		// Return current source sequence's current value.
		return s.Sources[s.Current].Value()
	} else {
		// If we're finished with all our sources, keep outputting last sequence's
		// last value.
		return s.Sources[len(s.Sources)-1].Value()
	}
}

/*
DescendentSequence is a rather tricky sequence whose job it is to return every
descendant of a file. It does this in a depth-first manner by directory.
*/
type DescendantSequence struct {
	Source  Sequence
	ToVisit []*FileItem
}

/*
Return a sequence of all descendant files of start.
*/
func newDescendantSequence(start *FileItem) *DescendantSequence {
	return &DescendantSequence{Source: nil, ToVisit: []*FileItem{start}}
}

func (s *DescendantSequence) Next(ctx *Context) (bool, error) {
	var err error = nil
	var hasNext bool
	for {
		if s.Source != nil {
			// Try to yield from the source sequence, which is the list of files in
			// the current directory.
			hasNext, err = s.Source.Next(ctx)
			if err != nil {
				return false, err
			} else if hasNext {
				// If there is a next item, get it and add it to the visit
				// stack when it's a directory.
				it := s.Source.Value().(*FileItem)
				if it.Info.IsDir() {
					log.WithFields(log.Fields{
						"axis": "DescendantAxis",
						"size": len(s.ToVisit),
						"item": it,
					}).Debug("Adding item to visit stack.")
					s.ToVisit = append(s.ToVisit, it)
				}
				return true, nil
			}
			// Continue on if no error and the source sequence is empty.
		}

		// If the visit stack is empty, we're done
		if len(s.ToVisit) <= 0 {
			log.Debug("Iteration ending (visit stack empty).")
			return false, nil
		}

		// Grab the next directory from our depth-first stack of directories.
		oldCtx := ctx.ContextItem
		ctx.ContextItem = s.ToVisit[len(s.ToVisit)-1]
		s.ToVisit = s.ToVisit[:len(s.ToVisit)-1]
		log.WithFields(log.Fields{
			"axis": "DescendantAxis",
			"size": len(s.ToVisit),
			"item": ctx.ContextItem,
		}).Debug("Starting on new source for children.")
		s.Source, err = AXIS_CHILD.Iterate(ctx)
		ctx.ContextItem = oldCtx

		if err != nil {
			return false, err
		}

		// Fall through back to the top of the loop to try to get stuff from
		// the source again.
	}
}

func (s *DescendantSequence) Value() Item {
	if s.Source != nil {
		// Make sure to add directories to the visit queue as we see them.
		return s.Source.Value().(*FileItem)
	} else {
		return nil
	}
}
