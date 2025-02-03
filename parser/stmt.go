package parser

import (
	"fmt"

	"github.com/tobiashort/gox/ast"
	"github.com/tobiashort/gox/lexer"
)

func ParseBlockStmt(parser *Parser) ast.Stmt {
	blockStmt := ast.BlockStmt{
		Body: make([]ast.Stmt, 0),
	}

	for {
		token := parser.Advance()
		if token.Type == lexer.TokenEOF {
			panic("reached unexpected EOF")
		}
		if token.Type == lexer.TokenBraceClose {
			break
		}
		parseStmt := StmtParserForTokenType(token.Type)
		if parseStmt != nil {
			blockStmt.Body = append(blockStmt.Body, parseStmt(parser))
		}
	}

	return blockStmt
}

func ParsePackageStmt(parser *Parser) ast.Stmt {
	token := parser.Expect(lexer.TokenIdentifier)
	value := token.Value.(string)
	parser.Expect(lexer.TokenNewLine)

	return ast.PackageStmt{
		PackageName: value,
	}
}

func ParseImportStmt(parser *Parser) ast.Stmt {
	token := parser.Expect(lexer.TokenString)
	value := token.Value.(string)
	parser.Expect(lexer.TokenNewLine)

	return ast.ImportStmt{
		PackageName: value,
	}
}

func ParseFuncDeclStmt(parser *Parser) ast.Stmt {
	funcDeclStmt := ast.FuncDeclStmt{}
	token := parser.Expect(lexer.TokenIdentifier)
	funcDeclStmt.FuncName = token.Value.(string)

	// TODO
	parser.Expect(lexer.TokenParenOpen)
	parser.Expect(lexer.TokenParenClose)
	funcDeclStmt.Parameters = make([]ast.FuncParameter, 0)
	funcDeclStmt.ReturnTypes = make([]string, 0)

	parser.Expect(lexer.TokenBraceOpen)
	funcDeclStmt.FuncBlock = ParseBlockStmt(parser).(ast.BlockStmt)

	return funcDeclStmt
}

func StmtParserForTokenType(tokenType lexer.TokenType) func(*Parser) ast.Stmt {
	switch tokenType {
	case lexer.TokenNewLine:
		return nil
	case lexer.TokenPackage:
		return ParsePackageStmt
	case lexer.TokenImport:
		return ParseImportStmt
	case lexer.TokenFunc:
		return ParseFuncDeclStmt
	default:
		panic(fmt.Sprintf("invalid token %s", tokenType))
	}
}
