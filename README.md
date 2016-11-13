DPath
=====

DPath is a work-in-progress project for my EECS 433 Databases course. It is an
implementation of an XQuery-like language for file system queries. For example,
the DPath query below would match every PNG file in the current directory tree:

```
.//file()[ends-with(name(), ".png")]
```

Currently I have a lexer, parser, and support for evaluating basic arithmetic
expressions. The main program takes a DPath expression as its first argument and
outputs a parse tree followed by the result of evaluating the expression.

To try it, clone this repo and run the following.

```
$ go get github.com/blynn/nex
$ go generate
$ go build
$ ./dpath '1 * (3 - 1)'
PARSE TREE:
  1
*
    3
  -
    1
OUTPUT:
integer:2
```
