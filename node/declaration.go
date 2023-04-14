package node

/** 声明 **/

// VariableDeclaration 变量声明
type VariableDeclaration struct {
	Var  int
	List []*Binding
}

func (this_ *VariableDeclaration) Start() int { return this_.Var }
func (this_ *VariableDeclaration) End() int {
	if len(this_.List) > 0 {
		return this_.List[len(this_.List)-1].End()
	}
	return this_.Var + (+3)
}

// ClassElement 类元素
type ClassElement interface {
	Node
	IsClassElementNode()
}

// FieldDefinition 字段定义
type FieldDefinition struct {
	Idx         int
	Key         Expression
	Initializer Expression
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

func (*MethodDefinition) IsClassElementNode() {}
func (this_ *MethodDefinition) Start() int    { return this_.Idx }
func (this_ *MethodDefinition) End() int      { return this_.Body.End() }

// ClassStaticBlock 类静态块
type ClassStaticBlock struct {
	Static          int
	Block           *BlockStatement
	Source          string
	DeclarationList []*VariableDeclaration
}

func (*ClassStaticBlock) IsClassElementNode() {}
func (this_ *ClassStaticBlock) Start() int    { return this_.Static }
func (this_ *ClassStaticBlock) End() int      { return this_.Block.End() }

// ForLoopInitializer 循环初始化程序
type ForLoopInitializer interface {
	Node
	IsForLoopInitializerNode()
}

// ForLoopInitializerExpression 循环初始化程序表达式
type ForLoopInitializerExpression struct {
	Expression Expression
}

func (*ForLoopInitializerExpression) IsForLoopInitializerNode() {}
func (this_ *ForLoopInitializerExpression) Start() int          { return this_.Expression.Start() }
func (this_ *ForLoopInitializerExpression) End() int            { return this_.Expression.End() }

// ForLoopInitializerVarDeclList 对于循环初始化程序变量声明列表
type ForLoopInitializerVarDeclList struct {
	Var  int
	List []*Binding
}

func (*ForLoopInitializerVarDeclList) IsForLoopInitializerNode() {}
func (this_ *ForLoopInitializerVarDeclList) Start() int          { return this_.List[0].Start() }
func (this_ *ForLoopInitializerVarDeclList) End() int {
	return this_.List[len(this_.List)-1].End()
}

// ForLoopInitializerLexicalDecl 对于循环初始化程序词法声明
type ForLoopInitializerLexicalDecl struct {
	LexicalDeclaration LexicalDeclaration
}

func (*ForLoopInitializerLexicalDecl) IsForLoopInitializerNode() {}
func (this_ *ForLoopInitializerLexicalDecl) Start() int {
	return this_.LexicalDeclaration.Start()
}
func (this_ *ForLoopInitializerLexicalDecl) End() int { return this_.LexicalDeclaration.End() }

// ForInto 循环
type ForInto interface {
	Node
	IsForIntoNode()
}

// ForIntoVar 循环变量
type ForIntoVar struct {
	Binding *Binding
}

func (*ForIntoVar) IsForIntoNode()   {}
func (this_ *ForIntoVar) Start() int { return this_.Binding.Start() }
func (this_ *ForIntoVar) End() int   { return this_.Binding.End() }

// ForDeclaration 循环声明
type ForDeclaration struct {
	Idx     int
	IsConst bool
	Target  BindingTarget
}

func (*ForDeclaration) IsForIntoNode()   {}
func (this_ *ForDeclaration) Start() int { return this_.Idx }
func (this_ *ForDeclaration) End() int   { return this_.Target.End() }

// ForIntoExpression 循环表达式
type ForIntoExpression struct {
	Expression Expression
}

func (*ForIntoExpression) IsForIntoNode()   {}
func (this_ *ForIntoExpression) Start() int { return this_.Expression.Start() }
func (this_ *ForIntoExpression) End() int   { return this_.Expression.End() }
