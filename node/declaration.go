package node

/** 声明 **/

type VariableDeclaration struct {
	Var  *Position
	List []*Binding
}

type ClassElement interface {
	Node
	isClassElement()
}

type FieldDefinition struct {
	Idx         *Position
	Key         Expression
	Initializer Expression
	Computed    bool
	Static      bool
}

type PropertyKind string

const (
	PropertyKindValue  PropertyKind = "value"
	PropertyKindGet    PropertyKind = "get"
	PropertyKindSet    PropertyKind = "set"
	PropertyKindMethod PropertyKind = "method"
)

type MethodDefinition struct {
	Idx      *Position
	Key      Expression
	Kind     PropertyKind // "method", "get" or "set"
	Body     *FunctionLiteral
	Computed bool
	Static   bool
}

type ClassStaticBlock struct {
	Static          *Position
	Block           *BlockStatement
	Source          string
	DeclarationList []*VariableDeclaration
}

type ForLoopInitializer interface {
	Node
	isForLoopInitializer()
}

type ForLoopInitializerExpression struct {
	Expression Expression
}

type ForLoopInitializerVarDeclList struct {
	Var  *Position
	List []*Binding
}

type ForLoopInitializerLexicalDecl struct {
	LexicalDeclaration LexicalDeclaration
}

type ForInto interface {
	Node
	isForInto()
}

type ForIntoVar struct {
	Binding *Binding
}

type ForDeclaration struct {
	Idx     *Position
	IsConst bool
	Target  BindingTarget
}

type ForIntoExpression struct {
	Expression Expression
}

func (*ForLoopInitializerExpression) isForLoopInitializer()  {}
func (*ForLoopInitializerVarDeclList) isForLoopInitializer() {}
func (*ForLoopInitializerLexicalDecl) isForLoopInitializer() {}

func (*ForIntoVar) isForInto()        {}
func (*ForDeclaration) isForInto()    {}
func (*ForIntoExpression) isForInto() {}

func (*ArrayPattern) isPattern()       {}
func (*ArrayPattern) isBindingTarget() {}

func (*ObjectPattern) isPattern()       {}
func (*ObjectPattern) isBindingTarget() {}

func (*BadExpression) isBindingTarget() {}

func (*PropertyShort) isProperty() {}
func (*PropertyKeyed) isProperty() {}
func (*SpreadElement) isProperty() {}

func (*Identifier) isBindingTarget() {}

func (*BlockStatement) isConciseBody() {}
func (*ExpressionBody) isConciseBody() {}

func (*FieldDefinition) isClassElement()  {}
func (*MethodDefinition) isClassElement() {}
func (*ClassStaticBlock) isClassElement() {}

func (this_ *ForLoopInitializerExpression) Start() *Position  { return this_.Expression.Start() }
func (this_ *ForLoopInitializerVarDeclList) Start() *Position { return this_.List[0].Start() }
func (this_ *ForLoopInitializerLexicalDecl) Start() *Position {
	return this_.LexicalDeclaration.Start()
}
func (this_ *PropertyShort) Start() *Position  { return this_.Name.Idx }
func (this_ *PropertyKeyed) Start() *Position  { return this_.Key.Start() }
func (this_ *ExpressionBody) Start() *Position { return this_.Expression.Start() }

func (this_ *VariableDeclaration) Start() *Position { return this_.Var }
func (this_ *FieldDefinition) Start() *Position     { return this_.Idx }
func (this_ *MethodDefinition) Start() *Position    { return this_.Idx }
func (this_ *ClassStaticBlock) Start() *Position    { return this_.Static }

func (this_ *ForDeclaration) Start() *Position    { return this_.Idx }
func (this_ *ForIntoVar) Start() *Position        { return this_.Binding.Start() }
func (this_ *ForIntoExpression) Start() *Position { return this_.Expression.Start() }

func (this_ *ForLoopInitializerExpression) End() *Position { return this_.Expression.End() }
func (this_ *ForLoopInitializerVarDeclList) End() *Position {
	return this_.List[len(this_.List)-1].End()
}
func (this_ *ForLoopInitializerLexicalDecl) End() *Position { return this_.LexicalDeclaration.End() }

func (this_ *PropertyShort) End() *Position {
	if this_.Initializer != nil {
		return this_.Initializer.End()
	}
	return this_.Name.End()
}

func (this_ *PropertyKeyed) End() *Position { return this_.Value.End() }

func (this_ *ExpressionBody) End() *Position { return this_.Expression.End() }

func (this_ *VariableDeclaration) End() *Position {
	if len(this_.List) > 0 {
		return this_.List[len(this_.List)-1].End()
	}

	//return this_.Var + 3
	return this_.Var.NewByColumnOffset(+3)
}

func (this_ *FieldDefinition) End() *Position {
	if this_.Initializer != nil {
		return this_.Initializer.End()
	}
	return this_.Key.End()
}

func (this_ *MethodDefinition) End() *Position {
	return this_.Body.End()
}

func (this_ *ClassStaticBlock) End() *Position {
	return this_.Block.End()
}

func (this_ *YieldExpression) End() *Position {
	if this_.Argument != nil {
		return this_.Argument.End()
	}
	//return this_.Yield + 5
	return this_.Yield.NewByColumnOffset(+5)
}

func (this_ *ForDeclaration) End() *Position    { return this_.Target.End() }
func (this_ *ForIntoVar) End() *Position        { return this_.Binding.End() }
func (this_ *ForIntoExpression) End() *Position { return this_.Expression.End() }
