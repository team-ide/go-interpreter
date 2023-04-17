package node

// Interface 接口
type Interface interface {
	Node
	Statement
	Expression
	IsInterfaceNode()
}

// InterfaceElement 接口 元素
type InterfaceElement interface {
	Node
	IsInterfaceElementNode()
}
