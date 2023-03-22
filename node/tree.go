package node

// Tree 属性结构
type Tree struct {
	Children        []Statement
	DeclarationList []*VariableDeclaration
}

func (this_ *Tree) Start() *Position {
	if len(this_.Children) == 0 {
		return nil
	}
	return this_.Children[0].Start()
}

func (this_ *Tree) End() *Position {
	if len(this_.Children) == 0 {
		return nil
	}
	return this_.Children[len(this_.Children)-1].End()
}
