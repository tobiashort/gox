package parser

import (
	"github.com/tobiashort/gox/ast"
	"github.com/tobiashort/gox/lexer"
)

func ParseBlockStmt(parser *Parser) ast.Stmt {
	parser.Expect(lexer.TokenBraceOpen)

	blockStmt := ast.BlockStmt{
		Body: make([]ast.Stmt, 0),
	}

	for {
		nextToken := parser.Peek()
		if nextToken.Type == lexer.TokenEOF {
			panic("reached unexpected EOF")
		}
		if nextToken.Type == lexer.TokenBraceClose {
			parser.Advance()
			break
		}
		stmt := ParseStmt(parser, nextToken)
		if stmt != nil {
			blockStmt.Body = append(blockStmt.Body, stmt)
		} else {
			parser.Advance()
		}
	}

	return blockStmt
}

func ParsePackageStmt(parser *Parser) ast.Stmt {
	parser.Expect(lexer.TokenPackage)
	nextToken := parser.Expect(lexer.TokenIdentifier)
	value := nextToken.Value
	parser.Expect(lexer.TokenNewLine)

	return ast.PackageStmt{
		PackageName: value,
	}
}

func ParseImportStmt(parser *Parser) ast.Stmt {
	parser.Expect(lexer.TokenImport)
	nextToken := parser.Expect(lexer.TokenString)
	value := nextToken.Value
	parser.Expect(lexer.TokenNewLine)

	return ast.ImportStmt{
		PackageName: value,
	}
}

func ParseFuncDeclStmt(parser *Parser) ast.Stmt {
	parser.Expect(lexer.TokenFunc)
	funcDeclStmt := ast.FuncDeclStmt{}
	nextToken := parser.Expect(lexer.TokenIdentifier)
	funcDeclStmt.FuncName = nextToken.Value

	// TODO
	parser.Expect(lexer.TokenParenOpen)
	parser.Expect(lexer.TokenParenClose)
	funcDeclStmt.Parameters = make([]ast.FuncParameter, 0)
	funcDeclStmt.ReturnTypes = make([]string, 0)

	funcDeclStmt.FuncBlock = ParseBlockStmt(parser).(ast.BlockStmt)

	return funcDeclStmt
}

func ParseExprStmt(parser *Parser) ast.Stmt {
	bindingPower := BindingPower(parser, parser.Peek())
	expr := ParseExpr(parser, bindingPower)

	return ast.ExprStmt{
		Expr: expr,
	}
}

func ParseStmt(parser *Parser, token lexer.Token) ast.Stmt {
	switch token.Type {
	case lexer.TokenNewLine:
		return nil
	case lexer.TokenIdentifier:
		return ParseExprStmt(parser)
	case lexer.TokenPackage:
		return ParsePackageStmt(parser)
	case lexer.TokenImport:
		return ParseImportStmt(parser)
	case lexer.TokenFunc:
		return ParseFuncDeclStmt(parser)
	default:
		parser.InvalidToken(token)
		return nil
	}
}
