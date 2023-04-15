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
	Namespace *node.ChainNameStatement
}

func (*NamespaceStatement) IsStatementNode() {}
func (this_ *NamespaceStatement) Start() int { return this_.From }
func (this_ *NamespaceStatement) End() int   { return this_.To }

// ExceptionStatement thrift 异常
type ExceptionStatement struct {
	From    int
	To      int
	Name    string
	Extends *node.ChainNameStatement
	Fields  []*FieldDefinition
}

func (*ExceptionStatement) IsStatementNode() {}
func (this_ *ExceptionStatement) Start() int { return this_.From }
func (this_ *ExceptionStatement) End() int   { return this_.To }

// StructStatement thrift 结构体
type StructStatement struct {
	From    int
	To      int
	Name    string
	Extends *node.ChainNameStatement
	Fields  []*FieldDefinition
}

func (*StructStatement) IsStatementNode() {}
func (this_ *StructStatement) Start() int { return this_.From }
func (this_ *StructStatement) End() int   { return this_.To }

// ServiceStatement 导入
type ServiceStatement struct {
	From    int
	To      int
	Name    string
	Extends *node.ChainNameStatement
	Methods []*IFaceMethodDefinition
}

func (*ServiceStatement) IsStatementNode() {}
func (this_ *ServiceStatement) Start() int { return this_.From }
func (this_ *ServiceStatement) End() int   { return this_.To }

// EnumStatement 导入
type EnumStatement struct {
	From    int
	To      int
	Name    string
	Extends *node.ChainNameStatement
	Fields  []*FieldDefinition
}

func (*EnumStatement) IsStatementNode() {}
func (this_ *EnumStatement) Start() int { return this_.From }
func (this_ *EnumStatement) End() int   { return this_.To }

// FieldType 字段类型
type FieldType interface {
	node.Node
	IsFieldTypeNode()
}

type FieldTypeName struct {
	From         int
	To           int
	Name         string
	GenericTypes []FieldType
}

func (*FieldTypeName) IsFieldTypeNode() {}
func (this_ *FieldTypeName) Start() int { return this_.From }
func (this_ *FieldTypeName) End() int   { return this_.To }

type FieldTypeDot struct {
	From         int
	To           int
	Names        []string
	GenericTypes []FieldType
}

func (*FieldTypeDot) IsFieldTypeNode() {}
func (this_ *FieldTypeDot) Start() int { return this_.From }
func (this_ *FieldTypeDot) End() int   { return this_.To }

// IFaceMethodDefinition 字段定义
type IFaceMethodDefinition struct {
	From   int
	To     int
	Return FieldType
	Name   string
	Params []*FieldDefinition
}

func (*IFaceMethodDefinition) IsExpressionNode() {}
func (*IFaceMethodDefinition) IsDefinitionNode() {}
func (this_ *IFaceMethodDefinition) Start() int  { return this_.From }
func (this_ *IFaceMethodDefinition) End() int    { return this_.To }

// FieldDefinition 字段定义
type FieldDefinition struct {
	From  int
	To    int
	Num   int
	Type  FieldType
	Name  string
	Value string
}

func (*FieldDefinition) IsClassElementNode() {}
func (this_ *FieldDefinition) Start() int    { return this_.From }
func (this_ *FieldDefinition) End() int      { return this_.To }
