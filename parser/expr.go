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
		Number: token,
	}
}

func ParseStringExpr(token lexer.Token) ast.Expr {
	return ast.StringExpr{
		String: token,
	}
}

func ParseDeclAssignExpr(parser *Parser, left ast.Expr, token lexer.Token) ast.Expr {
	return ast.DeclAssignExpr{
		Left:  left,
		Right: ParseExpr(parser, BindingPower(parser, token)),
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

func ParseDotExpr(parser *Parser, left ast.Expr, token lexer.Token) ast.Expr {
	_, leftIsSymbol := left.(ast.SymbolExpr)

	right := ParseExpr(parser, BindingPower(parser, token))
	_, rightIsSymbol := right.(ast.SymbolExpr)

	if leftIsSymbol && rightIsSymbol {
		return ast.AccessExpr{
			Instance: left,
			Field:    right,
		}
	}

	parser.InvalidToken(token)
	return nil
}

func ParseParenOpenExpr(parser *Parser, left ast.Expr, token lexer.Token) ast.Expr {
	if left == nil {
		expr := ParseExpr(parser, 1)
		parser.Expect(lexer.TokenParenClose)
		return expr
	}

	_, isSymbolExpr := left.(ast.SymbolExpr)
	_, isAccessExpr := left.(ast.AccessExpr)

	if isSymbolExpr || isAccessExpr {
		funcCallExpr := ast.FuncCallExpr{}
		funcCallExpr.Func = left
		if parser.Peek().Type == lexer.TokenParenClose {
			parser.Advance()
			return funcCallExpr
		}
		args := ParseExpr(parser, BindingPower(parser, token))
		funcCallExpr.Args = args
		parser.Expect(lexer.TokenParenClose)
		return funcCallExpr
	}

	parser.InvalidToken(token)
	return nil
}

func ParseListExpr(parser *Parser, left ast.Expr, token lexer.Token) ast.Expr {
	return ast.ListExpr{
		Value: left,
		Next:  ParseExpr(parser, BindingPower(parser, token)),
	}
}

func ParseOrPanicExpr(parser *Parser, left ast.Expr, token lexer.Token) ast.Expr {
	funcCallExpr := left.(ast.FuncCallExpr)
	funcCallExpr.OrPanic = true
	return funcCallExpr
}

func BindingPower(parser *Parser, token lexer.Token) int {
	switch token.Type {
	case lexer.TokenOrPanic:
		fallthrough
	case lexer.TokenParenOpen:
		fallthrough
	case lexer.TokenDot:
		return 14
	case lexer.TokenStar:
		return 12
	case lexer.TokenPlus:
		return 11
	case lexer.TokenAssign:
		fallthrough
	case lexer.TokenDeclAssign:
		return 2
	case lexer.TokenComma:
		return 1
	case lexer.TokenNumber:
		fallthrough
	case lexer.TokenIdentifier:
		fallthrough
	case lexer.TokenParenClose:
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
	case lexer.TokenString:
		return ParseStringExpr(token)
	case lexer.TokenIdentifier:
		return ParseSymbolExpr(token)
	case lexer.TokenNumber:
		return ParseNumberExpr(token)
	case lexer.TokenParenOpen:
		return ParseParenOpenExpr(parser, nil, token)
	default:
		parser.InvalidToken(token)
		return nil
	}
}

func LED(parser *Parser, left ast.Expr, token lexer.Token) ast.Expr {
	switch token.Type {
	case lexer.TokenOrPanic:
		return ParseOrPanicExpr(parser, left, token)
	case lexer.TokenParenOpen:
		return ParseParenOpenExpr(parser, left, token)
	case lexer.TokenDot:
		return ParseDotExpr(parser, left, token)
	case lexer.TokenStar:
		fallthrough
	case lexer.TokenPlus:
		return ParseBinaryExpr(parser, left, token)
	case lexer.TokenAssign:
		return ParseAssignmentExpr(parser, left, token)
	case lexer.TokenDeclAssign:
		return ParseDeclAssignExpr(parser, left, token)
	case lexer.TokenComma:
		return ParseListExpr(parser, left, token)
	default:
		parser.InvalidToken(token)
		return nil
	}
}
