package node

/** 表达式 **/

// Expression 所有表达式节点都实现 Expression 接口
type Expression interface {
	Node
	IsExpressionNode()
}
