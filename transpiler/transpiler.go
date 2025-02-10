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
	transpiler.TranspileWithDepth(_ast, 0)
	return strings.TrimSpace(transpiler.StringBuilder.String())
}

func (transpiler *Transpiler) Write(str string) {
	transpiler.StringBuilder.WriteString(str)
}

func (transpiler *Transpiler) Writef(format string, args ...any) {
	transpiler.StringBuilder.WriteString(fmt.Sprintf(format, args...))
}

func (transpiler *Transpiler) TranspileSymbolExpr(expr ast.SymbolExpr) {
	transpiler.Write(expr.Symbol.Value)
}

func (transpiler *Transpiler) TranspileStringExpr(expr ast.StringExpr) {
	transpiler.Writef("\"%s\"", expr.String.Value)
}

func (transpiler *Transpiler) TranspileAccessExpr(stmtInterface ast.Stmt, expr ast.AccessExpr, indent string) {
	transpiler.TranspileExpr(stmtInterface, expr.Instance, indent)
	transpiler.Write(".")
	transpiler.TranspileExpr(stmtInterface, expr.Field, indent)
}

func (transpiler *Transpiler) TranspileFuncCallExpr(stmtInterface ast.Stmt, expr ast.FuncCallExpr, indent string) {
	switch stmtInterface.(type) {
	case ast.ReturnStmt:
		transpiler.TranspileExpr(stmtInterface, expr.Func, indent)
		transpiler.Write("(")
		transpiler.TranspileExpr(stmtInterface, expr.Args, indent)
		transpiler.Write(")\n")
		if expr.OrPanic {
			panic(fmt.Sprintf("\n%s... <--- illegal or_panic", transpiler.StringBuilder.String()))
		}
	case ast.ExprStmt:
		transpiler.Write(indent)
		if expr.OrPanic {
			transpiler.StringBuilder.WriteString("err := ")
		}
		transpiler.TranspileExpr(stmtInterface, expr.Func, indent)
		transpiler.Write("(")
		transpiler.TranspileExpr(stmtInterface, expr.Args, indent)
		transpiler.Write(")\n")
		if expr.OrPanic {
			transpiler.Writef("%sif err != nil { \n", indent)
			transpiler.Writef("%s%spanic(err)\n", indent, indent)
			transpiler.Writef("%s}\n", indent)
		}
	default:
		panic(fmt.Sprintf("\n%s... <--- unhandled %s", transpiler.StringBuilder.String(), reflect.TypeOf(expr)))
	}
}

func (transpiler *Transpiler) TranspileAssignDeclExpr(stmtInterface ast.Stmt, expr ast.AssignDeclExpr, indent string) {
	transpiler.Write(indent)
	transpiler.TranspileExpr(stmtInterface, expr.Left, indent)
	transpiler.Write(" := ")
	transpiler.TranspileExpr(stmtInterface, expr.Right, indent)
	transpiler.Write("\n")
}

func (transpiler *Transpiler) TranspileListExpr(stmtInterface ast.Stmt, expr ast.ListExpr, indent string) {
	transpiler.TranspileExpr(stmtInterface, expr.Value, indent)
	if expr.Next != nil {
		transpiler.Write(", ")
		transpiler.TranspileExpr(stmtInterface, expr.Next, "")
	}
	transpiler.Write("\n")
}

func (transpiler *Transpiler) TranspileExpr(stmtInterface ast.Stmt, exprInterface ast.Expr, indent string) {
	switch expr := exprInterface.(type) {
	case nil:
		return
	case ast.SymbolExpr:
		transpiler.TranspileSymbolExpr(expr)
	case ast.StringExpr:
		transpiler.TranspileStringExpr(expr)
	case ast.AccessExpr:
		transpiler.TranspileAccessExpr(stmtInterface, expr, indent)
	case ast.FuncCallExpr:
		transpiler.TranspileFuncCallExpr(stmtInterface, expr, indent)
	case ast.AssignDeclExpr:
		transpiler.TranspileAssignDeclExpr(stmtInterface, expr, indent)
	case ast.ListExpr:
		transpiler.TranspileListExpr(stmtInterface, expr, indent)
	default:
		panic(fmt.Sprintf("\n%s... <--- unhandled %s", transpiler.StringBuilder.String(), reflect.TypeOf(expr)))
	}
}

func (transpiler *Transpiler) TranspilePackageStmt(stmt ast.PackageStmt, indent string) {
	transpiler.Writef("%spackage %s\n\n", indent, stmt.PackageName.Value)
}

func (transpiler *Transpiler) TranspileImportStmt(stmt ast.ImportStmt, indent string) {
	transpiler.Writef("%simport \"%s\"\n\n", indent, stmt.PackageName.Value)
}

func (transpiler *Transpiler) TranspileFuncDeclStmt(stmt ast.FuncDeclStmt, indent string, depth int) {
	transpiler.Writef("%sfunc %s", indent, stmt.Name.Value)
	paramNameAndType := make([]string, 0)
	for _, param := range stmt.Parameters {
		paramNameAndType = append(paramNameAndType, fmt.Sprintf("%s %s", param.Name.Value, param.Type.Value))
	}
	transpiler.Writef("(%s) ", strings.Join(paramNameAndType, ", "))
	if len(stmt.ReturnTypes) > 0 {
		returnTypes := make([]string, 0)
		for _, _type := range stmt.ReturnTypes {
			returnTypes = append(returnTypes, _type.Value)
		}
		if len(returnTypes) > 1 {
			transpiler.Writef("(%s) ", strings.Join(returnTypes, ", "))
		} else {
			transpiler.Writef("%s ", strings.Join(returnTypes, ", "))
		}
	}
	transpiler.Write("{\n")
	transpiler.TranspileWithDepth(stmt.Block.(ast.BlockStmt).Body, depth+1)
	transpiler.Write("}\n\n")
}

func (transpiler *Transpiler) TranspileReturnStmt(stmt ast.ReturnStmt, indent string) {
	transpiler.Writef("%sreturn ", indent)
	transpiler.TranspileExpr(stmt, stmt.Values, indent)
}

func (transpiler *Transpiler) TranspileVarDeclStmt(stmt ast.VarDeclStmt, indent string) {
	transpiler.Writef("%svar %s %s\n", indent, stmt.Name.Value, stmt.Type.Value)
}

func (transpiler *Transpiler) TranspileWithDepth(_ast []ast.Stmt, depth int) {
	indent := strings.Repeat("\t", depth)
	for _, stmtInterface := range _ast {
		switch stmt := stmtInterface.(type) {
		case ast.PackageStmt:
			transpiler.TranspilePackageStmt(stmt, indent)
		case ast.ImportStmt:
			transpiler.TranspileImportStmt(stmt, indent)
		case ast.FuncDeclStmt:
			transpiler.TranspileFuncDeclStmt(stmt, indent, depth)
		case ast.ReturnStmt:
			transpiler.TranspileReturnStmt(stmt, indent)
		case ast.VarDeclStmt:
			transpiler.TranspileVarDeclStmt(stmt, indent)
		case ast.ExprStmt:
			transpiler.TranspileExpr(stmt, stmt.Expr, indent)
		default:
			panic(fmt.Sprintf("\n%s%s--- here\n%sunhandled %s", transpiler.StringBuilder.String(), indent, indent, reflect.TypeOf(stmt)))
		}
	}
}
