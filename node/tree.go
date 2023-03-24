package node

// Tree 属性结构
type Tree struct {
	Children        []Statement
	DeclarationList []*VariableDeclaration
}

func (this_ *Tree) Start() int {
	if len(this_.Children) == 0 {
		return 0
	}
	return this_.Children[0].Start()
}

func (this_ *Tree) End() int {
	if len(this_.Children) == 0 {
		return 0
	}
	return this_.Children[len(this_.Children)-1].End()
}
