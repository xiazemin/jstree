package call

import (
	"fmt"
	"go/ast"
	"go/token"
)

type BranchVisitor struct {
	fset  *token.FileSet
	Funcs []string
}

func NewBranchVisitor(fset *token.FileSet) *BranchVisitor {
	return &BranchVisitor{
		fset: fset,
	}
}

func (v *BranchVisitor) Visit(node ast.Node) ast.Visitor {
	switch no := node.(type) {
	case *ast.IfStmt:
		fmt.Println(no.Cond)
	case *ast.SwitchStmt:
		fmt.Println(no.Init)
	case *ast.ForStmt:
		fmt.Println(no.Cond)
	}
	return v
}
