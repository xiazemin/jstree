package call

type BranchNode struct {
	Comment  string
	Code     string
	Children []*BranchNode
	Parent   *BranchNode
	Return   bool
	Else     bool
}

func (n *BranchNode) AddChild(child *BranchNode) {
	n.Children = append(n.Children, child)
	// child.Parent = n
}
