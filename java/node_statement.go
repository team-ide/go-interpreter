package java

// ImportStatement 导入
type ImportStatement struct {
	From   int
	To     int
	Import string
}

func (*ImportStatement) IsStatementNode() {}
func (this_ *ImportStatement) Start() int { return this_.From }
func (this_ *ImportStatement) End() int   { return this_.To }
