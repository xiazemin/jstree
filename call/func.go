package call

import (
	"fmt"
	"go/ast"
	"go/token"
)

type FuncVisitor struct {
	fset    *token.FileSet
	content []byte // 解析前的
	Funcs   map[string]*BranchNode
}

func NewFuncVisitor(fset *token.FileSet, content []byte) *FuncVisitor {
	return &FuncVisitor{
		fset:    fset,
		content: content,
		Funcs:   make(map[string]*BranchNode),
	}
}

func (v *FuncVisitor) Visit(node ast.Node) ast.Visitor {
	switch no := node.(type) {
	case *ast.FuncDecl:
		fmt.Println(no.Name.Name)
		if no.Name.Name[0] >= 'A' && no.Name.Name[0] <= 'Z' {
			fmt.Println(no.Name.Name)
			fmt.Println(no.Body)
			fbv := &FuncBranchVisitor{
				fset:    v.fset,
				content: v.content,
				edit:    NewBuffer(v.content), // QINIU
				Root: &BranchNode{
					Comment: no.Name.Name,
				},
			}
			fbv.Current = fbv.Root
			ast.Walk(fbv, no)
			fbv.Root.Children = append(fbv.Root.Children, &BranchNode{
				Comment: "return",
				Return:  true,
			})
			v.Funcs[no.Name.Name] = fbv.Root
		}
	}
	return v
}
