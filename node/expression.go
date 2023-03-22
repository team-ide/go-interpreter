package node

import (
	"github.com/dop251/goja/unistring"
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
	Yield    *Position
	Argument Expression
	Delegate bool
}

// AwaitExpression 等待表达式
type AwaitExpression struct {
	Await    *Position
	Argument Expression
}

// ArrayLiteral 数组
type ArrayLiteral struct {
	LeftBracket  *Position
	RightBracket *Position
	Value        []Expression
}

// ArrayPattern 阵列模式
type ArrayPattern struct {
	LeftBracket  *Position
	RightBracket *Position
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
	From *Position
	To   *Position
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
	Idx     *Position
	Literal string
	Value   bool
}

// BracketExpression 括号表达式
type BracketExpression struct {
	Left         Expression
	Member       Expression
	LeftBracket  *Position
	RightBracket *Position
}

// CallExpression 调用表达式
type CallExpression struct {
	Callee           Expression
	LeftParenthesis  *Position
	ArgumentList     []Expression
	RightParenthesis *Position
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
	Function      *Position
	Name          *Identifier
	ParameterList *ParameterList
	Body          *BlockStatement
	Source        string

	DeclarationList []*VariableDeclaration

	Async, Generator bool
}

// ClassLiteral 类
type ClassLiteral struct {
	Class      *Position
	RightBrace *Position
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
	Start_          *Position
	ParameterList   *ParameterList
	Body            ConciseBody
	Source          string
	DeclarationList []*VariableDeclaration
	Async           bool
}

// Identifier 标识符
type Identifier struct {
	Name unistring.String
	Idx  *Position
}

// PrivateIdentifier 私有标识符
type PrivateIdentifier struct {
	Identifier
}

// NewExpression new标识符
type NewExpression struct {
	New              *Position
	Callee           Expression
	LeftParenthesis  *Position
	ArgumentList     []Expression
	RightParenthesis *Position
}

// NullLiteral null
type NullLiteral struct {
	Idx     *Position
	Literal string
}

// NumberLiteral 数字
type NumberLiteral struct {
	Idx     *Position
	Literal string
	Value   interface{}
}

// ObjectLiteral 对象
type ObjectLiteral struct {
	LeftBrace  *Position
	RightBrace *Position
	Value      []Property
}

// ObjectPattern 对象模式
type ObjectPattern struct {
	LeftBrace  *Position
	RightBrace *Position
	Properties []Property
	Rest       Expression
}

// ParameterList 参数列表
type ParameterList struct {
	Opening *Position
	List    []*Binding
	Rest    Expression
	Closing *Position
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
	Idx     *Position
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
	Idx     *Position
	Literal string
	Value   unistring.String
}

// TemplateElement 模板元素
type TemplateElement struct {
	Idx     *Position
	Literal string
	Parsed  unistring.String
	Valid   bool
}

// TemplateLiteral 模板
type TemplateLiteral struct {
	OpenQuote   *Position
	CloseQuote  *Position
	Tag         Expression
	Elements    []*TemplateElement
	Expressions []Expression
}

// ThisExpression this
type ThisExpression struct {
	Idx *Position
}

// SuperExpression super
type SuperExpression struct {
	Idx *Position
}

// UnaryExpression 一元表达式
type UnaryExpression struct {
	Operator token.Token
	Idx      *Position // If a prefix operation
	Operand  Expression
	Postfix  bool
}

// MetaProperty 元属性
type MetaProperty struct {
	Meta, Property *Identifier
	Idx            *Position
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

func (this_ *ArrayLiteral) Start() *Position          { return this_.LeftBracket }
func (this_ *ArrayPattern) Start() *Position          { return this_.LeftBracket }
func (this_ *YieldExpression) Start() *Position       { return this_.Yield }
func (this_ *AwaitExpression) Start() *Position       { return this_.Await }
func (this_ *ObjectPattern) Start() *Position         { return this_.LeftBrace }
func (this_ *ParameterList) Start() *Position         { return this_.Opening }
func (this_ *AssignExpression) Start() *Position      { return this_.Left.Start() }
func (this_ *BadExpression) Start() *Position         { return this_.From }
func (this_ *BinaryExpression) Start() *Position      { return this_.Left.Start() }
func (this_ *BooleanLiteral) Start() *Position        { return this_.Idx }
func (this_ *BracketExpression) Start() *Position     { return this_.Left.Start() }
func (this_ *CallExpression) Start() *Position        { return this_.Callee.Start() }
func (this_ *ConditionalExpression) Start() *Position { return this_.Test.Start() }
func (this_ *DotExpression) Start() *Position         { return this_.Left.Start() }
func (this_ *PrivateDotExpression) Start() *Position  { return this_.Left.Start() }
func (this_ *FunctionLiteral) Start() *Position       { return this_.Function }
func (this_ *ClassLiteral) Start() *Position          { return this_.Class }
func (this_ *ArrowFunctionLiteral) Start() *Position  { return this_.Start_ }
func (this_ *Identifier) Start() *Position            { return this_.Idx }
func (this_ *NewExpression) Start() *Position         { return this_.New }
func (this_ *NullLiteral) Start() *Position           { return this_.Idx }
func (this_ *NumberLiteral) Start() *Position         { return this_.Idx }
func (this_ *ObjectLiteral) Start() *Position         { return this_.LeftBrace }
func (this_ *RegExpLiteral) Start() *Position         { return this_.Idx }
func (this_ *SequenceExpression) Start() *Position    { return this_.Sequence[0].Start() }
func (this_ *StringLiteral) Start() *Position         { return this_.Idx }
func (this_ *TemplateElement) Start() *Position       { return this_.Idx }
func (this_ *TemplateLiteral) Start() *Position       { return this_.OpenQuote }
func (this_ *ThisExpression) Start() *Position        { return this_.Idx }
func (this_ *SuperExpression) Start() *Position       { return this_.Idx }
func (this_ *UnaryExpression) Start() *Position       { return this_.Idx }
func (this_ *MetaProperty) Start() *Position          { return this_.Idx }

/* 实现 Node End 接口 */

func (this_ *ArrayLiteral) End() *Position          { return this_.RightBracket.NewByColumnOffset(+1) }
func (this_ *ArrayPattern) End() *Position          { return this_.RightBracket.NewByColumnOffset(+1) }
func (this_ *AssignExpression) End() *Position      { return this_.Right.End() }
func (this_ *AwaitExpression) End() *Position       { return this_.Argument.End() }
func (this_ *BadExpression) End() *Position         { return this_.To }
func (this_ *BinaryExpression) End() *Position      { return this_.Right.End() }
func (this_ *BooleanLiteral) End() *Position        { return this_.Idx.NewByColumnOffset(len(this_.Literal)) }
func (this_ *BracketExpression) End() *Position     { return this_.RightBracket.NewByColumnOffset(+1) }
func (this_ *CallExpression) End() *Position        { return this_.RightParenthesis.NewByColumnOffset(+1) }
func (this_ *ConditionalExpression) End() *Position { return this_.Test.End() }
func (this_ *DotExpression) End() *Position         { return this_.Identifier.End() }
func (this_ *PrivateDotExpression) End() *Position  { return this_.Identifier.End() }
func (this_ *FunctionLiteral) End() *Position       { return this_.Body.End() }
func (this_ *ClassLiteral) End() *Position          { return this_.RightBrace.NewByColumnOffset(+1) }
func (this_ *ArrowFunctionLiteral) End() *Position  { return this_.Body.End() }
func (this_ *Identifier) End() *Position            { return this_.Idx.NewByColumnOffset(len(this_.Name)) }
func (this_ *NewExpression) End() *Position {
	if this_.ArgumentList != nil {
		return this_.RightParenthesis.NewByColumnOffset(+1)
	} else {
		return this_.Callee.End()
	}
}
func (this_ *NullLiteral) End() *Position        { return this_.Idx.NewByColumnOffset(+4) } // "null"
func (this_ *NumberLiteral) End() *Position      { return this_.Idx.NewByColumnOffset(len(this_.Literal)) }
func (this_ *ObjectLiteral) End() *Position      { return this_.RightBrace.NewByColumnOffset(+1) }
func (this_ *ObjectPattern) End() *Position      { return this_.RightBrace.NewByColumnOffset(+1) }
func (this_ *ParameterList) End() *Position      { return this_.Closing.NewByColumnOffset(+1) }
func (this_ *RegExpLiteral) End() *Position      { return this_.Idx.NewByColumnOffset(len(this_.Literal)) }
func (this_ *SequenceExpression) End() *Position { return this_.Sequence[len(this_.Sequence)-1].End() }
func (this_ *StringLiteral) End() *Position      { return this_.Idx.NewByColumnOffset(len(this_.Literal)) }
func (this_ *TemplateElement) End() *Position    { return this_.Idx.NewByColumnOffset(len(this_.Literal)) }
func (this_ *TemplateLiteral) End() *Position    { return this_.CloseQuote.NewByColumnOffset(+1) }
func (this_ *ThisExpression) End() *Position     { return this_.Idx.NewByColumnOffset(+4) }
func (this_ *SuperExpression) End() *Position    { return this_.Idx.NewByColumnOffset(+5) }
func (this_ *UnaryExpression) End() *Position {
	if this_.Postfix {
		return this_.Operand.End().NewByColumnOffset(+2) // ++ --
	}
	return this_.Operand.End()
}
func (this_ *MetaProperty) End() *Position {
	return this_.Property.End()
}
