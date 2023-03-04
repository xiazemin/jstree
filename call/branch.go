package call

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
)

type FuncBranchVisitor struct {
	fset    *token.FileSet //解析后的
	content []byte         // 解析前的
	edit    *Buffer        // QINIU
	Root    *BranchNode
	Current *BranchNode
}

func NewFuncBranchVisitor(fset *token.FileSet, content []byte) *FuncBranchVisitor {
	return &FuncBranchVisitor{
		fset:    fset,
		content: content,
		edit:    NewBuffer(content), // QINIU
		Root: &BranchNode{
			Comment: "root",
		},
	}
}

func (f *FuncBranchVisitor) String() string {
	return f.edit.String()
}

// Visit implements the ast.Visitor interface.
func (f *FuncBranchVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	// case *ast.FuncDecl:
	// 	cur := f.addNode(int(n.Name.Pos()-1), int(n.Name.End()), n.Name.Name)
	// 	ast.Walk(f, n.Body)
	// 	f.Current = cur
	// 	return nil
	case *ast.BlockStmt:
		// If it's a switch or select, the body is a list of case clauses; don't tag the block itself.
		if len(n.List) > 0 {
			switch n.List[0].(type) {
			case *ast.CaseClause: // switch
				for _, n := range n.List {
					clause := n.(*ast.CaseClause)
					f.addCounters(clause.Colon+1, clause.Colon+1, clause.End(), clause.Body, false)
				}
				return f
			case *ast.CommClause: // select
				for _, n := range n.List {
					clause := n.(*ast.CommClause)
					f.addCounters(clause.Colon+1, clause.Colon+1, clause.End(), clause.Body, false)
				}
				return f
			}
		}
		f.addCounters(n.Lbrace, n.Lbrace+1, n.Rbrace+1, n.List, true) // +1 to step past closing brace.
	case *ast.IfStmt:
		// cur := f.addNode(int(n.Pos()-1), int(n.End()))
		// defer func() {
		// 	f.Current = cur
		// }()
		if n.Init != nil {
			// cur := f.addNode(int(n.Pos()-1), int(n.End()))
			ast.Walk(f, n.Init)
			// f.Current = cur
		}
		ast.Walk(f, n.Cond)
		cur := f.Current
		if f.Current.Else {
			cur = f.addNode(n.Init, int(n.Cond.Pos()-1), int(n.Cond.End()), "else if")
		} else {
			cur = f.addNode(n.Init, int(n.Cond.Pos()-1), int(n.Cond.End()), "if")
		}
		ast.Walk(f, n.Body)
		// f.Current = cur
		if n.Else == nil {
			f.Current = cur
			return nil
		}
		// The elses are special, because if we have
		//	if x {
		//	} else if y {
		//	}
		// we want to cover the "if y". To do this, we need a place to drop the counter,
		// so we add a hidden block:
		//	if x {
		//	} else {
		//		if y {
		//		}
		//	}
		elseOffset := f.findText(n.Body.End(), "else")
		if elseOffset < 0 {
			panic("lost else")
		}
		f.edit.Insert(elseOffset+4, "{")
		f.edit.Insert(f.offset(n.Else.End()), "}")

		// We just created a block, now walk it.
		// Adjust the position of the new block to start after
		// the "else". That will cause it to follow the "{"
		// we inserted above.
		pos := f.fset.File(n.Body.End()).Pos(elseOffset + 4)
		switch stmt := n.Else.(type) {
		case *ast.IfStmt:
			block := &ast.BlockStmt{
				Lbrace: pos,
				List:   []ast.Stmt{stmt},
				Rbrace: stmt.End(),
			}
			n.Else = block
			f.Current.Else = true
			// cur = f.addNode(int(stmt.Cond.Pos()-1), int(stmt.Cond.End()), "else")
			// f.Current = cur
		case *ast.BlockStmt:
			stmt.Lbrace = pos
			cur = f.addNode(n.Init, int(n.Else.End()), int(n.Else.End()), "else")
		default:
			panic("unexpected node type in if")
		}
		// cur = f.addNode(int(n.Else.Pos()-1), int(n.Else.End()), "else")
		ast.Walk(f, n.Else)
		f.Current = cur
		// f.Current = cur
		return nil
	case *ast.SelectStmt:
		// Don't annotate an empty select - creates a syntax error.
		if n.Body == nil || len(n.Body.List) == 0 {
			return nil
		}
	case *ast.SwitchStmt:
		// Don't annotate an empty switch - creates a syntax error.
		if n.Body == nil || len(n.Body.List) == 0 {
			if n.Init != nil {
				// cur := f.addNode(int(n.Pos()-1), int(n.End()))
				ast.Walk(f, n.Init)
				// f.Current = cur
			}
			if n.Tag != nil {
				// cur := f.addNode(int(n.Pos()-1), int(n.End()))
				ast.Walk(f, n.Tag)
				// f.Current = cur
			}
			return nil
		}
	case *ast.TypeSwitchStmt:
		// Don't annotate an empty type switch - creates a syntax error.
		if n.Body == nil || len(n.Body.List) == 0 {
			if n.Init != nil {
				// cur := f.addNode(int(n.Pos()-1), int(n.End()))
				ast.Walk(f, n.Init)
				// f.Current = cur
			}
			// cur := f.addNode(int(n.Pos()-1), int(n.End()))
			ast.Walk(f, n.Assign)
			// f.Current = cur
			return nil
		}
	}
	return f
}

func (f *FuncBranchVisitor) addNode(ini ast.Stmt, start, end int, name string) *BranchNode {
	prefix := string(f.content[start:end])
	if ini != nil {
		prefix = string(f.content[ini.Pos()-1:ini.End()]) + ";" + string(f.content[start:end])
	}
	cur := &BranchNode{
		Comment: name,
		Code:    prefix,
	}
	oldCur := f.Current
	f.Current.AddChild(cur)
	f.Current = cur
	return oldCur
}

func (f *FuncBranchVisitor) addCounters(pos, insertPos, blockEnd token.Pos, list []ast.Stmt, extendToClosingBrace bool) {
	fmt.Println(pos)
	for _, s := range list {
		switch no := s.(type) {
		case *ast.ReturnStmt:
			fmt.Println("find ruturn:", no.Return, no.Results)
			if f.Current.Return {
				fmt.Println("return twice")
			}
			f.Current.Return = true
		}
		fmt.Println(s)
	}
}
func (f *FuncBranchVisitor) findText(pos token.Pos, text string) int {
	b := []byte(text)
	start := f.offset(pos)
	i := start
	s := f.content
	for i < len(s) {
		if bytes.HasPrefix(s[i:], b) {
			return i
		}
		if i+2 <= len(s) && s[i] == '/' && s[i+1] == '/' {
			for i < len(s) && s[i] != '\n' {
				i++
			}
			continue
		}
		if i+2 <= len(s) && s[i] == '/' && s[i+1] == '*' {
			for i += 2; ; i++ {
				if i+2 > len(s) {
					return 0
				}
				if s[i] == '*' && s[i+1] == '/' {
					i += 2
					break
				}
			}
			continue
		}
		i++
	}
	return -1
}

func (f *FuncBranchVisitor) offset(pos token.Pos) int {
	return f.fset.Position(pos).Offset
}
