package ast

import "github.com/tobiashort/gox/lexer"

type Expr interface {
	_NOP_expr()
}

type SymbolExpr struct {
	Symbol lexer.Token
}

type NumberExpr struct {
	Number string
}

type AssignmentExpr struct {
	Left  Expr
	Right Expr
}

type BinaryExpr struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

func (SymbolExpr) _NOP_expr()     {}
func (NumberExpr) _NOP_expr()     {}
func (AssignmentExpr) _NOP_expr() {}
func (BinaryExpr) _NOP_expr()     {}
