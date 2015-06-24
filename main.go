package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var input string

	flag.StringVar(&input, "input", "", "input file")
	flag.Parse()

	if input == "" {
		flag.Usage()
		os.Exit(1)
	}

	var err error

	f, err := os.Open(input)
	if err != nil {
		panic(err)
	}

	p := NewParser(f)

	col, err := p.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Println("Read", len(col.Games), "games")
}
