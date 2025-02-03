package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sanity-io/litter"

	"github.com/tobiashort/gox/assert"
	"github.com/tobiashort/gox/ast"
	"github.com/tobiashort/gox/lexer"
	"github.com/tobiashort/gox/parser"
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage:
	gox tokenize FILE
`)
}

func tokenize(file string) []lexer.Token {
	data, err := os.ReadFile(file)
	assert.Nil(err)
	lexer := lexer.NewLexer()
	lexer.Tokenize(string(data))
	return lexer.Tokens
}

func parse(file string) []ast.Stmt {
	tokens := tokenize(file)
	parser := parser.NewParser()
	parser.Parse(tokens)
	return parser.Stmts
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
		tokens := tokenize(file)
		for _, token := range tokens {
			fmt.Println(token)
		}
	case "parse":
		if flag.NArg() != 2 {
			fmt.Fprintln(os.Stderr, "must provide file")
			os.Exit(1)
		}
		file := flag.Arg(1)
		ast := parse(file)
		litter.Dump(ast)
	default:
		usage()
		os.Exit(1)
	}
}
