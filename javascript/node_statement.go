package javascript

import "github.com/team-ide/go-interpreter/node"

// DebuggerStatement debugger语句
type DebuggerStatement struct {
	Debugger int
}

func (*DebuggerStatement) IsStatementNode() {}
func (this_ *DebuggerStatement) Start() int { return this_.Debugger }
func (this_ *DebuggerStatement) End() int   { return this_.Debugger + (8) }

// LabelledStatement 带标签的 语句
type LabelledStatement struct {
	Label     *node.Identifier
	Colon     int
	Statement node.Statement
}

func (*LabelledStatement) IsStatementNode() {}
func (this_ *LabelledStatement) Start() int { return this_.Label.Start() }
func (this_ *LabelledStatement) End() int   { return this_.Colon + (1) }
