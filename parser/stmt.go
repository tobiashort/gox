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
	parser.Expect(lexer.TokenNewLine)

	return ast.PackageStmt{
		PackageName: nextToken,
	}
}

func ParseImportStmt(parser *Parser) ast.Stmt {
	parser.Expect(lexer.TokenImport)
	nextToken := parser.Expect(lexer.TokenString)
	parser.Expect(lexer.TokenNewLine)

	return ast.ImportStmt{
		PackageName: nextToken,
	}
}

func ParseFuncDeclStmt(parser *Parser) ast.Stmt {
	parser.Expect(lexer.TokenFunc)
	funcDeclStmt := ast.FuncDeclStmt{}
	nextToken := parser.Expect(lexer.TokenIdentifier)
	funcDeclStmt.Name = nextToken

	// parse parameters
	parser.Expect(lexer.TokenParenOpen)
	funcDeclStmt.Parameters = make([]ast.FuncParameter, 0)
	for parser.Peek().Type != lexer.TokenParenClose {
		param := ast.FuncParameter{}
		param.Name = parser.Expect(lexer.TokenIdentifier)
		param.Type = parser.Expect(lexer.TokenIdentifier)
		funcDeclStmt.Parameters = append(funcDeclStmt.Parameters, param)
		nextToken := parser.Peek()
		if nextToken.Type == lexer.TokenComma {
			parser.Advance()
			continue
		} else if nextToken.Type == lexer.TokenParenClose {
			break
		} else {
			parser.InvalidToken(nextToken)
		}
	}
	parser.Advance()

	funcDeclStmt.ReturnTypes = make([]lexer.Token, 0)
	funcDeclStmt.Block = ParseBlockStmt(parser).(ast.BlockStmt)

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
