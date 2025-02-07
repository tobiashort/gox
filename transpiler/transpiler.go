package transpiler

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tobiashort/gox/ast"
)

type Transpiler struct {
	StringBuilder strings.Builder
}

func NewTranspiler() *Transpiler {
	return &Transpiler{
		StringBuilder: strings.Builder{},
	}
}

func (transpiler *Transpiler) Transpile(_ast []ast.Stmt) string {
	transpiler.TranspileIndent(_ast, 0)
	return strings.TrimSpace(transpiler.StringBuilder.String())
}

func (transpiler *Transpiler) TranspileIndent(_ast []ast.Stmt, depth int) {
	indent := strings.Repeat("\t", depth)
	for _, stmtInterface := range _ast {
		switch stmt := stmtInterface.(type) {
		case ast.PackageStmt:
			transpiler.StringBuilder.WriteString(fmt.Sprintf("%spackage %s\n\n", indent, stmt.PackageName.Value))
		case ast.ImportStmt:
			transpiler.StringBuilder.WriteString(fmt.Sprintf("%simport \"%s\"\n\n", indent, stmt.PackageName.Value))
		case ast.FuncDeclStmt:
			transpiler.StringBuilder.WriteString(fmt.Sprintf("%sfunc %s", indent, stmt.Name.Value))
			paramNameAndType := make([]string, 0)
			for _, param := range stmt.Parameters {
				paramNameAndType = append(paramNameAndType, fmt.Sprintf("%s %s", param.Name.Value, param.Type.Value))
			}
			transpiler.StringBuilder.WriteString(fmt.Sprintf("(%s) ", strings.Join(paramNameAndType, ", ")))
			if len(stmt.ReturnTypes) > 0 {
				returnTypes := make([]string, 0)
				for _, _type := range stmt.ReturnTypes {
					returnTypes = append(returnTypes, _type.Value)
				}
				if len(returnTypes) > 1 {
					transpiler.StringBuilder.WriteString(fmt.Sprintf("(%s) ", strings.Join(returnTypes, ", ")))
				} else {
					transpiler.StringBuilder.WriteString(fmt.Sprintf("%s ", strings.Join(returnTypes, ", ")))
				}
			}
			transpiler.StringBuilder.WriteString("{\n")
			transpiler.TranspileIndent(stmt.Block.(ast.BlockStmt).Body, depth+1)
			transpiler.StringBuilder.WriteString("}\n\n")
		case ast.ReturnStmt:
			transpiler.StringBuilder.WriteString(fmt.Sprintf("%sreturn ", indent))
			transpiler.TranspileExpr(stmt, stmt.Values, indent)
		case ast.ExprStmt:
			transpiler.TranspileExpr(stmt, stmt.Expr, indent)
		default:
			panic(fmt.Sprintf("\n%s%s--- here\n%sunhandled %s", transpiler.StringBuilder.String(), indent, indent, reflect.TypeOf(stmt)))
		}
	}
}

func (transpiler *Transpiler) TranspileExpr(stmtInterface ast.Stmt, exprInterface ast.Expr, indent string) {
	switch expr := exprInterface.(type) {
	case nil:
		return
	case ast.SymbolExpr:
		transpiler.StringBuilder.WriteString(expr.Symbol.Value)
	case ast.StringExpr:
		transpiler.StringBuilder.WriteString(fmt.Sprintf("\"%s\"", expr.String.Value))
	case ast.AccessExpr:
		transpiler.TranspileExpr(stmtInterface, expr.Instance, indent)
		transpiler.StringBuilder.WriteString(".")
		transpiler.TranspileExpr(stmtInterface, expr.Field, indent)
	case ast.FuncCallExpr:
		switch stmtInterface.(type) {
		case ast.ReturnStmt:
			transpiler.TranspileExpr(stmtInterface, expr.Func, indent)
			transpiler.StringBuilder.WriteString("(")
			transpiler.TranspileExpr(stmtInterface, expr.Args, indent)
			transpiler.StringBuilder.WriteString(")\n")
			if expr.OrPanic {
				panic(fmt.Sprintf("\n%s... <--- illegal or_panic", transpiler.StringBuilder.String()))
			}
		case ast.ExprStmt:
			transpiler.StringBuilder.WriteString(indent)
			if expr.OrPanic {
				transpiler.StringBuilder.WriteString("err = ")
			}
			transpiler.TranspileExpr(stmtInterface, expr.Func, indent)
			transpiler.StringBuilder.WriteString("(")
			transpiler.TranspileExpr(stmtInterface, expr.Args, indent)
			transpiler.StringBuilder.WriteString(")\n")
			if expr.OrPanic {
				transpiler.StringBuilder.WriteString(fmt.Sprintf("%sif err != nil { \n", indent))
				transpiler.StringBuilder.WriteString(fmt.Sprintf("%s%spanic(err)\n", indent, indent))
				transpiler.StringBuilder.WriteString(fmt.Sprintf("%s}\n", indent))
			}
		default:
			panic(fmt.Sprintf("\n%s... <--- unhandled %s", transpiler.StringBuilder.String(), reflect.TypeOf(expr)))
		}
	default:
		panic(fmt.Sprintf("\n%s... <--- unhandled %s", transpiler.StringBuilder.String(), reflect.TypeOf(expr)))
	}
}
