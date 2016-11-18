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

$ dpath '/home/stephen//.[ends-with(name(), ".png")]'
# lists PNG files in home directory

$ dpath '/home/stephen//.[contains(name(), "smb196")]'
# finds files I named with my school ID, likely for homework submissions

$ dpath '/home/stephen//.[matches(name(), ".*\w{3}\d{1,3}.*")]'
# finds files that have any school ID in them
# this turns up lots of junk because it's a common pattern

$ dpath './/.[starts-with(name(), parent::*/name())]'
# or just:
$ dpath './/.[starts-with(name(), ../name())]'
# finds files that start with their containing directory's name
file:/home/stephen/go/src/github.com/brenns10/dpath/dpath.nex
file:/home/stephen/go/src/github.com/brenns10/dpath/dpath.y
file:/home/stephen/go/src/github.com/brenns10/dpath/dpath
file:/home/stephen/go/src/github.com/brenns10/dpath/dpath.nn.go
file:/home/stephen/go/src/github.com/brenns10/dpath/.git/objects/e9/e9f542b2423e029b7adc72f71265e2eabb63a6
```

If you're interested in how this implementation works, I maintain the
file [GUIDE.md](GUIDE.md), which should give some high-level explanation of how
the language comes together. The low-level details can be read about in comments
and code.

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
* Boolean logic expressions (`and or`)
* Predicate syntax on sequences, e.g. `(1 to 5)[. mod 2 eq 0]` which evaluates
  to `(2, 4)`.
* Path expressions on the child, parent, and descendant axes
* The shorthand notations `*`, `..`, `//`, `#"spaces etc here"`
* Functions: `boolean()`, `concat()`, `round()`, `substring()`, `string()`,
  `string-length()`, `ends-with()`, `starts-with()`, `contains()`, `matches()`,
  `empty()`, `exists()`, `name()`, `path()`, `count()`.
* Selectors: `file()`, `dir()`
