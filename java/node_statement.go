package java

import "github.com/team-ide/go-interpreter/node"

// ImportStatement 导入
type ImportStatement struct {
	From   int
	To     int
	Import *node.ChainNameStatement
}

func (*ImportStatement) IsStatementNode() {}
func (this_ *ImportStatement) Start() int { return this_.From }
func (this_ *ImportStatement) End() int   { return this_.To }
