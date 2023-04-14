package thrift

import "github.com/team-ide/go-interpreter/node"

// IncludeStatement thrift 导入
type IncludeStatement struct {
	From    int
	To      int
	Include string
}

func (*IncludeStatement) IsStatementNode() {}
func (this_ *IncludeStatement) Start() int { return this_.From }
func (this_ *IncludeStatement) End() int   { return this_.To }

// NamespaceStatement thrift 生成语言命名空间
type NamespaceStatement struct {
	From      int
	To        int
	Language  string
	Namespace string
}

func (*NamespaceStatement) IsStatementNode() {}
func (this_ *NamespaceStatement) Start() int { return this_.From }
func (this_ *NamespaceStatement) End() int   { return this_.To }

// ExceptionStatement thrift 异常
type ExceptionStatement struct {
	From   int
	To     int
	Name   string
	Fields []*FieldDefinition
}

func (*ExceptionStatement) IsStatementNode() {}
func (this_ *ExceptionStatement) Start() int { return this_.From }
func (this_ *ExceptionStatement) End() int   { return this_.To }

// StructStatement thrift 结构体
type StructStatement struct {
	From   int
	To     int
	Name   string
	Fields []*FieldDefinition
}

func (*StructStatement) IsStatementNode() {}
func (this_ *StructStatement) Start() int { return this_.From }
func (this_ *StructStatement) End() int   { return this_.To }

// ServiceStatement 导入
type ServiceStatement struct {
	From    int
	To      int
	Name    string
	Methods []*IFaceMethodDefinition
}

func (*ServiceStatement) IsStatementNode() {}
func (this_ *ServiceStatement) Start() int { return this_.From }
func (this_ *ServiceStatement) End() int   { return this_.To }

// EnumStatement 导入
type EnumStatement struct {
	From   int
	To     int
	Name   string
	Fields []*FieldDefinition
}

func (*EnumStatement) IsStatementNode() {}
func (this_ *EnumStatement) Start() int { return this_.From }
func (this_ *EnumStatement) End() int   { return this_.To }

// IFaceMethodDefinition 字段定义
type IFaceMethodDefinition struct {
	From   int
	To     int
	Return node.Expression
	Name   node.Expression
	Params []*FieldDefinition
}

func (*IFaceMethodDefinition) IsExpressionNode() {}
func (*IFaceMethodDefinition) IsDefinitionNode() {}
func (this_ *IFaceMethodDefinition) Start() int  { return this_.From }
func (this_ *IFaceMethodDefinition) End() int    { return this_.To }

// FieldDefinition 字段定义
type FieldDefinition struct {
	Idx         int
	FieldNum    int
	Type        node.Expression
	Key         node.Expression
	Initializer node.Expression
	Computed    bool
	Static      bool
}

func (*FieldDefinition) IsClassElementNode() {}
func (this_ *FieldDefinition) Start() int    { return this_.Idx }
func (this_ *FieldDefinition) End() int {
	if this_.Initializer != nil {
		return this_.Initializer.End()
	}
	if this_.Key != nil {
		return this_.Key.End()
	}
	return this_.Idx
}
