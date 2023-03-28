package parser

import (
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/syntax"
	"github.com/team-ide/go-interpreter/token"
	"unicode/utf8"
)

func New(src string, syntax syntax.Syntax) (p *Parser) {
	p = &Parser{
		Chr:    ' ',
		Str:    src,
		Length: len(src),
		Syntax: syntax,
		Base:   1,
	}
	return
}

type Parser struct {
	syntax.Syntax
	Str    string
	Length int

	Chr       rune // 当前 字符
	ChrOffset int  // 当前 字符 偏移量
	Offset    int  // 当前 字符 偏移量

	Base int

	Idx           int         // The index of token
	Token         token.Token // The token
	Literal       string      // The literal of the token, if any
	ParsedLiteral node.String

	Scope             *Scope
	InsertSemicolon   bool // If we see a newline, then insert an implicit semicolon
	ImplicitSemicolon bool // An implicit semicolon exists

	Recover struct {
		// Scratch when trying to seek to the next statement, etc.
		Idx   int
		Count int
	}

	Errors ErrorList
}

// ImplicitRead 隐式读取下一个
func (this_ *Parser) ImplicitRead() rune {
	if this_.Offset < this_.Length {
		return rune(this_.Str[this_.Offset])
	}
	return -1
}

// Read 读取下一个 将重新设定偏移量
func (this_ *Parser) Read() {
	if this_.Offset < this_.Length {
		this_.ChrOffset = this_.Offset
		chr, width := rune(this_.Str[this_.Offset]), 1
		// 检查 编码 是否 是 ASCII
		if chr >= utf8.RuneSelf { // !ASCII
			chr, width = utf8.DecodeRuneInString(this_.Str[this_.Offset:])
			if chr == utf8.RuneError && width == 1 {
				_ = this_.Error("read char utf8.RuneError chr:"+string(chr), this_.ChrOffset, "Invalid UTF-8 character")
			}
		}
		this_.Offset += width
		this_.Chr = chr
	} else {
		this_.ChrOffset = this_.Length
		this_.Chr = -1 // EOF 读取结束
	}
}

type Scope struct {
	Outer           *Scope
	AllowIn         bool
	AllowLet        bool
	InIteration     bool
	InSwitch        bool
	InFuncParams    bool
	InFunction      bool
	InAsync         bool
	AllowAwait      bool
	AllowYield      bool
	DeclarationList []*node.VariableDeclaration

	Labels []node.String
}

func (this_ *Parser) OpenScope() {
	this_.Scope = &Scope{
		Outer:   this_.Scope,
		AllowIn: true,
	}
}

func (this_ *Parser) CloseScope() {
	this_.Scope = this_.Scope.Outer
}

func (this_ *Scope) Declare(declaration *node.VariableDeclaration) {
	this_.DeclarationList = append(this_.DeclarationList, declaration)
}

func (this_ *Scope) HasLabel(name node.String) bool {
	for _, label := range this_.Labels {
		if label == name {
			return true
		}
	}
	if this_.Outer != nil && !this_.InFunction {
		// Crossing a function boundary to look for a label is verboten
		return this_.Outer.HasLabel(name)
	}
	return false
}
