package node

import (
	"github.com/team-ide/go-interpreter/token"
)

/** 表达式 **/

// Expression 所有表达式节点都实现 Expression 接口
type Expression interface {
	Node
	isExpression()
}

// BindingTarget 绑定目标
type BindingTarget interface {
	Expression
	isBindingTarget()
}

// Binding 绑定
type Binding struct {
	Target      BindingTarget
	Initializer Expression
}

// Pattern 模式
type Pattern interface {
	BindingTarget
	isPattern()
}

// YieldExpression 收益率表达式
type YieldExpression struct {
	Yield    int
	Argument Expression
	Delegate bool
}

// AwaitExpression 等待表达式
type AwaitExpression struct {
	Await    int
	Argument Expression
}

// ArrayLiteral 数组
type ArrayLiteral struct {
	LeftBracket  int
	RightBracket int
	Value        []Expression
}

// ArrayPattern 阵列模式
type ArrayPattern struct {
	LeftBracket  int
	RightBracket int
	Elements     []Expression
	Rest         Expression
}

// AssignExpression 指定表达式
type AssignExpression struct {
	Operator token.Token
	Left     Expression
	Right    Expression
}

// BadExpression 错误的表达式
type BadExpression struct {
	From int
	To   int
}

// BinaryExpression 二进制表达式
type BinaryExpression struct {
	Operator   token.Token
	Left       Expression
	Right      Expression
	Comparison bool
}

// BooleanLiteral 布尔
type BooleanLiteral struct {
	Idx     int
	Literal string
	Value   bool
}

// BracketExpression 括号表达式
type BracketExpression struct {
	Left         Expression
	Member       Expression
	LeftBracket  int
	RightBracket int
}

// CallExpression 调用表达式
type CallExpression struct {
	Callee           Expression
	LeftParenthesis  int
	ArgumentList     []Expression
	RightParenthesis int
}

// ConditionalExpression 条件表达式
type ConditionalExpression struct {
	Test       Expression
	Consequent Expression
	Alternate  Expression
}

// DotExpression 点表达式
type DotExpression struct {
	Left       Expression
	Identifier Identifier
}

// PrivateDotExpression 私有点表达式
type PrivateDotExpression struct {
	Left       Expression
	Identifier PrivateIdentifier
}

// OptionalChain 可选链条
type OptionalChain struct {
	Expression
}

// Optional 可选
type Optional struct {
	Expression
}

// FunctionLiteral 函数
type FunctionLiteral struct {
	Function      int
	Name          *Identifier
	ParameterList *ParameterList
	Body          *BlockStatement
	Source        string

	DeclarationList []*VariableDeclaration

	Async, Generator bool
}

// ClassLiteral 类
type ClassLiteral struct {
	Class      int
	RightBrace int
	Name       *Identifier
	SuperClass Expression
	Body       []ClassElement
	Source     string
}

// ConciseBody 简明正文
type ConciseBody interface {
	Node
	isConciseBody()
}

// ExpressionBody 表达式正文
type ExpressionBody struct {
	Expression Expression
}

// ArrowFunctionLiteral 箭头函数
type ArrowFunctionLiteral struct {
	Idx             int
	ParameterList   *ParameterList
	Body            ConciseBody
	Source          string
	DeclarationList []*VariableDeclaration
	Async           bool
}

// Identifier 标识符
type Identifier struct {
	Name String
	Idx  int
}

// PrivateIdentifier 私有标识符
type PrivateIdentifier struct {
	Identifier
}

// NewExpression new标识符
type NewExpression struct {
	New              int
	Callee           Expression
	LeftParenthesis  int
	ArgumentList     []Expression
	RightParenthesis int
}

// NullLiteral null
type NullLiteral struct {
	Idx     int
	Literal string
}

// NumberLiteral 数字
type NumberLiteral struct {
	Idx     int
	Literal string
	Value   interface{}
}

// ObjectLiteral 对象
type ObjectLiteral struct {
	LeftBrace  int
	RightBrace int
	Value      []Property
}

// ObjectPattern 对象模式
type ObjectPattern struct {
	LeftBrace  int
	RightBrace int
	Properties []Property
	Rest       Expression
}

// ParameterList 参数列表
type ParameterList struct {
	Opening int
	List    []*Binding
	Rest    Expression
	Closing int
}

// Property 属性
type Property interface {
	Expression
	isProperty()
}

// PropertyShort 属性缩写
type PropertyShort struct {
	Name        Identifier
	Initializer Expression
}

// PropertyKeyed 属性映射
type PropertyKeyed struct {
	Key      Expression
	Kind     PropertyKind
	Value    Expression
	Computed bool
}

// SpreadElement 排列元素
type SpreadElement struct {
	Expression
}

// RegExpLiteral REG分解
type RegExpLiteral struct {
	Idx     int
	Literal string
	Pattern string
	Flags   string
}

// SequenceExpression 序列表达式
type SequenceExpression struct {
	Sequence []Expression
}

// StringLiteral 字符串
type StringLiteral struct {
	Idx     int
	Literal string
	Value   String
}

// TemplateElement 模板元素
type TemplateElement struct {
	Idx     int
	Literal string
	Parsed  String
	Valid   bool
}

// TemplateLiteral 模板
type TemplateLiteral struct {
	OpenQuote   int
	CloseQuote  int
	Tag         Expression
	Elements    []*TemplateElement
	Expressions []Expression
}

// ThisExpression this
type ThisExpression struct {
	Idx int
}

// SuperExpression super
type SuperExpression struct {
	Idx int
}

// UnaryExpression 一元表达式
type UnaryExpression struct {
	Operator token.Token
	Idx      int // If a prefix operation
	Operand  Expression
	Postfix  bool
}

// MetaProperty 元属性
type MetaProperty struct {
	Meta, Property *Identifier
	Idx            int
}

/* 实现 Expression 接口 */
func (*ArrayLiteral) isExpression()          {}
func (*AssignExpression) isExpression()      {}
func (*YieldExpression) isExpression()       {}
func (*AwaitExpression) isExpression()       {}
func (*BadExpression) isExpression()         {}
func (*BinaryExpression) isExpression()      {}
func (*BooleanLiteral) isExpression()        {}
func (*BracketExpression) isExpression()     {}
func (*CallExpression) isExpression()        {}
func (*ConditionalExpression) isExpression() {}
func (*DotExpression) isExpression()         {}
func (*PrivateDotExpression) isExpression()  {}
func (*FunctionLiteral) isExpression()       {}
func (*ClassLiteral) isExpression()          {}
func (*ArrowFunctionLiteral) isExpression()  {}
func (*Identifier) isExpression()            {}
func (*NewExpression) isExpression()         {}
func (*NullLiteral) isExpression()           {}
func (*NumberLiteral) isExpression()         {}
func (*ObjectLiteral) isExpression()         {}
func (*RegExpLiteral) isExpression()         {}
func (*SequenceExpression) isExpression()    {}
func (*StringLiteral) isExpression()         {}
func (*TemplateLiteral) isExpression()       {}
func (*ThisExpression) isExpression()        {}
func (*SuperExpression) isExpression()       {}
func (*UnaryExpression) isExpression()       {}
func (*MetaProperty) isExpression()          {}
func (*ObjectPattern) isExpression()         {}
func (*ArrayPattern) isExpression()          {}
func (*Binding) isExpression()               {}

func (*PropertyShort) isExpression() {}
func (*PropertyKeyed) isExpression() {}

/* 实现 Node Start 接口 */

func (this_ *ArrayLiteral) Start() int          { return this_.LeftBracket }
func (this_ *ArrayPattern) Start() int          { return this_.LeftBracket }
func (this_ *YieldExpression) Start() int       { return this_.Yield }
func (this_ *AwaitExpression) Start() int       { return this_.Await }
func (this_ *ObjectPattern) Start() int         { return this_.LeftBrace }
func (this_ *ParameterList) Start() int         { return this_.Opening }
func (this_ *AssignExpression) Start() int      { return this_.Left.Start() }
func (this_ *BadExpression) Start() int         { return this_.From }
func (this_ *BinaryExpression) Start() int      { return this_.Left.Start() }
func (this_ *BooleanLiteral) Start() int        { return this_.Idx }
func (this_ *BracketExpression) Start() int     { return this_.Left.Start() }
func (this_ *CallExpression) Start() int        { return this_.Callee.Start() }
func (this_ *ConditionalExpression) Start() int { return this_.Test.Start() }
func (this_ *DotExpression) Start() int         { return this_.Left.Start() }
func (this_ *PrivateDotExpression) Start() int  { return this_.Left.Start() }
func (this_ *FunctionLiteral) Start() int       { return this_.Function }
func (this_ *ClassLiteral) Start() int          { return this_.Class }
func (this_ *ArrowFunctionLiteral) Start() int  { return this_.Idx }
func (this_ *Identifier) Start() int            { return this_.Idx }
func (this_ *NewExpression) Start() int         { return this_.New }
func (this_ *NullLiteral) Start() int           { return this_.Idx }
func (this_ *NumberLiteral) Start() int         { return this_.Idx }
func (this_ *ObjectLiteral) Start() int         { return this_.LeftBrace }
func (this_ *RegExpLiteral) Start() int         { return this_.Idx }
func (this_ *SequenceExpression) Start() int    { return this_.Sequence[0].Start() }
func (this_ *StringLiteral) Start() int         { return this_.Idx }
func (this_ *TemplateElement) Start() int       { return this_.Idx }
func (this_ *TemplateLiteral) Start() int       { return this_.OpenQuote }
func (this_ *ThisExpression) Start() int        { return this_.Idx }
func (this_ *SuperExpression) Start() int       { return this_.Idx }
func (this_ *UnaryExpression) Start() int       { return this_.Idx }
func (this_ *MetaProperty) Start() int          { return this_.Idx }

/* 实现 Node End 接口 */

func (this_ *ArrayLiteral) End() int          { return this_.RightBracket + (+1) }
func (this_ *ArrayPattern) End() int          { return this_.RightBracket + (+1) }
func (this_ *AssignExpression) End() int      { return this_.Right.End() }
func (this_ *AwaitExpression) End() int       { return this_.Argument.End() }
func (this_ *BadExpression) End() int         { return this_.To }
func (this_ *BinaryExpression) End() int      { return this_.Right.End() }
func (this_ *BooleanLiteral) End() int        { return this_.Idx + (len(this_.Literal)) }
func (this_ *BracketExpression) End() int     { return this_.RightBracket + (+1) }
func (this_ *CallExpression) End() int        { return this_.RightParenthesis + (+1) }
func (this_ *ConditionalExpression) End() int { return this_.Test.End() }
func (this_ *DotExpression) End() int         { return this_.Identifier.End() }
func (this_ *PrivateDotExpression) End() int  { return this_.Identifier.End() }
func (this_ *FunctionLiteral) End() int       { return this_.Body.End() }
func (this_ *ClassLiteral) End() int          { return this_.RightBrace + (+1) }
func (this_ *ArrowFunctionLiteral) End() int  { return this_.Body.End() }
func (this_ *Identifier) End() int            { return this_.Idx + (len(this_.Name)) }
func (this_ *NewExpression) End() int {
	if this_.ArgumentList != nil {
		return this_.RightParenthesis + (+1)
	} else {
		return this_.Callee.End()
	}
}
func (this_ *NullLiteral) End() int        { return this_.Idx + (+4) } // "null"
func (this_ *NumberLiteral) End() int      { return this_.Idx + (len(this_.Literal)) }
func (this_ *ObjectLiteral) End() int      { return this_.RightBrace + (+1) }
func (this_ *ObjectPattern) End() int      { return this_.RightBrace + (+1) }
func (this_ *ParameterList) End() int      { return this_.Closing + (+1) }
func (this_ *RegExpLiteral) End() int      { return this_.Idx + (len(this_.Literal)) }
func (this_ *SequenceExpression) End() int { return this_.Sequence[len(this_.Sequence)-1].End() }
func (this_ *StringLiteral) End() int      { return this_.Idx + (len(this_.Literal)) }
func (this_ *TemplateElement) End() int    { return this_.Idx + (len(this_.Literal)) }
func (this_ *TemplateLiteral) End() int    { return this_.CloseQuote + (+1) }
func (this_ *ThisExpression) End() int     { return this_.Idx + (+4) }
func (this_ *SuperExpression) End() int    { return this_.Idx + (+5) }
func (this_ *UnaryExpression) End() int {
	if this_.Postfix {
		return this_.Operand.End() + (+2) // ++ --
	}
	return this_.Operand.End()
}
func (this_ *MetaProperty) End() int {
	return this_.Property.End()
}
