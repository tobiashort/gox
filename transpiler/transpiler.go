package transpiler

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tobiashort/gox/ast"
	"github.com/tobiashort/gox/lexer"
	"github.com/tobiashort/gox/parser"
)

type Transpiler struct {
	StringBuilder strings.Builder
}

func NewTranspiler() *Transpiler {
	return &Transpiler{
		StringBuilder: strings.Builder{},
	}
}

func (transpiler *Transpiler) Transpile(source string) string {
	parser := parser.NewParser()
	parser.Parse(source)
	transpiler.TranspileWithDepth(parser.Stmts, 0)
	return strings.TrimSpace(transpiler.StringBuilder.String())
}

func (transpiler *Transpiler) Write(str string) {
	transpiler.StringBuilder.WriteString(str)
}

func (transpiler *Transpiler) Writef(format string, args ...any) {
	transpiler.StringBuilder.WriteString(fmt.Sprintf(format, args...))
}

func (transpiler *Transpiler) String() string {
	return transpiler.StringBuilder.String()
}

func (transpiler *Transpiler) TranspileSymbolExpr(expr ast.SymbolExpr) {
	transpiler.Write(expr.Symbol.Value)
}

func (transpiler *Transpiler) TranspileStringExpr(expr ast.StringExpr) {
	transpiler.Writef("\"%s\"", expr.String.Value)
}

func (transpiler *Transpiler) TranspileNumberExpr(expr ast.NumberExpr) {
	transpiler.Write(expr.Number.Value)
}

func (transpiler *Transpiler) TranspileAccessExpr(stmtInterface ast.Stmt, expr ast.AccessExpr, indent string) {
	transpiler.TranspileExpr(stmtInterface, expr.Instance, indent)
	transpiler.Write(".")
	transpiler.TranspileExpr(stmtInterface, expr.Field, indent)
}

func (transpiler *Transpiler) TranspileBinaryExpr(stmtInterface ast.Stmt, expr ast.BinaryExpr, indent string) {
	transpiler.TranspileExpr(stmtInterface, expr.Left, indent)
	switch expr.Operator.Type {
	case lexer.TokenPlus:
		transpiler.Write(" + ")
	case lexer.TokenStar:
		transpiler.Write(" * ")
	default:
		panic(fmt.Sprintf("\n%s... <--- unhandled %s", transpiler.String(), reflect.TypeOf(expr.Operator)))
	}
	transpiler.TranspileExpr(stmtInterface, expr.Right, indent)
}

func (transpiler *Transpiler) TranspileFuncCallExpr(stmtInterface ast.Stmt, expr ast.FuncCallExpr, indent string) {
	switch stmt := stmtInterface.(type) {
	case ast.ReturnStmt:
		if expr.OrPanic {
			transpiler.Writef("%sret, err := ", indent)
			transpiler.TranspileExpr(stmtInterface, expr.Func, indent)
			transpiler.Write("(")
			transpiler.TranspileExpr(stmtInterface, expr.Args, indent)
			transpiler.Write(")\n")
			transpiler.Writef("%sif err != nil {\n", indent)
			transpiler.Writef("%s%spanic(err)\n", indent, indent)
			transpiler.Writef("%s}\n", indent)
			transpiler.Writef("%sreturn ret\n", indent)
		} else {
			transpiler.TranspileExpr(stmtInterface, expr.Func, indent)
			transpiler.Write("(")
			transpiler.TranspileExpr(stmtInterface, expr.Args, indent)
			transpiler.Write(")\n")
		}
	case ast.ExprStmt:
		switch stmt.Expr.(type) {
		case ast.DeclAssignExpr:
			transpiler.TranspileExpr(stmtInterface, expr.Func, indent)
			transpiler.Write("(")
			transpiler.TranspileExpr(stmtInterface, expr.Args, indent)
			transpiler.Write(")\n")
			if expr.OrPanic {
				transpiler.Writef("%sif err != nil { \n", indent)
				transpiler.Writef("%s%spanic(err)\n", indent, indent)
				transpiler.Writef("%s}\n", indent)
			}
		case ast.FuncCallExpr:
			transpiler.Write(indent)
			if expr.OrPanic {
				transpiler.Write("err := ")
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
			panic(fmt.Sprintf("\n%s... <--- unhandled %s", transpiler.String(), reflect.TypeOf(expr)))
		}
	default:
		panic(fmt.Sprintf("\n%s... <--- unhandled %s", transpiler.String(), reflect.TypeOf(stmt)))
	}
}

func (transpiler *Transpiler) TranspileDeclAssignExpr(stmtInterface ast.Stmt, expr ast.DeclAssignExpr, indent string) {
	transpiler.Write(indent)
	transpiler.TranspileExpr(stmtInterface, expr.Left, indent)
	funcCallExpr, isFuncCallExpr := expr.Right.(ast.FuncCallExpr)
	if isFuncCallExpr && funcCallExpr.OrPanic {
		transpiler.Write(", err")
	}
	transpiler.Write(" := ")
	transpiler.TranspileExpr(stmtInterface, expr.Right, indent)
	transpiler.Write("\n")
}

func (transpiler *Transpiler) TranspileAssignmentExpr(stmtInterface ast.Stmt, expr ast.AssignmentExpr, indent string) {
	transpiler.Write(indent)
	transpiler.TranspileExpr(stmtInterface, expr.Left, indent)
	transpiler.Write(" = ")
	transpiler.TranspileExpr(stmtInterface, expr.Right, indent)
	transpiler.Write("\n")
}

func (transpiler *Transpiler) TranspileListExpr(stmtInterface ast.Stmt, expr ast.ListExpr, indent string) {
	transpiler.TranspileExpr(stmtInterface, expr.Value, indent)
	if expr.Next != nil {
		transpiler.Write(", ")
		transpiler.TranspileExpr(stmtInterface, expr.Next, "")
	}
}

func (transpiler *Transpiler) TranspileExpr(stmtInterface ast.Stmt, exprInterface ast.Expr, indent string) {
	switch expr := exprInterface.(type) {
	case nil:
		return
	case ast.SymbolExpr:
		transpiler.TranspileSymbolExpr(expr)
	case ast.StringExpr:
		transpiler.TranspileStringExpr(expr)
	case ast.NumberExpr:
		transpiler.TranspileNumberExpr(expr)
	case ast.AccessExpr:
		transpiler.TranspileAccessExpr(stmtInterface, expr, indent)
	case ast.BinaryExpr:
		transpiler.TranspileBinaryExpr(stmtInterface, expr, indent)
	case ast.FuncCallExpr:
		transpiler.TranspileFuncCallExpr(stmtInterface, expr, indent)
	case ast.AssignmentExpr:
		transpiler.TranspileAssignmentExpr(stmtInterface, expr, indent)
	case ast.DeclAssignExpr:
		transpiler.TranspileDeclAssignExpr(stmtInterface, expr, indent)
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
	funcCallExpr, isFuncCallExpr := stmt.Values.(ast.FuncCallExpr)
	if isFuncCallExpr && funcCallExpr.OrPanic {
		// defer the return keyword
	} else {
		transpiler.Writef("%sreturn ", indent)
	}
	transpiler.TranspileExpr(stmt, stmt.Values, indent)
	transpiler.Write("\n")
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
