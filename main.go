package main

import (
	"fmt"
	"os"
)

func main() {
	if _, err := ParseString(os.Args[1]); err != nil {
		fmt.Println(err)
	}
}
