package ast

type Stmt interface {
	_NOP_stmt()
}

type BlockStmt struct {
	Body []Stmt
}

type PackageStmt struct {
	PackageName string
}

type ImportStmt struct {
	PackageName string
}

type FuncParameter struct {
	ParamName string
	ParamType string
}

type FuncDeclStmt struct {
	FuncName    string
	Parameters  []FuncParameter
	ReturnTypes []string
	FuncBlock   BlockStmt
}

type ExprStmt struct {
	Expr Expr
}

func (BlockStmt) _NOP_stmt()    {}
func (PackageStmt) _NOP_stmt()  {}
func (ImportStmt) _NOP_stmt()   {}
func (FuncDeclStmt) _NOP_stmt() {}
func (ExprStmt) _NOP_stmt()     {}
