package transpiler

import "github.com/tobiashort/gox/ast"

type Transpiler struct{}

func NewTranspiler() *Transpiler {
	return &Transpiler{}
}

func (transpiler *Transpiler) Transpile(ast []ast.Stmt) string {
	return ""
}
