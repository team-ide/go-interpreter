package node

// Struct 结构体
type Struct interface {
	Node
	Statement
	Expression
	IsStructNode()
}

// StructElement 结构元素
type StructElement interface {
	Node
	IsStructElementNode()
}
