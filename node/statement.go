package node

/** 语句 **/

// Statement 所有语句节点都实现了 Statement 接口
type Statement interface {
	Node
	IsStatementNode()
}

// BadStatement 错误的语句
type BadStatement struct {
	From int
	To   int
}

func (*BadStatement) IsStatementNode() {}
func (this_ *BadStatement) Start() int { return this_.From }
func (this_ *BadStatement) End() int   { return this_.To }

// ChainNameStatement 链式名称 如 x、x.xx、xx.xx.xxx
type ChainNameStatement struct {
	From  int
	To    int
	Names []string
}

func (*ChainNameStatement) IsStatementNode() {}
func (this_ *ChainNameStatement) Start() int { return this_.From }
func (this_ *ChainNameStatement) End() int   { return this_.To }
