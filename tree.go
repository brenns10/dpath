package main

// The following comments instruct go's build system on how to generate
// the lexer and parser.
//go:generate nex dpath.nex
//go:generate go tool yacc dpath.y

type ParseTree struct {
}
