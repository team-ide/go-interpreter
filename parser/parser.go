package parser

import (
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/syntax"
	"github.com/team-ide/go-interpreter/token"
	"unicode/utf8"
)

type parser struct {
	syntax.Syntax
	str    string
	length int

	chr       rune // 当前 字符
	chrOffset int  // 当前 字符 偏移量
	offset    int  // 当前 字符 偏移量

	idx           int         // The index of token
	token         token.Token // The token
	literal       string      // The literal of the token, if any
	parsedLiteral node.String

	scope             *_scope
	insertSemicolon   bool // If we see a newline, then insert an implicit semicolon
	implicitSemicolon bool // An implicit semicolon exists

	recover struct {
		// Scratch when trying to seek to the next statement, etc.
		idx   int
		count int
	}

	errors ErrorList
}

// 隐式读取下一个
func (this_ *parser) implicitRead() rune {
	if this_.offset < this_.length {
		return rune(this_.str[this_.offset])
	}
	return -1
}

// 读取下一个 将重新设定偏移量
func (this_ *parser) read() {
	if this_.offset < this_.length {
		this_.chrOffset = this_.offset
		chr, width := rune(this_.str[this_.offset]), 1
		// 检查 编码 是否 是 ASCII
		if chr >= utf8.RuneSelf { // !ASCII
			chr, width = utf8.DecodeRuneInString(this_.str[this_.offset:])
			if chr == utf8.RuneError && width == 1 {
				_ = this_.error(this_.chrOffset, "Invalid UTF-8 character")
			}
		}
		this_.offset += width
		this_.chr = chr
	} else {
		this_.chrOffset = this_.length
		this_.chr = -1 // EOF 读取结束
	}
}

type _scope struct {
	outer           *_scope
	allowIn         bool
	allowLet        bool
	inIteration     bool
	inSwitch        bool
	inFuncParams    bool
	inFunction      bool
	inAsync         bool
	allowAwait      bool
	allowYield      bool
	declarationList []*node.VariableDeclaration

	labels []node.String
}

func (this_ *parser) openScope() {
	this_.scope = &_scope{
		outer:   this_.scope,
		allowIn: true,
	}
}

func (this_ *parser) closeScope() {
	this_.scope = this_.scope.outer
}

func (this_ *_scope) declare(declaration *node.VariableDeclaration) {
	this_.declarationList = append(this_.declarationList, declaration)
}

func (this_ *_scope) hasLabel(name node.String) bool {
	for _, label := range this_.labels {
		if label == name {
			return true
		}
	}
	if this_.outer != nil && !this_.inFunction {
		// Crossing a function boundary to look for a label is verboten
		return this_.outer.hasLabel(name)
	}
	return false
}
