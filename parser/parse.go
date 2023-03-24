package parser

import (
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/syntax"
	"github.com/team-ide/go-interpreter/token"
	"unicode"
	"unicode/utf8"
)

func Parse(src string, syntax syntax.Syntax) (tree *node.Tree, err error) {
	p := &parser{
		chr:    ' ',
		str:    src,
		length: len(src),
		Syntax: syntax,
	}
	return p.parse()
}

// 解析
func (this_ *parser) parse() (tree *node.Tree, err error) {
	this_.openScope()
	defer this_.closeScope()

	this_.next()
	tree = this_.parseTree()
	//this_.errors.Sort()
	err = this_.errors.Err()
	return
}

func (this_ *parser) parseTree() (tree *node.Tree) {
	tree = &node.Tree{
		Children:        this_.parseSourceElements(),
		DeclarationList: this_.scope.declarationList,
	}
	//this_.file.SetSourceMap(this_.parseSourceMap())
	return
}

func (this_ *parser) parseSourceElements() (statements []node.Statement) {
	for this_.token != token.Eof {
		this_.scope.allowLet = true
		statements = append(statements, this_.parseStatement())
	}

	return statements
}

func (this_ *parser) skipWhiteSpace() {
	for {
		switch this_.chr {
		case ' ', '\t', '\f', '\v', '\u00a0', '\ufeff':
			this_.read()
			continue
		case '\r':
			if this_.implicitRead() == '\n' {
				this_.read()
			}
			fallthrough
		case '\u2028', '\u2029', '\n':
			if this_.insertSemicolon {
				return
			}
			this_.read()
			continue
		}
		if this_.chr >= utf8.RuneSelf {
			if unicode.IsSpace(this_.chr) {
				this_.read()
				continue
			}
		}
		break
	}
}

func (this_ *parser) Position(offset int) (position *node.Position) {
	position = &node.Position{
		Idx: offset,
	}
	return
}

func (this_ *parser) next() {
	this_.token, this_.literal, this_.parsedLiteral, this_.idx = this_.scan()
}

func (this_ *parser) slice(start, end int) string {
	from := start
	to := end
	if from >= 0 && to <= len(this_.str) {
		return this_.str[from:to]
	}

	return ""
}

func (this_ *parser) optionalSemicolon() {
	if this_.token == token.Semicolon {
		this_.next()
		return
	}

	if this_.implicitSemicolon {
		this_.implicitSemicolon = false
		return
	}

	if this_.token != token.Eof && this_.token != token.RightBrace {
		this_.expect(token.Semicolon)
	}
}

func (this_ *parser) semicolon() {
	if this_.token != token.RightParenthesis && this_.token != token.RightBrace {
		if this_.implicitSemicolon {
			this_.implicitSemicolon = false
			return
		}

		this_.expect(token.Semicolon)
	}
}

type parserState struct {
	idx                                int
	tok                                token.Token
	literal                            string
	parsedLiteral                      node.String
	implicitSemicolon, insertSemicolon bool
	chr                                rune
	chrOffset, offset                  int
	errorCount                         int
}

func (this_ *parser) mark(state *parserState) *parserState {
	if state == nil {
		state = &parserState{}
	}
	state.idx, state.tok, state.literal, state.parsedLiteral, state.implicitSemicolon, state.insertSemicolon, state.chr, state.chrOffset, state.offset =
		this_.idx, this_.token, this_.literal, this_.parsedLiteral, this_.implicitSemicolon, this_.insertSemicolon, this_.chr, this_.chrOffset, this_.offset

	state.errorCount = len(this_.errors)
	return state
}

func (this_ *parser) restore(state *parserState) {
	this_.idx, this_.token, this_.literal, this_.parsedLiteral, this_.implicitSemicolon, this_.insertSemicolon, this_.chr, this_.chrOffset, this_.offset =
		state.idx, state.tok, state.literal, state.parsedLiteral, state.implicitSemicolon, state.insertSemicolon, state.chr, state.chrOffset, state.offset
	this_.errors = this_.errors[:state.errorCount]
}

func (this_ *parser) peek() token.Token {
	implicitSemicolon, insertSemicolon, chr, chrOffset, offset := this_.implicitSemicolon, this_.insertSemicolon, this_.chr, this_.chrOffset, this_.offset
	tok, _, _, _ := this_.scan()
	this_.implicitSemicolon, this_.insertSemicolon, this_.chr, this_.chrOffset, this_.offset = implicitSemicolon, insertSemicolon, chr, chrOffset, offset
	return tok
}
