package thrift

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/team-ide/go-interpreter/node"
)

// IncludeStatement thrift 导入
type IncludeStatement struct {
	From    int    `json:"from"`
	To      int    `json:"to"`
	Include string `json:"include"`
}

func (*IncludeStatement) IsStatementNode() {}
func (this_ *IncludeStatement) Start() int { return this_.From }
func (this_ *IncludeStatement) End() int   { return this_.To }

// NamespaceStatement thrift 生成语言命名空间
type NamespaceStatement struct {
	From      int                      `json:"from"`
	To        int                      `json:"to"`
	Language  string                   `json:"language"`
	Namespace *node.ChainNameStatement `json:"namespace"`
}

func (*NamespaceStatement) IsStatementNode() {}
func (this_ *NamespaceStatement) Start() int { return this_.From }
func (this_ *NamespaceStatement) End() int   { return this_.To }

// ExceptionStatement thrift 异常
type ExceptionStatement struct {
	*StructStatement
}

func (*ExceptionStatement) IsStatementNode() {}
func (this_ *ExceptionStatement) Start() int { return this_.From }
func (this_ *ExceptionStatement) End() int   { return this_.To }

// StructStatement thrift 结构体
type StructStatement struct {
	From           int          `json:"from"`
	To             int          `json:"to"`
	Name           string       `json:"name"`
	ExtendsInclude string       `json:"extendsInclude"`
	ExtendsName    string       `json:"extendsName"`
	Fields         []*FieldNode `json:"fields"`
}

func (*StructStatement) IsStatementNode() {}
func (this_ *StructStatement) Start() int { return this_.From }
func (this_ *StructStatement) End() int   { return this_.To }

// ServiceStatement 导入
type ServiceStatement struct {
	From           int                  `json:"from"`
	To             int                  `json:"to"`
	Name           string               `json:"name"`
	ExtendsInclude string               `json:"extendsInclude"`
	ExtendsName    string               `json:"extendsName"`
	Methods        []*ServiceMethodNode `json:"methods"`
}

func (*ServiceStatement) IsStatementNode() {}
func (this_ *ServiceStatement) Start() int { return this_.From }
func (this_ *ServiceStatement) End() int   { return this_.To }

// EnumStatement 导入
type EnumStatement struct {
	From           int    `json:"from"`
	To             int    `json:"to"`
	Name           string `json:"name"`
	ExtendsInclude string `json:"extendsInclude"`
	ExtendsName    string `json:"extendsName"`
	Fields         []*FieldNode
}

func (*EnumStatement) IsStatementNode() {}
func (this_ *EnumStatement) Start() int { return this_.From }
func (this_ *EnumStatement) End() int   { return this_.To }

// FieldType 字段类型
type FieldType struct {
	From          int          `json:"from"`
	To            int          `json:"to"`
	TypeId        thrift.TType `json:"typeId"`
	TypeName      string       `json:"typeName"`
	StructInclude string       `json:"structInclude"`
	StructName    string       `json:"structName"`
	ListType      *FieldType   `json:"listType"`
	SetType       *FieldType   `json:"setType"`
	MapKeyType    *FieldType   `json:"mapKeyType"`
	MapValueType  *FieldType   `json:"mapValueType"`
}

func (this_ *FieldType) Start() int { return this_.From }
func (this_ *FieldType) End() int   { return this_.To }

// ServiceMethodNode 服务接口方法
type ServiceMethodNode struct {
	From       int          `json:"from"`
	To         int          `json:"to"`
	Return     *FieldType   `json:"return"`
	Name       string       `json:"name"`
	Exceptions []*FieldNode `json:"exceptions"`
	Params     []*FieldNode `json:"params"`
}

func (this_ *ServiceMethodNode) Start() int { return this_.From }
func (this_ *ServiceMethodNode) End() int   { return this_.To }

// FieldNode 字段定义
type FieldNode struct {
	From     int        `json:"from"`
	To       int        `json:"to"`
	Num      int16      `json:"num"`
	Optional bool       `json:"optional"`
	Type     *FieldType `json:"type"`
	Name     string     `json:"name"`
	Value    string     `json:"value"`
}

func (this_ *FieldNode) Start() int { return this_.From }
func (this_ *FieldNode) End() int   { return this_.To }
