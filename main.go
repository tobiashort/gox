package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tobiashort/gox/lexer"
	"github.com/tobiashort/gox/lexer/assert"
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage:
	gox tokenize FILE
`)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	task := ""
	if flag.NArg() > 0 {
		task = flag.Arg(0)
	}

	switch task {
	case "tokenize":
		if flag.NArg() != 2 {
			fmt.Fprintln(os.Stderr, "must provide file")
			os.Exit(1)
		}
		file := flag.Arg(1)
		data, err := os.ReadFile(file)
		assert.Nil(err)
		lexer := lexer.NewLexer()
		lexer.Tokenize(string(data))
		for _, token := range lexer.Tokens {
			fmt.Println(token)
		}
	}
}
