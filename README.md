DPath
=====

DPath is a work-in-progress project for my EECS 433 Databases course. It is an
implementation of an XQuery-like language for file system queries. For example,
the DPath query below would match every PNG file in the current directory tree:

```
.//.[ends-with(name(), ".png")]
```

Usage
-----

### Setup

This will get you where you need to go, although there may be a "better way".

```bash
$ export GOPATH=~/go
$ export PATH=$GOPATH/bin:$PATH
$ go get -t github.com/blynn/nex github.com/brenns10/dpath
# this will give an error, ignore it
$ cd $GOPATH/src/github.com/brenns10/dpath
$ go generate
$ go build
$ go test
```

### Usage

You can `go install` once you've done `go generate`, which will put the `dpath`
command in your Go binary directory, which is hopefully in your `$PATH`. From
there, try some queries:

```bash
$ dpath './/.'
# recursively lists this subdirectory and everything under it

$ dpath '//.'
# recursively lists everything in the filesystem
# on second thought, don't do that

$ dpath '(1 to 10)[. mod 3 eq 1]'
# lists numbers with remainder 1 when divided by 3 :)
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
