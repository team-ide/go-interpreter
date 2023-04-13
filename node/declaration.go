package node

/** 声明 **/

// VariableDeclaration 变量声明
type VariableDeclaration struct {
	Var  int
	List []*Binding
}

// ClassElement 类元素
type ClassElement interface {
	Node
	isClassElement()
}

// FieldDefinition 字段定义
type FieldDefinition struct {
	Idx         int
	FieldNum    int
	Type        Expression
	Key         Expression
	Initializer Expression
	Computed    bool
	Static      bool
}

// IFaceDefinition 字段定义
type IFaceDefinition struct {
	From   int
	To     int
	Return Expression
	Name   Expression
	Params []*FieldDefinition
}

func (*IFaceDefinition) isExpression()    {}
func (*IFaceDefinition) isDefinition()    {}
func (this_ *IFaceDefinition) Start() int { return this_.From }
func (this_ *IFaceDefinition) End() int   { return this_.To }

// PropertyKind 属性类型
type PropertyKind string

const (
	PropertyKindValue  PropertyKind = "value"
	PropertyKindGet    PropertyKind = "get"
	PropertyKindSet    PropertyKind = "set"
	PropertyKindMethod PropertyKind = "method"
)

// MethodDefinition 方法定义
type MethodDefinition struct {
	Idx      int
	Key      Expression
	Kind     PropertyKind // "method", "get" or "set"
	Body     *FunctionLiteral
	Computed bool
	Static   bool
}

// ClassStaticBlock 类静态块
type ClassStaticBlock struct {
	Static          int
	Block           *BlockStatement
	Source          string
	DeclarationList []*VariableDeclaration
}

// ForLoopInitializer 循环初始化程序
type ForLoopInitializer interface {
	Node
	isForLoopInitializer()
}

// ForLoopInitializerExpression 循环初始化程序表达式
type ForLoopInitializerExpression struct {
	Expression Expression
}

// ForLoopInitializerVarDeclList 对于循环初始化程序变量声明列表
type ForLoopInitializerVarDeclList struct {
	Var  int
	List []*Binding
}

// ForLoopInitializerLexicalDecl 对于循环初始化程序词法声明
type ForLoopInitializerLexicalDecl struct {
	LexicalDeclaration LexicalDeclaration
}

// ForInto 循环
type ForInto interface {
	Node
	isForInto()
}

// ForIntoVar 循环变量
type ForIntoVar struct {
	Binding *Binding
}

// ForDeclaration 循环声明
type ForDeclaration struct {
	Idx     int
	IsConst bool
	Target  BindingTarget
}

// ForIntoExpression 循环表达式
type ForIntoExpression struct {
	Expression Expression
}

/* 实现 ForLoopInitializer 接口 */
func (*ForLoopInitializerExpression) isForLoopInitializer()  {}
func (*ForLoopInitializerVarDeclList) isForLoopInitializer() {}
func (*ForLoopInitializerLexicalDecl) isForLoopInitializer() {}

/* 实现 ForInto 接口 */
func (*ForIntoVar) isForInto()        {}
func (*ForDeclaration) isForInto()    {}
func (*ForIntoExpression) isForInto() {}

/* 实现 Pattern 接口 */
func (*ArrayPattern) isPattern() {}

/* 实现 BindingTarget 接口 */
func (*ArrayPattern) isBindingTarget() {}

/* 实现 Pattern 接口 */
func (*ObjectPattern) isPattern() {}

/* 实现 BindingTarget 接口 */
func (*ObjectPattern) isBindingTarget() {}

/* 实现 BindingTarget 接口 */
func (*BadExpression) isBindingTarget() {}

/* 实现 Property 接口 */
func (*PropertyShort) isProperty() {}
func (*PropertyKeyed) isProperty() {}
func (*SpreadElement) isProperty() {}

/* 实现 BindingTarget 接口 */
func (*Identifier) isBindingTarget() {}

/* 实现 ConciseBody 接口 */
func (*BlockStatement) isConciseBody() {}
func (*ExpressionBody) isConciseBody() {}

/* 实现 ClassElement 接口 */
func (*FieldDefinition) isClassElement()  {}
func (*MethodDefinition) isClassElement() {}
func (*ClassStaticBlock) isClassElement() {}

/* 实现 Node Start 接口 */

func (this_ *ForLoopInitializerExpression) Start() int  { return this_.Expression.Start() }
func (this_ *ForLoopInitializerVarDeclList) Start() int { return this_.List[0].Start() }
func (this_ *ForLoopInitializerLexicalDecl) Start() int {
	return this_.LexicalDeclaration.Start()
}
func (this_ *PropertyShort) Start() int  { return this_.Name.Idx }
func (this_ *PropertyKeyed) Start() int  { return this_.Key.Start() }
func (this_ *ExpressionBody) Start() int { return this_.Expression.Start() }

func (this_ *VariableDeclaration) Start() int { return this_.Var }
func (this_ *FieldDefinition) Start() int     { return this_.Idx }
func (this_ *MethodDefinition) Start() int    { return this_.Idx }
func (this_ *ClassStaticBlock) Start() int    { return this_.Static }

func (this_ *ForDeclaration) Start() int    { return this_.Idx }
func (this_ *ForIntoVar) Start() int        { return this_.Binding.Start() }
func (this_ *ForIntoExpression) Start() int { return this_.Expression.Start() }

/* 实现 Node End 接口 */

func (this_ *ForLoopInitializerExpression) End() int { return this_.Expression.End() }
func (this_ *ForLoopInitializerVarDeclList) End() int {
	return this_.List[len(this_.List)-1].End()
}
func (this_ *ForLoopInitializerLexicalDecl) End() int { return this_.LexicalDeclaration.End() }

func (this_ *PropertyShort) End() int {
	if this_.Initializer != nil {
		return this_.Initializer.End()
	}
	return this_.Name.End()
}

func (this_ *PropertyKeyed) End() int { return this_.Value.End() }

func (this_ *ExpressionBody) End() int { return this_.Expression.End() }

func (this_ *VariableDeclaration) End() int {
	if len(this_.List) > 0 {
		return this_.List[len(this_.List)-1].End()
	}

	return this_.Var + (+3)
}

func (this_ *FieldDefinition) End() int {
	if this_.Initializer != nil {
		return this_.Initializer.End()
	}
	if this_.Key != nil {
		return this_.Key.End()
	}
	return this_.Idx
}

func (this_ *MethodDefinition) End() int {
	return this_.Body.End()
}

func (this_ *ClassStaticBlock) End() int {
	return this_.Block.End()
}

func (this_ *YieldExpression) End() int {
	if this_.Argument != nil {
		return this_.Argument.End()
	}
	return this_.Yield + (+5)
}

func (this_ *ForDeclaration) End() int    { return this_.Target.End() }
func (this_ *ForIntoVar) End() int        { return this_.Binding.End() }
func (this_ *ForIntoExpression) End() int { return this_.Expression.End() }
