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

type BindingTarget interface {
	Expression
	isBindingTarget()
}

type Binding struct {
	Target      BindingTarget
	Initializer Expression
}

type Pattern interface {
	BindingTarget
	isPattern()
}

type YieldExpression struct {
	Yield    *Position
	Argument Expression
	Delegate bool
}

type AwaitExpression struct {
	Await    *Position
	Argument Expression
}

type ArrayLiteral struct {
	LeftBracket  *Position
	RightBracket *Position
	Value        []Expression
}

type ArrayPattern struct {
	LeftBracket  *Position
	RightBracket *Position
	Elements     []Expression
	Rest         Expression
}

type AssignExpression struct {
	Operator token.Token
	Left     Expression
	Right    Expression
}

type BadExpression struct {
	From *Position
	To   *Position
}

type BinaryExpression struct {
	Operator   token.Token
	Left       Expression
	Right      Expression
	Comparison bool
}

type BooleanLiteral struct {
	Idx     *Position
	Literal string
	Value   bool
}

type BracketExpression struct {
	Left         Expression
	Member       Expression
	LeftBracket  *Position
	RightBracket *Position
}

type CallExpression struct {
	Callee           Expression
	LeftParenthesis  *Position
	ArgumentList     []Expression
	RightParenthesis *Position
}

type ConditionalExpression struct {
	Test       Expression
	Consequent Expression
	Alternate  Expression
}

type DotExpression struct {
	Left       Expression
	Identifier Identifier
}

type PrivateDotExpression struct {
	Left       Expression
	Identifier PrivateIdentifier
}
type OptionalChain struct {
	Expression
}
type Optional struct {
	Expression
}

type FunctionLiteral struct {
	Function      *Position
	Name          *Identifier
	ParameterList *ParameterList
	Body          *BlockStatement
	Source        string

	DeclarationList []*VariableDeclaration

	Async, Generator bool
}

type ClassLiteral struct {
	Class      *Position
	RightBrace *Position
	Name       *Identifier
	SuperClass Expression
	Body       []ClassElement
	Source     string
}

type ConciseBody interface {
	Node
	isConciseBody()
}

type ExpressionBody struct {
	Expression Expression
}

type ArrowFunctionLiteral struct {
	Start_          *Position
	ParameterList   *ParameterList
	Body            ConciseBody
	Source          string
	DeclarationList []*VariableDeclaration
	Async           bool
}

type Identifier struct {
	Name unistring.String
	Idx  *Position
}

type PrivateIdentifier struct {
	Identifier
}

type NewExpression struct {
	New              *Position
	Callee           Expression
	LeftParenthesis  *Position
	ArgumentList     []Expression
	RightParenthesis *Position
}

type NullLiteral struct {
	Idx     *Position
	Literal string
}

type NumberLiteral struct {
	Idx     *Position
	Literal string
	Value   interface{}
}

type ObjectLiteral struct {
	LeftBrace  *Position
	RightBrace *Position
	Value      []Property
}

type ObjectPattern struct {
	LeftBrace  *Position
	RightBrace *Position
	Properties []Property
	Rest       Expression
}

type ParameterList struct {
	Opening *Position
	List    []*Binding
	Rest    Expression
	Closing *Position
}

type Property interface {
	Expression
	isProperty()
}

type PropertyShort struct {
	Name        Identifier
	Initializer Expression
}

type PropertyKeyed struct {
	Key      Expression
	Kind     PropertyKind
	Value    Expression
	Computed bool
}

type SpreadElement struct {
	Expression
}

type RegExpLiteral struct {
	Idx     *Position
	Literal string
	Pattern string
	Flags   string
}

type SequenceExpression struct {
	Sequence []Expression
}

type StringLiteral struct {
	Idx     *Position
	Literal string
	Value   unistring.String
}

type TemplateElement struct {
	Idx     *Position
	Literal string
	Parsed  unistring.String
	Valid   bool
}

type TemplateLiteral struct {
	OpenQuote   *Position
	CloseQuote  *Position
	Tag         Expression
	Elements    []*TemplateElement
	Expressions []Expression
}

type ThisExpression struct {
	Idx *Position
}

type SuperExpression struct {
	Idx *Position
}

type UnaryExpression struct {
	Operator token.Token
	Idx      *Position // If a prefix operation
	Operand  Expression
	Postfix  bool
}

type MetaProperty struct {
	Meta, Property *Identifier
	Idx            *Position
}

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
