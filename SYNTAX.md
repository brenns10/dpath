Usage of DPath
==============

This document contains a comprehensive overview of the syntax and features DPath
offers to the user. While simply saying that DPath is "XPath for files" is
enough for those familiar with XPath to begin working with DPath, this document
exists to describe all of the features and syntax, and

Command Line
------------

DPath must be invoked on the command line with a single argument. To avoid the
shell applying globbing to the argument, enclose all DPath queries with single
quotes.

There are no command line options.

The output will be a (possibly empty) sequence of DPath Items, one per line.
They are printed in the format `type:value`, although this could be subject to
change.

Syntax
------

### Numeric Expressions

DPath supports two numeric types: integer, a 64-bit signed integer, and double,
a 64-bit double precision floating point value. Integer literals may only be in
base 10. Floating point literals may be written as decimals (e.g. `1.5`) or in
"scientific" notation (e.g. `5e-1`).

The principle numeric operators are `+ - * div idiv mod`. Most numeric operators
return the same type as their input. When an integer and double are used in a
binary operator, the integer is cast to a double, and the result is a double.
The only operators that do not obey these rules are `div`, which always returns
a double, and `idiv`, which always returns an integer.

Operators are left associative. Order of operations is "as expected" and
operations may be grouped with parentheses to enforce a particular order of
operations.

There are also unary operators `+` and `-` for negating values.

The function `round()` will round a double to the nearest integer (while still
returning a double). It returns integers unmodified. If there are two such
integers, it returns the one closest to positive infinity.

The syntax `1 to 5` (using the "to" operator) will return a sequence of numbers
starting at the first number and ending at the last, incrementing by 1. For
example, `(1 to 5) = (1, 2, 3, 4, 5)`, and `(1.0 to 3.3) = (1.0, 2.0, 3.0)`.

### Strings

DPath supports a string type. String literals may be created with single or
double quotes. There is no support for escape sequences, except that to escape
the string delimiter, it may appear twice within the string. For instance, the
string literal `'madam i''m adam'` will become `madam i'm adam`, and the string
literal `"""murder,"" she wrote"` will become `"murder", she wrote`.

Strings may also be created by calling the `name()` or `path()` function on a
file. `name()` returns a file's base name, and `path()` returns the full path of
the file.

Several string functions are supported for manipulating and creating strings:
- `string(x)` converts its argument to a string
- `concat(a, b, ...)` takes one or more arguments, converts them to strings if
  they aren't already, and concatenates them
- `string-length(x)` returns the length of a string as an integer
- `substring(s, start)`, `substring(s, start, length)` returns a substring of
  `s` starting from `start`, with `length` characters. If length is provided, it
  is assumed to be infinite. Note that indices are **one based** in XPath and
  DPath. Doubles are rounded to the nearest whole number. The exact semantics of
  substring mean that a character is included if `start <= index < start +
  length`, which means that `substring("hello", -1, 3) = "h"`, even though we
  asked for a length 3 string. It's a weird function, but that's how the XPath
  spec [defines](https://www.w3.org/TR/xpath-functions/#func-substring) it.
- `starts-with(s, prefix)`, `ends-with(s, suffix)`, `contains(s, sub)` are all
  very self explanatory, returning booleans
- `matches(s, pattern)` returns true when the entire string matches the regular
  expression given as the pattern. Note that the pattern should conform to Go's
  regular expression syntax.

### Sequence

Every expression returns a sequence in DPath. Sequences may contain zero or more
items. Most expressions return a "singleton" containing exactly one item. You
can write an empty sequence with `()` and an arbitrary sequence with `(1, 2,
3)`.

The function `count()` returns the number of items in a sequence. The function
`empty()` returns true if a sequence is empty, and `exists()` returns true if a
sequence has at least one item.

### Booleans, Comparisons, etc

There are two sets of comparison operators with an important semantic
distinction.

#### Value comparison

These are the `eq, ne, ge, gt, le, lt` operators. They compare *only* items, and
so their operands *must* be singleton sequences. For example `1 == 5` is a valid
comparison, but `1 == (1, 5)` will return an error.

#### General comparison

These are the `=, !=, >=, >, <=, <` operators. They compare sequences,
*existential* semantics. That is, a comparison will return true if it is true
for **any** pair of items from the operand sequences.

So, `1 = 1` will be true. But `1 = (1, 2)` will also be true. And even `(5, 6,
7, 8) = (8, 9, 10, 11)`. This raises some oddities. You normally expect that if
`a = b` returns true, then `a != b` returns false. But this is not necessarily
the case for general comparison! For instance, `1 = (1, 2)` is true, and so is
`1 != (1, 2)`, because in both cases, there exists a pair of items that
satisfies the comparison.

In general, if you know you're comparing two items, just use the value
comparisons!

#### Boolean Operators and Functions

Booleans are another data type, and they can be manipulated with the binary
operators `and` and `or`. The function `not()` takes a boolean and inverts it.
There are no boolean literals, but `true()` and `false()` are functions which
will return their corresponding boolean values.

The `boolean()` function takes one argument and casts it to a boolean according
to some very specific semantics:
1. If the argument is the empty sequence, return false.
2. If the first item in the sequence is a file, return true.
3. If the sequence is a singleton, and its type is boolean, return the value of
   that boolean.
4. If the sequence is a singleton, and its type is string, return false if the
   string has zero length, otherwize true.
5. If the sequence is a singleton, and its type is numeric, return false if the
   number is zero or NaN, otherwise true.
6. Finally, if none of these apply, raise an error. This case could be hit by
   something like `boolean((1, 2, 3))`.

### Paths

Finally, the meat of DPath. Path expressions have higher precedence than
multiplicative operators, but lower precedence than unary operators. They may be
rooted or not (i.e. starting with a slash). They consist of a number of "step
expressions", each separated by a slash.

Step expressions may specify an axis using `axis-name::<the rest here>`, or they
may use the default axis, which is `child::`. The attribute axis (which can only
tell you a file's size) can be accessed with the shorthand `@`. The parent axis
can be used with the shorthand `..`. An axis tells DPath what "direction" it
should "step" in. The child axis finds children of a directory. Parent has its
parent directory. Descendant is the transitive closure of child, and ancestor is
the transitive closure of parent. Descendant and Ancestor don't normally include
the object they operate on, but `descendant-or-self` and `ancestor-or-self`
exist to solve that. Finally, the `attribute` axis contains (only) the size
attribute of a file or directory.

Step expressions (other than `..`) must express some sort of test, either on the
name of the node, or on its type. A name test involves simply writing the name
as an identifier. For example `/bin` finds children of the root directory named
`bin`. If a name cannot be expressed as an identifier, the special (non-XPath)
syntax `#"literal here"` may be used in place of an identifier. For example
`./#".git"` returns the `.git` directory.

In place of an identifier, `*` may be specified, so that the path searches over
every item in the axis. For example, `//*` returns every file and directory.

A "kind" test will filter what kind of node is returned. `file()` returns files,
and `dir()` returns directories.

The `//` shorthand syntax is shorthand for `descendant-or-self::*`.

Several step expressions chained together form a path. Semantically, each step
is evaluated once for every output from the previous step. In particular, the
output of each step is used as the *context* for evaluating the next step. So in
`foo/bar`, the foo step finds a child of the current directory named foo. The
bar step takes that directory, makes it the context, and finds a child named
bar.

### Predicates

Predicates can be applied to any sequence, including a step expression. For
instance: `(1 to 10)[. mod 2 eq 0] = (2, 4, 6, 8, 10)`.

More practically, this can be used to filter a sequence of files by conditions
on their names, attributes, ancestors, etc. Predicates may contain any
expression. They are evaluated once for each sequence, and if the result is
true, that item is kept, otherwise, it is skipped.

**Divergence from XPath:** Predicates behave as if the `boolean()` function had
been called on the expression. In normal XPath, numeric predicates are compared
with the index in the sequence, which means that `(1 to 5)[2]` would return `2`
(essentially, indexing into the sequence). DPath does not have guaranteed order
for the sequences returned by its axes, and so indexing does not make sense.

Any expression may be used in a predicate, including another path, and even
another predicate.
