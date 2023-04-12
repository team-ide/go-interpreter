package node

import (
	"github.com/team-ide/go-interpreter/token"
)

/** 语句 **/

// Statement 所有语句节点都实现了 Statement 接口
type Statement interface {
	Node
	isStatement()
}

// BlankSpaceStatement 空白的语句
type BlankSpaceStatement struct {
	From int
	To   int
}

// BadStatement 错误的语句
type BadStatement struct {
	From int
	To   int
}

// BlockStatement block语句
type BlockStatement struct {
	LeftBrace  int
	List       []Statement
	RightBrace int
}

// ImportStatement 导入
type ImportStatement struct {
	From   int
	To     int
	Import string
}

// IncludeStatement 导入
type IncludeStatement struct {
	From    int
	To      int
	Include string
}

func (*IncludeStatement) isStatement()     {}
func (this_ *IncludeStatement) Start() int { return this_.From }
func (this_ *IncludeStatement) End() int   { return this_.To }

// NamespaceStatement 导入
type NamespaceStatement struct {
	From      int
	To        int
	Language  string
	Namespace string
}

func (*NamespaceStatement) isStatement()     {}
func (this_ *NamespaceStatement) Start() int { return this_.From }
func (this_ *NamespaceStatement) End() int   { return this_.To }

// ExceptionStatement 导入
type ExceptionStatement struct {
	From int
	To   int
	Name string
}

func (*ExceptionStatement) isStatement()     {}
func (this_ *ExceptionStatement) Start() int { return this_.From }
func (this_ *ExceptionStatement) End() int   { return this_.To }

// StructStatement 导入
type StructStatement struct {
	From int
	To   int
	Name string
}

func (*StructStatement) isStatement()     {}
func (this_ *StructStatement) Start() int { return this_.From }
func (this_ *StructStatement) End() int   { return this_.To }

// EnumStatement 导入
type EnumStatement struct {
	From int
	To   int
	Name string
}

func (*EnumStatement) isStatement()     {}
func (this_ *EnumStatement) Start() int { return this_.From }
func (this_ *EnumStatement) End() int   { return this_.To }

// ServiceStatement 导入
type ServiceStatement struct {
	From int
	To   int
	Name string
}

func (*ServiceStatement) isStatement()     {}
func (this_ *ServiceStatement) Start() int { return this_.From }
func (this_ *ServiceStatement) End() int   { return this_.To }

// BranchStatement branch语句
type BranchStatement struct {
	Idx    int
	EndIdx int
	Token  token.Token
	Label  *Identifier
}

// CaseStatement case语句
type CaseStatement struct {
	Case       int
	Test       Expression
	Consequent []Statement
}

// CatchStatement cache快
type CatchStatement struct {
	Catch     int
	Parameter BindingTarget
	Body      *BlockStatement
}

// DebuggerStatement debugger语句
type DebuggerStatement struct {
	Debugger int
}

// DoWhileStatement do while 语句
type DoWhileStatement struct {
	Do     int
	EndIdx int
	Test   Expression
	Body   Statement
}

// SemicolonStatement 分号语句
type SemicolonStatement struct {
	Semicolon int
}

// ExpressionStatement 表达式语句
type ExpressionStatement struct {
	Expression Expression
}

// ForInStatement for in 语句
type ForInStatement struct {
	For    int
	Into   ForInto
	Source Expression
	Body   Statement
}

// ForOfStatement for of 语句
type ForOfStatement struct {
	For    int
	Into   ForInto
	Source Expression
	Body   Statement
}

// ForStatement for 语句
type ForStatement struct {
	For         int
	Initializer ForLoopInitializer
	Update      Expression
	Test        Expression
	Body        Statement
}

// IfStatement if 语句
type IfStatement struct {
	If         int
	Test       Expression
	Consequent Statement
	Alternate  Statement
}

// LabelledStatement 带标签的 语句
type LabelledStatement struct {
	Label     *Identifier
	Colon     int
	Statement Statement
}

// ReturnStatement 返回 语句
type ReturnStatement struct {
	Return   int
	Argument Expression
}

// SwitchStatement switch 语句
type SwitchStatement struct {
	Switch       int
	Discriminant Expression
	Default      int
	Body         []*CaseStatement
}

// ThrowStatement throw 语句
type ThrowStatement struct {
	Throw    int
	Argument Expression
}

// TryStatement try 语句
type TryStatement struct {
	Try     int
	Body    *BlockStatement
	Catch   *CatchStatement
	Finally *BlockStatement
}

// VariableStatement 变量 语句
type VariableStatement struct {
	Var  int
	List []*Binding
}

// LexicalDeclaration 词汇声明
type LexicalDeclaration struct {
	Idx   int
	Token token.Token
	List  []*Binding
}

// WhileStatement while 语句
type WhileStatement struct {
	While int
	Test  Expression
	Body  Statement
}

// WithStatement with 语句
type WithStatement struct {
	With   int
	Object Expression
	Body   Statement
}

// FunctionDeclaration 函数 声明
type FunctionDeclaration struct {
	Function *FunctionLiteral
}

// ClassDeclaration 类 声明
type ClassDeclaration struct {
	Class *ClassLiteral
}

/* 实现 Statement 接口 */
func (*BadStatement) isStatement()        {}
func (*ImportStatement) isStatement()     {}
func (*BlankSpaceStatement) isStatement() {}
func (*BlockStatement) isStatement()      {}
func (*BranchStatement) isStatement()     {}
func (*CaseStatement) isStatement()       {}
func (*CatchStatement) isStatement()      {}
func (*DebuggerStatement) isStatement()   {}
func (*DoWhileStatement) isStatement()    {}
func (*SemicolonStatement) isStatement()  {}
func (*ExpressionStatement) isStatement() {}
func (*ForInStatement) isStatement()      {}
func (*ForOfStatement) isStatement()      {}
func (*ForStatement) isStatement()        {}
func (*IfStatement) isStatement()         {}
func (*LabelledStatement) isStatement()   {}
func (*ReturnStatement) isStatement()     {}
func (*SwitchStatement) isStatement()     {}
func (*ThrowStatement) isStatement()      {}
func (*TryStatement) isStatement()        {}
func (*VariableStatement) isStatement()   {}
func (*WhileStatement) isStatement()      {}
func (*WithStatement) isStatement()       {}
func (*LexicalDeclaration) isStatement()  {}
func (*FunctionDeclaration) isStatement() {}
func (*ClassDeclaration) isStatement()    {}

/* 实现 Node Start 接口 */

func (this_ *BadStatement) Start() int        { return this_.From }
func (this_ *ImportStatement) Start() int     { return this_.From }
func (this_ *BlankSpaceStatement) Start() int { return this_.From }
func (this_ *BlockStatement) Start() int      { return this_.LeftBrace }
func (this_ *BranchStatement) Start() int     { return this_.Idx }
func (this_ *CaseStatement) Start() int       { return this_.Case }
func (this_ *CatchStatement) Start() int      { return this_.Catch }
func (this_ *DebuggerStatement) Start() int   { return this_.Debugger }
func (this_ *DoWhileStatement) Start() int    { return this_.Do }
func (this_ *SemicolonStatement) Start() int  { return this_.Semicolon }
func (this_ *ExpressionStatement) Start() int { return this_.Expression.Start() }
func (this_ *ForInStatement) Start() int      { return this_.For }
func (this_ *ForOfStatement) Start() int      { return this_.For }
func (this_ *ForStatement) Start() int        { return this_.For }
func (this_ *IfStatement) Start() int         { return this_.If }
func (this_ *LabelledStatement) Start() int   { return this_.Label.Start() }
func (this_ *ReturnStatement) Start() int     { return this_.Return }
func (this_ *SwitchStatement) Start() int     { return this_.Switch }
func (this_ *ThrowStatement) Start() int      { return this_.Throw }
func (this_ *TryStatement) Start() int        { return this_.Try }
func (this_ *VariableStatement) Start() int   { return this_.Var }
func (this_ *WhileStatement) Start() int      { return this_.While }
func (this_ *WithStatement) Start() int       { return this_.With }
func (this_ *LexicalDeclaration) Start() int  { return this_.Idx }
func (this_ *FunctionDeclaration) Start() int { return this_.Function.Start() }
func (this_ *ClassDeclaration) Start() int    { return this_.Class.Start() }
func (this_ *Binding) Start() int             { return this_.Target.Start() }

/* 实现 Node End 接口 */

func (this_ *BadStatement) End() int        { return this_.To }
func (this_ *ImportStatement) End() int     { return this_.To }
func (this_ *BlankSpaceStatement) End() int { return this_.To }
func (this_ *BlockStatement) End() int      { return this_.RightBrace + (1) }
func (this_ *BranchStatement) End() int {
	return this_.EndIdx
}
func (this_ *CaseStatement) End() int       { return this_.Consequent[len(this_.Consequent)-1].End() }
func (this_ *CatchStatement) End() int      { return this_.Body.End() }
func (this_ *DebuggerStatement) End() int   { return this_.Debugger + (8) }
func (this_ *DoWhileStatement) End() int    { return this_.EndIdx }
func (this_ *SemicolonStatement) End() int  { return this_.Semicolon + (1) }
func (this_ *ExpressionStatement) End() int { return this_.Expression.End() }
func (this_ *ForInStatement) End() int      { return this_.Body.End() }
func (this_ *ForOfStatement) End() int      { return this_.Body.End() }
func (this_ *ForStatement) End() int        { return this_.Body.End() }
func (this_ *IfStatement) End() int {
	if this_.Alternate != nil {
		return this_.Alternate.End()
	}
	return this_.Consequent.End()
}
func (this_ *LabelledStatement) End() int { return this_.Colon + (1) }
func (this_ *ReturnStatement) End() int   { return this_.Return + (6) }
func (this_ *SwitchStatement) End() int   { return this_.Body[len(this_.Body)-1].End() }
func (this_ *ThrowStatement) End() int    { return this_.Argument.End() }
func (this_ *TryStatement) End() int {
	if this_.Finally != nil {
		return this_.Finally.End()
	}
	if this_.Catch != nil {
		return this_.Catch.End()
	}
	return this_.Body.End()
}
func (this_ *VariableStatement) End() int   { return this_.List[len(this_.List)-1].End() }
func (this_ *WhileStatement) End() int      { return this_.Body.End() }
func (this_ *WithStatement) End() int       { return this_.Body.End() }
func (this_ *LexicalDeclaration) End() int  { return this_.List[len(this_.List)-1].End() }
func (this_ *FunctionDeclaration) End() int { return this_.Function.End() }
func (this_ *ClassDeclaration) End() int    { return this_.Class.End() }
func (this_ *Binding) End() int {
	if this_.Initializer != nil {
		return this_.Initializer.End()
	}
	return this_.Target.End()
}
