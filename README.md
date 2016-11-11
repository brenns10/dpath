DPath
=====

DPath is a work-in-progress project for my EECS 433 Databases course. It is an
implementation of an XQuery-like language for file system queries. For example,
the DPath query below would match every PNG file in the current directory tree:

```
.//file()[ends-with(name(), ".png")]
```

Currently I have a lexer and parser implemented. The main program will attempt
to parse the expression, outputting tokens as it parses. After parsing, it will
output a (rather cryptic looking) parse tree, or an error.

To try it, clone this repo and run the following.

```
$ go get github.com/blynn/nex
$ go generate
$ go build
$ ./dpath './/file()[ends-with(name(), ".png")]'
PATH
  .
(ANY CHILD)
  FILTER EXPRESSION:
    file
  FILTER BY:
    ends-with()
      name()
      .png
```
