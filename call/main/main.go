package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"jstree-go/call"
)

func main() {
	fset := token.NewFileSet()
	path, _ := filepath.Abs("./call/demo/srv.go")
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		log.Println(err)
		return
	}
	bv := call.NewBranchVisitor(fset)
	ast.Walk(bv, f)

	decl := ""
	decl += "package model \n" + call.Annotate("./call/demo/srv.go", "count", "Xzm", "./call/demo/model") + "\n"
	ioutil.WriteFile("./call/demo/model/dao.go", []byte(decl), os.ModePerm)

	content, err := ioutil.ReadFile("./call/demo/srv.go")
	fbv := call.NewFuncBranchVisitor(fset, content)
	fbv.Current = fbv.Root
	ast.Walk(fbv, f)
	fbv.Current = nil
	fmt.Println(fbv)
	da, err := json.Marshal(fbv)
	fmt.Println(string(da), err)

	fv := call.NewFuncVisitor(fset, content)
	ast.Walk(fv, f)
	da, err = json.Marshal(fv.Funcs)
	fmt.Println(string(da), err)
}
