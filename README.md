DPath
=====

DPath is a work-in-progress project for my EECS 433 Databases course. It is an
implementation of an XQuery-like language for file system queries. For example,
the DPath query below would match every PNG file in the current directory tree:

```
.//*[ends-with(name(), ".png")]
```

Setup
-----

I have many features implemented (see below for details). To try it out, follow
the steps below. Make sure you enclose queries in single quotes so that the
shell doesn't do any funny business.

```
$ go get github.com/blynn/nex
$ go generate
$ go build
$ ./dpath '../../*/*'
PARSE TREE:
PATH
  ..
  ..
  *
  *
OUTPUT:
file:/home/stephen/go/src/github.com/brenns10/dpath
file:/home/stephen/go/src/github.com/brenns10/gochat
file:/home/stephen/go/src/github.com/stretchr/testify
file:/home/stephen/go/src/github.com/blynn/nex
```

Status
------

### Completed

* Lexer
* Parser
* Arithmetic expressions on int, double involving operators `+ - * div idiv mod`
* Range expressions for numeric types, e.g. `(1 to 5)` which evaluates to `(1,
  2, 3, 4, 5)`.
* Value (`= != <= < >= >`) and General (`eq ne le lt ge gt`) comparisons. The
  difference being that Value requires singletons, whereas General will look for
  any pair of atomics in the input sequences that satisfy the comparisons.
* Predicate syntax on sequences, e.g. `(1 to 5)[. mod 2 eq 0]` which evaluates
  to `(2, 4)`.
* Path expressions on the child, parent, and descendant axes
* The shorthand notations `*`, `..`, `//`
* The `boolean()` function

### To-Do

* Ancestor axis
* Attribute axes along with the shorthand `@`
* Children file and directory kinds
* String handling functions:
    * `concat(args as string...) as string`
    * `substring(s as string, start as integer, end as integer?) as string`
    * `len(s as string) as int`
    * `ends-with(s1 as string, suffix as string) as string`
    * `match(regex as string, target as string)`
* File functions:
    * `name(arg as file) as string`
    * `path(arg as file) as string`
* Shorthand syntax for using names that aren't supported XPath QNames
* Logging
* Testing

### Improvements

* Currently, the predicate feature is not properly implemented with regard to
  numeric indices. It doesn't really matter because I don't guarantee the order
  of items from axes. So a future improvement is to guarantee forward/reverse
  order of files, and then properly implement numeric indices.
