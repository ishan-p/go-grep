package main

import (
	"flag"
	"fmt"

	gogrep "github.com/ishan-p/go-grep/grep"
)

func main() {
	flag.Parse()
	searchQuery := flag.Arg(0)
	inputFile := flag.Arg(1)
	if len(searchQuery) == 0 || len(inputFile) == 0 {
		fmt.Println("Usage: ./go-grep search_query filename.txt")
		return
	}

	gogrep.Grep(searchQuery, inputFile)
}
