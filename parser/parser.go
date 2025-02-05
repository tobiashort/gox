package parser

import (
	"fmt"
	"strings"

	"github.com/tobiashort/gox/ast"
	"github.com/tobiashort/gox/lexer"
)

type Parser struct {
	Source string
	Stmts  []ast.Stmt
	Tokens []lexer.Token
	Pos    int
}

func NewParser() *Parser {
	return &Parser{
		Stmts:  make([]ast.Stmt, 0),
		Tokens: make([]lexer.Token, 0),
		Pos:    0,
	}
}

func (parser *Parser) Peek() lexer.Token {
	token := parser.Tokens[parser.Pos]
	return token
}

func (parser *Parser) Advance() lexer.Token {
	token := parser.Tokens[parser.Pos]
	parser.Pos += 1
	return token
}

func (parser *Parser) Parse(source string) {
	_lexer := lexer.NewLexer()
	_lexer.Tokenize(source)
	parser.Source = source
	parser.Tokens = _lexer.Tokens

	for {
		token := parser.Peek()
		if token.Type == lexer.TokenEOF {
			break
		}
		stmt := ParseStmt(parser, token)
		if stmt != nil {
			parser.Stmts = append(parser.Stmts, stmt)
		} else {
			parser.Advance()
		}
	}
}

func (parser *Parser) Expect(expected lexer.TokenType) lexer.Token {
	token := parser.Advance()
	if token.Type != expected {
		parser.InvalidToken(token)
	}
	return token
}

func (parser *Parser) InvalidToken(token lexer.Token) {
	line := strings.Split(parser.Source, "\n")[token.Line-1]
	line = strings.ReplaceAll(line, "\t", " ")
	cursor := "^"
	if token.Column > 0 {
		cursor = strings.Repeat("-", token.Column) + cursor
	}
	panic(fmt.Sprintf("invalid token %s at line %d column %d\n%s\n%s", token, token.Line, token.Column, line, cursor))
}
