package lexer

import (
	"fmt"
	"strings"
)

type Lexer struct {
	Tokens []Token
	Source string
	Pos    int
}

func NewLexer() *Lexer {
	return &Lexer{
		Tokens: make([]Token, 0),
		Source: "",
		Pos:    0,
	}
}

func (lexer *Lexer) HasMore() bool {
	return lexer.Pos < len(lexer.Source)
}

func (lexer *Lexer) Remainder() string {
	return lexer.Source[lexer.Pos:]
}

func (lexer *Lexer) Line() int {
	return len(strings.Split(lexer.Source[:lexer.Pos], "\n"))
}

func (lexer *Lexer) Column() int {
	column := lexer.Pos
	lines := strings.Split(lexer.Source[:lexer.Pos], "\n")
	for i := 0; i < len(lines)-1; i++ {
		column -= len(lines[i]) + 1
	}
	return column
}

func (lexer *Lexer) Add(token Token) {
	lexer.Tokens = append(lexer.Tokens, token)
}

func (lexer *Lexer) Tokenize(source string) {
	lexer.Tokens = make([]Token, 0)
	lexer.Source = source
	lexer.Pos = 0

	for lexer.HasMore() {
		matched := false

		for _, pattern := range Patterns {
			loc := pattern.Regex.FindStringIndex(lexer.Remainder())
			if loc != nil {
				matched = true
				pattern.Handler(lexer, pattern.Regex)
				break
			}
		}

		if !matched {
			panic(fmt.Sprintf("Unregognized token near: %s", lexer.Remainder()))
		}
	}

	lexer.Add(NewToken(TokenEOF, "", lexer.Line(), lexer.Column()))
}
