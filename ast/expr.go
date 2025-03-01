package ast

import "github.com/tobiashort/gox/lexer"

type Expr interface {
	_NOP_expr()
}

type SymbolExpr struct {
	Symbol lexer.Token
}

type NumberExpr struct {
	Number lexer.Token
}

type StringExpr struct {
	String lexer.Token
}

type AssignmentExpr struct {
	Left  Expr
	Right Expr
}

type DeclAssignExpr struct {
	Left  Expr
	Right Expr
}

type BinaryExpr struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

type AccessExpr struct {
	Instance Expr
	Field    Expr
}

type FuncCallExpr struct {
	Func    Expr
	Args    Expr
	OrPanic bool
}

type ListExpr struct {
	Value Expr
	Next  Expr
}

func (SymbolExpr) _NOP_expr()     {}
func (NumberExpr) _NOP_expr()     {}
func (StringExpr) _NOP_expr()     {}
func (AssignmentExpr) _NOP_expr() {}
func (DeclAssignExpr) _NOP_expr() {}
func (BinaryExpr) _NOP_expr()     {}
func (AccessExpr) _NOP_expr()     {}
func (FuncCallExpr) _NOP_expr()   {}
func (ListExpr) _NOP_expr()       {}
