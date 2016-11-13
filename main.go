package main

import (
	"fmt"
	"os"
)

func main() {
	tree, err := ParseString(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = tree.Print(os.Stdout, 0); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Output:")

	ctx := DefaultContext()
	seq, err := tree.Evaluate(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	for seq.Next() {
		if err = seq.Value().Print(os.Stdout); err != nil {
			fmt.Println(err)
		}
	}
}
