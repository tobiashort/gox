package ast

import "github.com/tobiashort/gox/lexer"

type Stmt interface {
	_NOP_stmt()
}

type BlockStmt struct {
	Body []Stmt
}

type PackageStmt struct {
	PackageName lexer.Token
}

type ImportStmt struct {
	PackageName lexer.Token
}

type FuncParameter struct {
	Name lexer.Token
	Type lexer.Token
}

type FuncDeclStmt struct {
	Name        lexer.Token
	Parameters  []FuncParameter
	ReturnTypes []lexer.Token
	Block       BlockStmt
}

type ExprStmt struct {
	Expr Expr
}

func (BlockStmt) _NOP_stmt()    {}
func (PackageStmt) _NOP_stmt()  {}
func (ImportStmt) _NOP_stmt()   {}
func (FuncDeclStmt) _NOP_stmt() {}
func (ExprStmt) _NOP_stmt()     {}
