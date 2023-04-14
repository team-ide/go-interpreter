package node

import (
	"github.com/team-ide/go-interpreter/token"
)

/** 语句 **/

// Statement 所有语句节点都实现了 Statement 接口
type Statement interface {
	Node
	IsStatementNode()
}

// BadStatement 错误的语句
type BadStatement struct {
	From int
	To   int
}

func (*BadStatement) IsStatementNode() {}
func (this_ *BadStatement) Start() int { return this_.From }
func (this_ *BadStatement) End() int   { return this_.To }

// BlankSpaceStatement 空白的语句
type BlankSpaceStatement struct {
	From int
	To   int
}

func (*BlankSpaceStatement) IsStatementNode() {}
func (this_ *BlankSpaceStatement) Start() int { return this_.From }
func (this_ *BlankSpaceStatement) End() int   { return this_.To }

// BlockStatement 块语句
type BlockStatement struct {
	LeftBrace  int
	List       []Statement
	RightBrace int
}

func (*BlockStatement) IsStatementNode()   {}
func (*BlockStatement) IsConciseBodyNode() {}
func (this_ *BlockStatement) Start() int   { return this_.LeftBrace }
func (this_ *BlockStatement) End() int     { return this_.RightBrace + (1) }

// BreakStatement 跳出语句
type BreakStatement struct {
	From  int
	To    int
	Label *Identifier
}

func (*BreakStatement) IsStatementNode() {}
func (this_ *BreakStatement) Start() int { return this_.From }
func (this_ *BreakStatement) End() int   { return this_.To }

// ContinueStatement 继续语句
type ContinueStatement struct {
	From  int
	To    int
	Label *Identifier
}

func (*ContinueStatement) IsStatementNode() {}
func (this_ *ContinueStatement) Start() int { return this_.From }
func (this_ *ContinueStatement) End() int   { return this_.To }

// CaseStatement case语句
type CaseStatement struct {
	Case       int
	Test       Expression
	Consequent []Statement
}

func (*CaseStatement) IsStatementNode() {}
func (this_ *CaseStatement) Start() int { return this_.Case }
func (this_ *CaseStatement) End() int   { return this_.Consequent[len(this_.Consequent)-1].End() }

// CatchStatement cache快
type CatchStatement struct {
	Catch     int
	Parameter BindingTarget
	Body      *BlockStatement
}

func (*CatchStatement) IsStatementNode() {}
func (this_ *CatchStatement) Start() int { return this_.Catch }
func (this_ *CatchStatement) End() int   { return this_.Body.End() }

// DoWhileStatement do while 语句
type DoWhileStatement struct {
	Do     int
	EndIdx int
	Test   Expression
	Body   Statement
}

func (*DoWhileStatement) IsStatementNode() {}
func (this_ *DoWhileStatement) Start() int { return this_.Do }
func (this_ *DoWhileStatement) End() int   { return this_.EndIdx }

// SemicolonStatement 分号语句
type SemicolonStatement struct {
	Semicolon int
}

func (*SemicolonStatement) IsStatementNode() {}
func (this_ *SemicolonStatement) Start() int { return this_.Semicolon }
func (this_ *SemicolonStatement) End() int   { return this_.Semicolon + (1) }

// ExpressionStatement 表达式语句
type ExpressionStatement struct {
	Expression Expression
}

func (*ExpressionStatement) IsStatementNode() {}
func (this_ *ExpressionStatement) Start() int { return this_.Expression.Start() }
func (this_ *ExpressionStatement) End() int   { return this_.Expression.End() }

// ForInStatement for in 语句
type ForInStatement struct {
	For    int
	Into   ForInto
	Source Expression
	Body   Statement
}

func (*ForInStatement) IsStatementNode() {}
func (this_ *ForInStatement) Start() int { return this_.For }
func (this_ *ForInStatement) End() int   { return this_.Body.End() }

// ForOfStatement for of 语句
type ForOfStatement struct {
	For    int
	Into   ForInto
	Source Expression
	Body   Statement
}

func (*ForOfStatement) IsStatementNode() {}
func (this_ *ForOfStatement) Start() int { return this_.For }
func (this_ *ForOfStatement) End() int   { return this_.Body.End() }

// ForStatement for 语句
type ForStatement struct {
	For         int
	Initializer ForLoopInitializer
	Update      Expression
	Test        Expression
	Body        Statement
}

func (*ForStatement) IsStatementNode() {}
func (this_ *ForStatement) Start() int { return this_.For }
func (this_ *ForStatement) End() int   { return this_.Body.End() }

// IfStatement if 语句
type IfStatement struct {
	If         int
	Test       Expression
	Consequent Statement
	Alternate  Statement
}

func (*IfStatement) IsStatementNode() {}
func (this_ *IfStatement) Start() int { return this_.If }
func (this_ *IfStatement) End() int {
	if this_.Alternate != nil {
		return this_.Alternate.End()
	}
	return this_.Consequent.End()
}

// ReturnStatement 返回 语句
type ReturnStatement struct {
	Return   int
	Argument Expression
}

func (*ReturnStatement) IsStatementNode() {}
func (this_ *ReturnStatement) Start() int { return this_.Return }
func (this_ *ReturnStatement) End() int   { return this_.Return + (6) }

// SwitchStatement switch 语句
type SwitchStatement struct {
	Switch       int
	Discriminant Expression
	Default      int
	Body         []*CaseStatement
}

func (*SwitchStatement) IsStatementNode() {}
func (this_ *SwitchStatement) Start() int { return this_.Switch }
func (this_ *SwitchStatement) End() int   { return this_.Body[len(this_.Body)-1].End() }

// ThrowStatement throw 语句
type ThrowStatement struct {
	Throw    int
	Argument Expression
}

func (*ThrowStatement) IsStatementNode() {}
func (this_ *ThrowStatement) Start() int { return this_.Throw }
func (this_ *ThrowStatement) End() int   { return this_.Argument.End() }

// TryStatement try 语句
type TryStatement struct {
	Try     int
	Body    *BlockStatement
	Catch   *CatchStatement
	Finally *BlockStatement
}

func (*TryStatement) IsStatementNode() {}
func (this_ *TryStatement) Start() int { return this_.Try }
func (this_ *TryStatement) End() int {
	if this_.Finally != nil {
		return this_.Finally.End()
	}
	if this_.Catch != nil {
		return this_.Catch.End()
	}
	return this_.Body.End()
}

// VariableStatement 变量 语句
type VariableStatement struct {
	Var  int
	List []*Binding
}

func (*VariableStatement) IsStatementNode() {}
func (this_ *VariableStatement) Start() int { return this_.Var }
func (this_ *VariableStatement) End() int   { return this_.List[len(this_.List)-1].End() }

// LexicalDeclaration 词汇声明
type LexicalDeclaration struct {
	Idx   int
	Token token.Token
	List  []*Binding
}

func (*LexicalDeclaration) IsStatementNode() {}
func (this_ *LexicalDeclaration) Start() int { return this_.Idx }
func (this_ *LexicalDeclaration) End() int   { return this_.List[len(this_.List)-1].End() }

// WhileStatement while 语句
type WhileStatement struct {
	While int
	Test  Expression
	Body  Statement
}

func (*WhileStatement) IsStatementNode() {}
func (this_ *WhileStatement) Start() int { return this_.While }
func (this_ *WhileStatement) End() int   { return this_.Body.End() }

// WithStatement with 语句
type WithStatement struct {
	With   int
	Object Expression
	Body   Statement
}

func (*WithStatement) IsStatementNode() {}
func (this_ *WithStatement) Start() int { return this_.With }
func (this_ *WithStatement) End() int   { return this_.Body.End() }

// FunctionDeclaration 函数 声明
type FunctionDeclaration struct {
	Function *FunctionLiteral
}

func (*FunctionDeclaration) IsStatementNode() {}
func (this_ *FunctionDeclaration) Start() int { return this_.Function.Start() }
func (this_ *FunctionDeclaration) End() int   { return this_.Function.End() }

// ClassDeclaration 类 声明
type ClassDeclaration struct {
	Class *ClassLiteral
}

func (*ClassDeclaration) IsStatementNode() {}
func (this_ *ClassDeclaration) Start() int { return this_.Class.Start() }
func (this_ *ClassDeclaration) End() int   { return this_.Class.End() }
