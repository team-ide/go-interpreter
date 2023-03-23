package parser

import (
	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/unistring"
	"github.com/team-ide/go-interpreter/syntax"
	"unicode/utf8"
)

type parser struct {
	syntax.Syntax
	str    string
	length int

	chr       rune // 当前 字符
	chrOffset int  // 当前 字符 偏移量
	offset    int  // 当前 字符 偏移量

	scope             *_scope
	insertSemicolon   bool // If we see a newline, then insert an implicit semicolon
	implicitSemicolon bool // An implicit semicolon exists

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
	declarationList []*ast.VariableDeclaration

	labels []unistring.String
}
