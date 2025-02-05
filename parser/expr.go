package parser

import (
	"github.com/tobiashort/gox/ast"
	"github.com/tobiashort/gox/lexer"
)

func ParseExpr(parser *Parser, bindingPower int) ast.Expr {
	token := parser.Advance()
	left := NUD(parser, token)
	for {
		nextBindingPower := BindingPower(parser, parser.Peek())
		if nextBindingPower <= bindingPower {
			break
		}
		token = parser.Advance()
		left = LED(parser, left, token)
	}
	return left
}

func ParseSymbolExpr(token lexer.Token) ast.Expr {
	return ast.SymbolExpr{
		Symbol: token,
	}
}

func ParseNumberExpr(token lexer.Token) ast.Expr {
	return ast.NumberExpr{
		Number: token.Value,
	}
}

func ParseAssignmentExpr(parser *Parser, left ast.Expr, token lexer.Token) ast.Expr {
	return ast.AssignmentExpr{
		Left:  left,
		Right: ParseExpr(parser, BindingPower(parser, token)),
	}
}

func ParseBinaryExpr(parser *Parser, left ast.Expr, token lexer.Token) ast.Expr {
	return ast.BinaryExpr{
		Left:     left,
		Operator: token,
		Right:    ParseExpr(parser, BindingPower(parser, token)),
	}
}

func BindingPower(parser *Parser, token lexer.Token) int {
	switch token.Type {
	case lexer.TokenStar:
		return 12
	case lexer.TokenPlus:
		return 11
	case lexer.TokenAssign:
		return 2
	case lexer.TokenNumber:
		fallthrough
	case lexer.TokenIdentifier:
		fallthrough
	case lexer.TokenNewLine:
		return 0
	default:
		parser.InvalidToken(token)
		return 0
	}
}

func NUD(parser *Parser, token lexer.Token) ast.Expr {
	switch token.Type {
	case lexer.TokenIdentifier:
		return ParseSymbolExpr(token)
	case lexer.TokenNumber:
		return ParseNumberExpr(token)
	default:
		parser.InvalidToken(token)
		return nil
	}
}

func LED(parser *Parser, left ast.Expr, token lexer.Token) ast.Expr {
	switch token.Type {
	case lexer.TokenStar:
		fallthrough
	case lexer.TokenPlus:
		return ParseBinaryExpr(parser, left, token)
	case lexer.TokenAssign:
		return ParseAssignmentExpr(parser, left, token)
	default:
		parser.InvalidToken(token)
		return nil
	}
}
