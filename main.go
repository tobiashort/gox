package main

import (
	"fmt"
	"os"

	"github.com/tobiashort/gox/lexer"
)

func assertNil(val any) {
	if val != nil {
		panic(val)
	}
}

func main() {
	data, err := os.ReadFile("./examples/or_panic.gox")
	assertNil(err)
	source := string(data)
	lexer := lexer.NewLexer()
	lexer.Tokenize(source)
	for _, token := range lexer.Tokens {
		fmt.Println(token.String())
	}
}
