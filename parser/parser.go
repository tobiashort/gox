package parser

import (
	"fmt"

	"github.com/tobiashort/gox/ast"
	"github.com/tobiashort/gox/lexer"
)

type Parser struct {
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

func (parser *Parser) Advance() lexer.Token {
	token := parser.Tokens[parser.Pos]
	parser.Pos += 1
	return token
}

func (parser *Parser) Parse(tokens []lexer.Token) {
	parser.Tokens = tokens

	for {
		token := parser.Advance()
		if token.Type == lexer.TokenEOF {
			break
		}
		parseStmt := StmtParserForTokenType(token.Type)
		if parseStmt != nil {
			stmt := parseStmt(parser)
			parser.Stmts = append(parser.Stmts, stmt)
		}
	}
}

func (parser *Parser) Expect(expected lexer.TokenType) lexer.Token {
	token := parser.Advance()
	if token.Type != expected {
		panic(fmt.Sprintf("invalid token %s", token))
	}
	return token
}
