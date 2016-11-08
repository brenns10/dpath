DPath
=====

DPath is a work-in-progress project for my EECS 433 Databases course. It is an
implementation of an XQuery-like language for file system queries. For example,
the DPath query below would match every PNG file in the current directory tree:

```
.//file()[ends-with(name(), ".png")]
```

Currently I have a lexer and parser implemented. The main program will attempt
to parse the expression, outputting tokens as it parses. If no error message is
produced, the expression parsed successfully.

To try it, clone this repo and run the following.

```
$ go get github.com/blynn/nex
$ go generate
$ go build
$ ./dpath './/file()[ends-with(name(), ".png")]'
DOT
SLASH
SLASH
FILE
LPAREN
RPAREN
LBRACKET
QNAME
LPAREN
QNAME
LPAREN
RPAREN
COMMA
STRING_LITERAL
RPAREN
RBRACKET
```
