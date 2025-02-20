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
	locals := make(map[string]bool)
	transpiler.TranspileWithDepth(parser.Stmts, 0, locals)
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

func (transpiler *Transpiler) TranspileAccessExpr(stmtInterface ast.Stmt, expr ast.AccessExpr, indent string, locals map[string]bool) {
	transpiler.TranspileExpr(stmtInterface, expr.Instance, indent, locals)
	transpiler.Write(".")
	transpiler.TranspileExpr(stmtInterface, expr.Field, indent, locals)
}

func (transpiler *Transpiler) TranspileBinaryExpr(stmtInterface ast.Stmt, expr ast.BinaryExpr, indent string, locals map[string]bool) {
	transpiler.TranspileExpr(stmtInterface, expr.Left, indent, locals)
	switch expr.Operator.Type {
	case lexer.TokenPlus:
		transpiler.Write(" + ")
	case lexer.TokenStar:
		transpiler.Write(" * ")
	default:
		panic(fmt.Sprintf("\n%s... <--- unhandled %s", transpiler.String(), reflect.TypeOf(expr.Operator)))
	}
	transpiler.TranspileExpr(stmtInterface, expr.Right, indent, locals)
}

func (transpiler *Transpiler) TranspileFuncCallExpr(stmtInterface ast.Stmt, expr ast.FuncCallExpr, indent string, locals map[string]bool) {
	switch stmt := stmtInterface.(type) {
	case ast.ReturnStmt:
		if expr.OrPanic {
			retExists := locals["ret"]
			if retExists {
				transpiler.Writef("%sret, err = ", indent)
			} else {
				transpiler.Writef("%sret, err := ", indent)
				locals["ret"] = true
				locals["err"] = true
			}
			transpiler.TranspileExpr(stmtInterface, expr.Func, indent, locals)
			transpiler.Write("(")
			transpiler.TranspileExpr(stmtInterface, expr.Args, indent, locals)
			transpiler.Write(")\n")
			transpiler.Writef("%sif err != nil {\n", indent)
			transpiler.Writef("%s%spanic(err)\n", indent, indent)
			transpiler.Writef("%s}\n", indent)
			transpiler.Writef("%sreturn ret", indent)
		} else {
			transpiler.TranspileExpr(stmtInterface, expr.Func, indent, locals)
			transpiler.Write("(")
			transpiler.TranspileExpr(stmtInterface, expr.Args, indent, locals)
			transpiler.Write(")")
		}
	case ast.ExprStmt:
		switch stmt.Expr.(type) {
		case ast.DeclAssignExpr:
			transpiler.TranspileExpr(stmtInterface, expr.Func, indent, locals)
			transpiler.Write("(")
			transpiler.TranspileExpr(stmtInterface, expr.Args, indent, locals)
			transpiler.Write(")\n")
			if expr.OrPanic {
				transpiler.Writef("%sif err != nil { \n", indent)
				transpiler.Writef("%s%spanic(err)\n", indent, indent)
				transpiler.Writef("%s}\n", indent)
			}
		case ast.FuncCallExpr:
			transpiler.Write(indent)
			if expr.OrPanic {
				errExists := locals["err"]
				if errExists {
					transpiler.Write("err = ")
				} else {
					transpiler.Write("err := ")
					locals["err"] = true
				}
			}
			transpiler.TranspileExpr(stmtInterface, expr.Func, indent, locals)
			transpiler.Write("(")
			transpiler.TranspileExpr(stmtInterface, expr.Args, indent, locals)
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

func (transpiler *Transpiler) TranspileDeclAssignExpr(stmtInterface ast.Stmt, expr ast.DeclAssignExpr, indent string, locals map[string]bool) {
	transpiler.Write(indent)
	transpiler.TranspileExpr(stmtInterface, expr.Left, indent, locals)
	funcCallExpr, isFuncCallExpr := expr.Right.(ast.FuncCallExpr)
	if isFuncCallExpr && funcCallExpr.OrPanic {
		transpiler.Write(", err")
		locals["err"] = true
	}
	transpiler.Write(" := ")
	transpiler.TranspileExpr(stmtInterface, expr.Right, indent, locals)
}

func (transpiler *Transpiler) TranspileAssignmentExpr(stmtInterface ast.Stmt, expr ast.AssignmentExpr, indent string, locals map[string]bool) {
	transpiler.Write(indent)
	transpiler.TranspileExpr(stmtInterface, expr.Left, indent, locals)
	transpiler.Write(" = ")
	transpiler.TranspileExpr(stmtInterface, expr.Right, indent, locals)
	transpiler.Write("\n")
}

func (transpiler *Transpiler) TranspileListExpr(stmtInterface ast.Stmt, expr ast.ListExpr, indent string, locals map[string]bool) {
	transpiler.TranspileExpr(stmtInterface, expr.Value, indent, locals)
	if expr.Next != nil {
		transpiler.Write(", ")
		transpiler.TranspileExpr(stmtInterface, expr.Next, "", locals)
	}
}

func (transpiler *Transpiler) TranspileExpr(stmtInterface ast.Stmt, exprInterface ast.Expr, indent string, locals map[string]bool) {
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
		transpiler.TranspileAccessExpr(stmtInterface, expr, indent, locals)
	case ast.BinaryExpr:
		transpiler.TranspileBinaryExpr(stmtInterface, expr, indent, locals)
	case ast.FuncCallExpr:
		transpiler.TranspileFuncCallExpr(stmtInterface, expr, indent, locals)
	case ast.AssignmentExpr:
		transpiler.TranspileAssignmentExpr(stmtInterface, expr, indent, locals)
	case ast.DeclAssignExpr:
		transpiler.TranspileDeclAssignExpr(stmtInterface, expr, indent, locals)
	case ast.ListExpr:
		transpiler.TranspileListExpr(stmtInterface, expr, indent, locals)
	default:
		panic(fmt.Sprintf("\n%s... <--- unhandled %s", transpiler.StringBuilder.String(), reflect.TypeOf(expr)))
	}
}

func (transpiler *Transpiler) TranspilePackageStmt(stmt ast.PackageStmt, indent string) {
	transpiler.Writef("%spackage %s\n\n", indent, stmt.PackageName.Value)
}

func (transpiler *Transpiler) TranspileImportStmt(stmt ast.ImportStmt, indent string) {
	transpiler.Writef("%simport ", indent)
	if len(stmt.PackageNames) == 0 {
		panic(fmt.Sprintf("\n%s... <--- ", transpiler.StringBuilder.String()))
	} else if len(stmt.PackageNames) == 1 {
		transpiler.Writef("\"%s\"\n\n", stmt.PackageNames[0].Value)
	} else {
		transpiler.Write("\n")
		for _, packageName := range stmt.PackageNames {
			transpiler.Writef("%s%s\"%s\"\n", indent, indent, packageName.Value)
		}
		transpiler.Writef("\n%s)\n\n", indent)
	}
}

func (transpiler *Transpiler) TranspileFuncDeclStmt(stmt ast.FuncDeclStmt, indent string, depth int, locals map[string]bool) {
	locals[stmt.Name.Value] = true
	innerLocals := make(map[string]bool)
	for k, v := range locals {
		innerLocals[k] = v
	}
	transpiler.Writef("%sfunc %s", indent, stmt.Name.Value)
	paramNameAndType := make([]string, 0)
	for _, param := range stmt.Parameters {
		innerLocals[param.Name.Value] = true
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
	transpiler.TranspileWithDepth(stmt.Block.(ast.BlockStmt).Body, depth+1, innerLocals)
	transpiler.Write("}\n\n")
}

func (transpiler *Transpiler) TranspileReturnStmt(stmt ast.ReturnStmt, indent string, locals map[string]bool) {
	funcCallExpr, isFuncCallExpr := stmt.Values.(ast.FuncCallExpr)
	if isFuncCallExpr && funcCallExpr.OrPanic {
		// defer the return keyword
	} else {
		transpiler.Writef("%sreturn ", indent)
	}
	transpiler.TranspileExpr(stmt, stmt.Values, indent, locals)
	transpiler.Write("\n")
}

func (transpiler *Transpiler) TranspileVarDeclStmt(stmt ast.VarDeclStmt, indent string, locals map[string]bool) {
	locals[stmt.Name.Value] = true
	transpiler.Writef("%svar %s %s\n", indent, stmt.Name.Value, stmt.Type.Value)
}

func (transpiler *Transpiler) TranspileWithDepth(_ast []ast.Stmt, depth int, locals map[string]bool) {
	indent := strings.Repeat("\t", depth)
	for _, stmtInterface := range _ast {
		switch stmt := stmtInterface.(type) {
		case ast.PackageStmt:
			transpiler.TranspilePackageStmt(stmt, indent)
		case ast.ImportStmt:
			transpiler.TranspileImportStmt(stmt, indent)
		case ast.FuncDeclStmt:
			transpiler.TranspileFuncDeclStmt(stmt, indent, depth, locals)
		case ast.ReturnStmt:
			transpiler.TranspileReturnStmt(stmt, indent, locals)
		case ast.VarDeclStmt:
			transpiler.TranspileVarDeclStmt(stmt, indent, locals)
		case ast.ExprStmt:
			transpiler.TranspileExpr(stmt, stmt.Expr, indent, locals)
		default:
			panic(fmt.Sprintf("\n%s%s--- here\n%sunhandled %s", transpiler.StringBuilder.String(), indent, indent, reflect.TypeOf(stmt)))
		}
	}
}
