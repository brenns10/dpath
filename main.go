package main

import (
	"fmt"
	"os"
)

/*
A command-line driver for evaluating DPath expressions.
*/
func main() {
	if len(os.Args) < 2 {
		fmt.Println("error: must provide DPath expression")
		return
	}

	// Parse the DPath expression.
	tree, err := ParseString(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print out the parse tree.
	fmt.Println("PARSE TREE:")
	if err = tree.Print(os.Stdout, 0); err != nil {
		fmt.Println(err)
		return
	}

	// Evaluate the expression and print the results.
	ctx := DefaultContext()
	seq, err := tree.Evaluate(ctx)
	if err != nil {
		fmt.Println("ERROR:")
		fmt.Println(err)
		return
	} else {
		fmt.Println("OUTPUT:")
	}
	for r, err := seq.Next(ctx); r && err == nil; r, err = seq.Next(ctx) {
		if err = seq.Value().Print(os.Stdout); err != nil {
			fmt.Println(err)
		}
	}
	if err != nil {
		fmt.Println("Error while iterating:")
		fmt.Println(err)
	}
}
