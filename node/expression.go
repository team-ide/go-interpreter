package node

import (
	"github.com/team-ide/go-interpreter/token"
)

/** 表达式 **/

// Expression 所有表达式节点都实现 Expression 接口
type Expression interface {
	Node
	IsExpressionNode()
}

// BadExpression 错误的表达式
type BadExpression struct {
	From int
	To   int
}

func (*BadExpression) IsExpressionNode()    {}
func (*BadExpression) IsBindingTargetNode() {}
func (this_ *BadExpression) Start() int     { return this_.From }
func (this_ *BadExpression) End() int       { return this_.To }

// BindingTarget 绑定目标
type BindingTarget interface {
	Expression
	IsBindingTargetNode()
}

// Binding 绑定
type Binding struct {
	Target      BindingTarget
	Initializer Expression
}

func (*Binding) IsExpressionNode() {}
func (this_ *Binding) Start() int  { return this_.Target.Start() }
func (this_ *Binding) End() int {
	if this_.Initializer != nil {
		return this_.Initializer.End()
	}
	return this_.Target.End()
}

// Pattern 模式
type Pattern interface {
	BindingTarget
	IsPatternNode()
}

// YieldExpression 收益率表达式
type YieldExpression struct {
	Yield    int
	Argument Expression
	Delegate bool
}

func (*YieldExpression) IsExpressionNode() {}
func (this_ *YieldExpression) Start() int  { return this_.Yield }

func (this_ *YieldExpression) End() int {
	if this_.Argument != nil {
		return this_.Argument.End()
	}
	return this_.Yield + (+5)
}

// AwaitExpression 等待表达式
type AwaitExpression struct {
	Await    int
	Argument Expression
}

func (*AwaitExpression) IsExpressionNode() {}
func (this_ *AwaitExpression) Start() int  { return this_.Await }
func (this_ *AwaitExpression) End() int    { return this_.Argument.End() }

// ArrayLiteral 数组
type ArrayLiteral struct {
	LeftBracket  int
	RightBracket int
	Value        []Expression
}

func (*ArrayLiteral) IsExpressionNode() {}
func (this_ *ArrayLiteral) Start() int  { return this_.LeftBracket }
func (this_ *ArrayLiteral) End() int    { return this_.RightBracket + (+1) }

// ArrayPattern 阵列模式
type ArrayPattern struct {
	LeftBracket  int
	RightBracket int
	Elements     []Expression
	Rest         Expression
}

func (*ArrayPattern) IsExpressionNode()    {}
func (*ArrayPattern) IsPatternNode()       {}
func (*ArrayPattern) IsBindingTargetNode() {}
func (this_ *ArrayPattern) Start() int     { return this_.LeftBracket }
func (this_ *ArrayPattern) End() int       { return this_.RightBracket + (+1) }

// AssignExpression 指定表达式
type AssignExpression struct {
	Operator token.Token
	Left     Expression
	Right    Expression
}

func (*AssignExpression) IsExpressionNode() {}
func (this_ *AssignExpression) Start() int  { return this_.Left.Start() }
func (this_ *AssignExpression) End() int    { return this_.Right.End() }

// BinaryExpression 二进制表达式
type BinaryExpression struct {
	Operator   token.Token
	Left       Expression
	Right      Expression
	Comparison bool
}

func (*BinaryExpression) IsExpressionNode() {}
func (this_ *BinaryExpression) Start() int  { return this_.Left.Start() }
func (this_ *BinaryExpression) End() int    { return this_.Right.End() }

// BooleanLiteral 布尔
type BooleanLiteral struct {
	Idx     int
	Literal string
	Value   bool
}

func (*BooleanLiteral) IsExpressionNode() {}
func (this_ *BooleanLiteral) Start() int  { return this_.Idx }
func (this_ *BooleanLiteral) End() int    { return this_.Idx + (len(this_.Literal)) }

// BracketExpression 括号表达式
type BracketExpression struct {
	Left         Expression
	Member       Expression
	LeftBracket  int
	RightBracket int
}

func (*BracketExpression) IsExpressionNode() {}
func (this_ *BracketExpression) Start() int  { return this_.Left.Start() }
func (this_ *BracketExpression) End() int    { return this_.RightBracket + (+1) }

// CallExpression 调用表达式
type CallExpression struct {
	Callee           Expression
	LeftParenthesis  int
	ArgumentList     []Expression
	RightParenthesis int
}

func (*CallExpression) IsExpressionNode() {}
func (this_ *CallExpression) Start() int  { return this_.Callee.Start() }
func (this_ *CallExpression) End() int    { return this_.RightParenthesis + (+1) }

// ConditionalExpression 条件表达式
type ConditionalExpression struct {
	Test       Expression
	Consequent Expression
	Alternate  Expression
}

func (*ConditionalExpression) IsExpressionNode() {}
func (this_ *ConditionalExpression) Start() int  { return this_.Test.Start() }
func (this_ *ConditionalExpression) End() int    { return this_.Test.End() }

// DotExpression 点表达式
type DotExpression struct {
	Left       Expression
	Identifier Identifier
}

func (*DotExpression) IsExpressionNode() {}
func (this_ *DotExpression) Start() int  { return this_.Left.Start() }
func (this_ *DotExpression) End() int    { return this_.Identifier.End() }

// PrivateDotExpression 私有点表达式
type PrivateDotExpression struct {
	Left       Expression
	Identifier PrivateIdentifier
}

func (*PrivateDotExpression) IsExpressionNode() {}
func (this_ *PrivateDotExpression) Start() int  { return this_.Left.Start() }
func (this_ *PrivateDotExpression) End() int    { return this_.Identifier.End() }

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

func (*FunctionLiteral) IsExpressionNode() {}
func (this_ *FunctionLiteral) Start() int  { return this_.Function }
func (this_ *FunctionLiteral) End() int    { return this_.Body.End() }

// ClassLiteral 类
type ClassLiteral struct {
	Class      int
	RightBrace int
	Name       *Identifier
	Extend     Expression
	Body       []ClassElement
	Implements []Expression // 实现的接口
	Source     string
}

func (*ClassLiteral) IsExpressionNode() {}
func (this_ *ClassLiteral) Start() int  { return this_.Class }
func (this_ *ClassLiteral) End() int    { return this_.RightBrace + (+1) }

// ConciseBody 简明正文
type ConciseBody interface {
	Node
	IsConciseBodyNode()
}

// ExpressionBody 表达式正文
type ExpressionBody struct {
	Expression Expression
}

func (*ExpressionBody) IsConciseBodyNode() {}
func (this_ *ExpressionBody) Start() int   { return this_.Expression.Start() }
func (this_ *ExpressionBody) End() int     { return this_.Expression.End() }

// ArrowFunctionLiteral 箭头函数
type ArrowFunctionLiteral struct {
	Idx             int
	ParameterList   *ParameterList
	Body            ConciseBody
	Source          string
	DeclarationList []*VariableDeclaration
	Async           bool
}

func (*ArrowFunctionLiteral) IsExpressionNode() {}
func (this_ *ArrowFunctionLiteral) Start() int  { return this_.Idx }
func (this_ *ArrowFunctionLiteral) End() int    { return this_.Body.End() }

// Identifier 标识符
type Identifier struct {
	Name string
	Idx  int
}

func (*Identifier) IsExpressionNode()    {}
func (*Identifier) IsBindingTargetNode() {}
func (this_ *Identifier) Start() int     { return this_.Idx }
func (this_ *Identifier) End() int       { return this_.Idx + (len(this_.Name)) }

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

func (*NewExpression) IsExpressionNode() {}
func (this_ *NewExpression) Start() int  { return this_.New }
func (this_ *NewExpression) End() int {
	end := this_.RightParenthesis + (+1)
	if end < this_.Callee.End() {
		end = this_.Callee.End()
	}
	return end
}

// NullLiteral null
type NullLiteral struct {
	Idx     int
	Literal string
}

func (*NullLiteral) IsExpressionNode() {}
func (this_ *NullLiteral) Start() int  { return this_.Idx }
func (this_ *NullLiteral) End() int    { return this_.Idx + (+4) } // "null"

// NumberLiteral 数字
type NumberLiteral struct {
	Idx     int
	Literal string
	Value   interface{}
}

func (*NumberLiteral) IsExpressionNode() {}
func (this_ *NumberLiteral) Start() int  { return this_.Idx }
func (this_ *NumberLiteral) End() int    { return this_.Idx + (len(this_.Literal)) }

// ObjectLiteral 对象
type ObjectLiteral struct {
	LeftBrace  int
	RightBrace int
	Value      []Property
}

func (*ObjectLiteral) IsExpressionNode() {}
func (this_ *ObjectLiteral) Start() int  { return this_.LeftBrace }
func (this_ *ObjectLiteral) End() int    { return this_.RightBrace + (+1) }

// ObjectPattern 对象模式
type ObjectPattern struct {
	LeftBrace  int
	RightBrace int
	Properties []Property
	Rest       Expression
}

func (*ObjectPattern) IsExpressionNode()    {}
func (*ObjectPattern) IsPatternNode()       {}
func (*ObjectPattern) IsBindingTargetNode() {}
func (this_ *ObjectPattern) Start() int     { return this_.LeftBrace }
func (this_ *ObjectPattern) End() int       { return this_.RightBrace + (+1) }

// ParameterList 参数列表
type ParameterList struct {
	Opening int
	List    []*Binding
	Rest    Expression
	Closing int
}

func (this_ *ParameterList) Start() int { return this_.Opening }
func (this_ *ParameterList) End() int   { return this_.Closing + (+1) }

// Property 属性
type Property interface {
	Expression
	IsPropertyNode()
}

// PropertyShort 属性缩写
type PropertyShort struct {
	Name        Identifier
	Initializer Expression
}

func (*PropertyShort) IsExpressionNode() {}
func (*PropertyShort) IsPropertyNode()   {}
func (this_ *PropertyShort) Start() int  { return this_.Name.Idx }
func (this_ *PropertyShort) End() int {
	if this_.Initializer != nil {
		return this_.Initializer.End()
	}
	return this_.Name.End()
}

// PropertyKeyed 属性映射
type PropertyKeyed struct {
	Key      Expression
	Kind     PropertyKind
	Value    Expression
	Computed bool
}

func (*PropertyKeyed) IsExpressionNode() {}
func (*PropertyKeyed) IsPropertyNode()   {}
func (this_ *PropertyKeyed) Start() int  { return this_.Key.Start() }
func (this_ *PropertyKeyed) End() int    { return this_.Value.End() }

// SpreadElement 排列元素
type SpreadElement struct {
	Expression
}

func (*SpreadElement) IsPropertyNode() {}

// RegExpLiteral REG分解
type RegExpLiteral struct {
	Idx     int
	Literal string
	Pattern string
	Flags   string
}

func (*RegExpLiteral) IsExpressionNode() {}
func (this_ *RegExpLiteral) Start() int  { return this_.Idx }
func (this_ *RegExpLiteral) End() int    { return this_.Idx + (len(this_.Literal)) }

// SequenceExpression 序列表达式
type SequenceExpression struct {
	Sequence []Expression
}

func (*SequenceExpression) IsExpressionNode() {}
func (this_ *SequenceExpression) Start() int  { return this_.Sequence[0].Start() }
func (this_ *SequenceExpression) End() int    { return this_.Sequence[len(this_.Sequence)-1].End() }

// StringLiteral 字符串
type StringLiteral struct {
	Idx     int
	Literal string
	Value   string
}

func (*StringLiteral) IsExpressionNode() {}
func (this_ *StringLiteral) Start() int  { return this_.Idx }
func (this_ *StringLiteral) End() int    { return this_.Idx + (len(this_.Literal)) }

// TemplateElement 模板元素
type TemplateElement struct {
	Idx     int
	Literal string
	Parsed  string
	Valid   bool
}

func (this_ *TemplateElement) Start() int { return this_.Idx }
func (this_ *TemplateElement) End() int   { return this_.Idx + (len(this_.Literal)) }

// TemplateLiteral 模板
type TemplateLiteral struct {
	OpenQuote   int
	CloseQuote  int
	Tag         Expression
	Elements    []*TemplateElement
	Expressions []Expression
}

func (*TemplateLiteral) IsExpressionNode() {}
func (this_ *TemplateLiteral) Start() int  { return this_.OpenQuote }
func (this_ *TemplateLiteral) End() int    { return this_.CloseQuote + (+1) }

// ThIsExpressionNode this
type ThIsExpressionNode struct {
	Idx int
}

func (*ThIsExpressionNode) IsExpressionNode() {}
func (this_ *ThIsExpressionNode) Start() int  { return this_.Idx }
func (this_ *ThIsExpressionNode) End() int    { return this_.Idx + (+4) }

// SuperExpression super
type SuperExpression struct {
	Idx int
}

func (*SuperExpression) IsExpressionNode() {}
func (this_ *SuperExpression) Start() int  { return this_.Idx }
func (this_ *SuperExpression) End() int    { return this_.Idx + (+5) }

// UnaryExpression 一元表达式
type UnaryExpression struct {
	Operator token.Token
	Idx      int // If a prefix operation
	Operand  Expression
	Postfix  bool
}

func (*UnaryExpression) IsExpressionNode() {}
func (this_ *UnaryExpression) Start() int {
	if this_.Postfix {
		return this_.Operand.Start()
	}
	return this_.Idx
}
func (this_ *UnaryExpression) End() int {
	if this_.Postfix {
		return this_.Operand.End() + (+2) // ++ --
	}
	return this_.Operand.End()
}

// MetaProperty 元属性
type MetaProperty struct {
	Meta, Property *Identifier
	Idx            int
}

func (*MetaProperty) IsExpressionNode() {}
func (this_ *MetaProperty) Start() int  { return this_.Idx }
func (this_ *MetaProperty) End() int {
	return this_.Property.End()
}
