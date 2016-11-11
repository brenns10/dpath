package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	tree, err := ParseString(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	w := bufio.NewWriter(os.Stdout)
	if err = tree.Print(w, 0); err != nil {
		fmt.Println(err)
	}
	w.Flush()

	w.WriteString("Output:\n")

	ctx := DefaultContext()
	seq := tree.Evaluate(ctx)
	for seq.Next() {
		if err = seq.Value().Print(w); err != nil {
			fmt.Println(err)
		}
	}

	w.Flush()
}
