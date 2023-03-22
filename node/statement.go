package node

import (
	"github.com/team-ide/go-interpreter/token"
)

/**表达式**/

// Statement 所有语句节点都实现了 Statement 接口
type Statement interface {
	Node
	isStatement()
}
type BadStatement struct {
	From *Position
	To   *Position
}

type BlockStatement struct {
	LeftBrace  *Position
	List       []Statement
	RightBrace *Position
}

type BranchStatement struct {
	Idx   *Position
	Token token.Token
	Label *Identifier
}

type CaseStatement struct {
	Case       *Position
	Test       Expression
	Consequent []Statement
}

type CatchStatement struct {
	Catch     *Position
	Parameter BindingTarget
	Body      *BlockStatement
}

type DebuggerStatement struct {
	Debugger *Position
}

type DoWhileStatement struct {
	Do   *Position
	Test Expression
	Body Statement
}

type EmptyStatement struct {
	Semicolon *Position
}

type ExpressionStatement struct {
	Expression Expression
}

type ForInStatement struct {
	For    *Position
	Into   ForInto
	Source Expression
	Body   Statement
}

type ForOfStatement struct {
	For    *Position
	Into   ForInto
	Source Expression
	Body   Statement
}

type ForStatement struct {
	For         *Position
	Initializer ForLoopInitializer
	Update      Expression
	Test        Expression
	Body        Statement
}

type IfStatement struct {
	If         *Position
	Test       Expression
	Consequent Statement
	Alternate  Statement
}

type LabelledStatement struct {
	Label     *Identifier
	Colon     *Position
	Statement Statement
}

type ReturnStatement struct {
	Return   *Position
	Argument Expression
}

type SwitchStatement struct {
	Switch       *Position
	Discriminant Expression
	Default      int
	Body         []*CaseStatement
}

type ThrowStatement struct {
	Throw    *Position
	Argument Expression
}

type TryStatement struct {
	Try     *Position
	Body    *BlockStatement
	Catch   *CatchStatement
	Finally *BlockStatement
}

type VariableStatement struct {
	Var  *Position
	List []*Binding
}

type LexicalDeclaration struct {
	Idx   *Position
	Token token.Token
	List  []*Binding
}

type WhileStatement struct {
	While *Position
	Test  Expression
	Body  Statement
}

type WithStatement struct {
	With   *Position
	Object Expression
	Body   Statement
}

type FunctionDeclaration struct {
	Function *FunctionLiteral
}

type ClassDeclaration struct {
	Class *ClassLiteral
}

func (*BadStatement) isStatement()        {}
func (*BlockStatement) isStatement()      {}
func (*BranchStatement) isStatement()     {}
func (*CaseStatement) isStatement()       {}
func (*CatchStatement) isStatement()      {}
func (*DebuggerStatement) isStatement()   {}
func (*DoWhileStatement) isStatement()    {}
func (*EmptyStatement) isStatement()      {}
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

func (this_ *BadStatement) Start() *Position        { return this_.From }
func (this_ *BlockStatement) Start() *Position      { return this_.LeftBrace }
func (this_ *BranchStatement) Start() *Position     { return this_.Idx }
func (this_ *CaseStatement) Start() *Position       { return this_.Case }
func (this_ *CatchStatement) Start() *Position      { return this_.Catch }
func (this_ *DebuggerStatement) Start() *Position   { return this_.Debugger }
func (this_ *DoWhileStatement) Start() *Position    { return this_.Do }
func (this_ *EmptyStatement) Start() *Position      { return this_.Semicolon }
func (this_ *ExpressionStatement) Start() *Position { return this_.Expression.Start() }
func (this_ *ForInStatement) Start() *Position      { return this_.For }
func (this_ *ForOfStatement) Start() *Position      { return this_.For }
func (this_ *ForStatement) Start() *Position        { return this_.For }
func (this_ *IfStatement) Start() *Position         { return this_.If }
func (this_ *LabelledStatement) Start() *Position   { return this_.Label.Start() }
func (this_ *ReturnStatement) Start() *Position     { return this_.Return }
func (this_ *SwitchStatement) Start() *Position     { return this_.Switch }
func (this_ *ThrowStatement) Start() *Position      { return this_.Throw }
func (this_ *TryStatement) Start() *Position        { return this_.Try }
func (this_ *VariableStatement) Start() *Position   { return this_.Var }
func (this_ *WhileStatement) Start() *Position      { return this_.While }
func (this_ *WithStatement) Start() *Position       { return this_.With }
func (this_ *LexicalDeclaration) Start() *Position  { return this_.Idx }
func (this_ *FunctionDeclaration) Start() *Position { return this_.Function.Start() }
func (this_ *ClassDeclaration) Start() *Position    { return this_.Class.Start() }
func (this_ *Binding) Start() *Position             { return this_.Target.Start() }

func (this_ *BadStatement) End() *Position        { return this_.To }
func (this_ *BlockStatement) End() *Position      { return this_.RightBrace.NewByColumnOffset(1) }
func (this_ *BranchStatement) End() *Position     { return this_.Idx }
func (this_ *CaseStatement) End() *Position       { return this_.Consequent[len(this_.Consequent)-1].End() }
func (this_ *CatchStatement) End() *Position      { return this_.Body.End() }
func (this_ *DebuggerStatement) End() *Position   { return this_.Debugger.NewByColumnOffset(8) }
func (this_ *DoWhileStatement) End() *Position    { return this_.Test.End() }
func (this_ *EmptyStatement) End() *Position      { return this_.Semicolon.NewByColumnOffset(1) }
func (this_ *ExpressionStatement) End() *Position { return this_.Expression.End() }
func (this_ *ForInStatement) End() *Position      { return this_.Body.End() }
func (this_ *ForOfStatement) End() *Position      { return this_.Body.End() }
func (this_ *ForStatement) End() *Position        { return this_.Body.End() }
func (this_ *IfStatement) End() *Position {
	if this_.Alternate != nil {
		return this_.Alternate.End()
	}
	return this_.Consequent.End()
}
func (this_ *LabelledStatement) End() *Position { return this_.Colon.NewByColumnOffset(1) }
func (this_ *ReturnStatement) End() *Position   { return this_.Return.NewByColumnOffset(6) }
func (this_ *SwitchStatement) End() *Position   { return this_.Body[len(this_.Body)-1].End() }
func (this_ *ThrowStatement) End() *Position    { return this_.Argument.End() }
func (this_ *TryStatement) End() *Position {
	if this_.Finally != nil {
		return this_.Finally.End()
	}
	if this_.Catch != nil {
		return this_.Catch.End()
	}
	return this_.Body.End()
}
func (this_ *VariableStatement) End() *Position   { return this_.List[len(this_.List)-1].End() }
func (this_ *WhileStatement) End() *Position      { return this_.Body.End() }
func (this_ *WithStatement) End() *Position       { return this_.Body.End() }
func (this_ *LexicalDeclaration) End() *Position  { return this_.List[len(this_.List)-1].End() }
func (this_ *FunctionDeclaration) End() *Position { return this_.Function.End() }
func (this_ *ClassDeclaration) End() *Position    { return this_.Class.End() }
func (this_ *Binding) End() *Position {
	if this_.Initializer != nil {
		return this_.Initializer.End()
	}
	return this_.Target.End()
}
