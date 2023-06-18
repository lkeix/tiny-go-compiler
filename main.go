package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "./testdata/testdata.go", nil, 0)
	if err != nil {
		panic(err)
	}

	// TODO: emit global variables
	fmt.Printf(".data\n")
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			emitBody(d.Body)
			fmt.Printf("\n")
			if d.Name.Name == "main" {
				fmt.Printf(".text\n")
				emitFunc("main.main", d.Body)
			}
		}
	}

	fmt.Printf(".global _start\n")
	fmt.Printf("_start:\n")
	fmt.Printf("\tcallq main.main\n")
	// emit Exit
	fmt.Printf("\tmovq %%rax, %%rdi\n")
	fmt.Printf("\tmovq $60, %%rax\n")
	fmt.Printf("\tsyscall\n")

}

func emitFunc(fname string, body *ast.BlockStmt) {
	fmt.Printf("%s:\n", fname)
	for _, stmt := range body.List {
		switch s := stmt.(type) {
		case *ast.ExprStmt:
			emitExpr(s.X)
		}
	}
	fmt.Printf("\tpushq %%rax\n")
	fmt.Printf("\tpopq %%rax\n")
	fmt.Printf("\tret\n")
}

func emitExpr(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		fmt.Printf("\tmovq %s+0(%%rip), %%rax\n", e.X.(*ast.Ident).Name)
		fmt.Printf("\tmovq %s+0(%%rip), %%rdi\n", e.Y.(*ast.Ident).Name)
		if e.Op == token.ADD {
			fmt.Printf("\taddq %%rdi, %%rax\n")
		}
		if e.Op == token.SUB {
			fmt.Printf("\tsubq %%rdi, %%rax\n")
		}
	}
}

func emitBody(body *ast.BlockStmt) {
	for _, stmt := range body.List {
		switch s := stmt.(type) {
		case *ast.AssignStmt:
			emitAssign(s.Lhs, s.Rhs)
		}
	}
}

func emitAssign(lhs, rhs []ast.Expr) {
	for i, l := range lhs {
		switch e := l.(type) {
		case *ast.Ident:
			rhs := rhs[i].(*ast.BasicLit)
			fmt.Printf("%s:\n", e.Name)
			fmt.Printf("\t.quad %s\n", rhs.Value)
		}
	}
}
