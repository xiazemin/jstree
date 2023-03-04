package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"

	"jstree-go/call"
)

type Node struct {
	ID     string `json:"id"`
	Parent string `json:"parent"`
	State  State  `json:"state,omitempty"`
	Text   string `json:"text"`
	AAttr  Attr   `json:"a_attr,omitempty"`
}

type State struct {
	Disabled bool `json:"disabled"`
	Opened   bool `json:"opened"`
	Selected bool `json:"selected"`
}

type Attr struct {
	Style string `json:"style"`
}

func BranchNodeToNode(root map[string]*call.BranchNode) (nodes []*Node) {
	if root == nil {
		return nil
	}
	var id int64
	nodes = []*Node{{
		ID:     strconv.FormatInt(id, 10),
		Parent: "#",
		Text:   "root",
		State: State{
			Opened: true,
		},
		AAttr: Attr{
			Style: "color:red",
		},
	}}

	globalId := id
	for _, v := range root {
		// id++
		// no := &Node{
		// 	ID:     strconv.FormatInt(id, 10),
		// 	Parent: "0",
		// 	Text:   k,
		// 	State: State{
		// 		Opened: true,
		// 	},
		// 	AAttr: Attr{
		// 		Style: "blue",
		// 	},
		// }
		// nodes = append(nodes, no)
		globalId, nodes = addChildren(id, globalId, v, nodes)
	}
	return nodes
}

func addChildren(parentId, globalId int64, node *call.BranchNode, nodes []*Node) (int64, []*Node) {
	globalId++
	text := node.Comment
	if node.Code != "" {
		text += ":" + node.Code
	}
	no := &Node{
		ID:     strconv.FormatInt(globalId, 10),
		Parent: strconv.FormatInt(parentId, 10),
		State: State{
			Opened: true,
		},
		Text: text,
		AAttr: Attr{
			Style: "color:blue",
		},
	}
	if node.Return {
		no.AAttr.Style = "color:green"
	}
	nodes = append(nodes, no)
	currId := globalId
	for _, c := range node.Children {
		globalId, nodes = addChildren(currId, globalId, c, nodes)
	}
	return globalId, nodes
}

func getNodes(name string) []byte {
	fset := token.NewFileSet()
	fileName, err := filepath.Abs(name)
	fmt.Println(err, name, fileName)
	// fileName = "./call/demo/srv.go"
	path, _ := filepath.Abs(fileName)
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		log.Println(err)
		return nil
	}
	content, err := ioutil.ReadFile(fileName)
	fv := call.NewFuncVisitor(fset, content)
	ast.Walk(fv, f)
	nodes := BranchNodeToNode(fv.Funcs)
	fun, err := json.Marshal(fv.Funcs)
	fmt.Println(string(fun), err)
	da, err := json.Marshal(nodes)
	fmt.Println(string(da), err)
	return da
}
