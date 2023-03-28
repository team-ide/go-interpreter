package parser

import (
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/token"
)

func (this_ *Parser) Position(offset int) (position *node.Position) {
	position = &node.Position{
		Idx: offset - this_.Base,
	}
	return
}

func (this_ *Parser) Next() {
	this_.Token, this_.Literal, this_.ParsedLiteral, this_.Idx = this_.Scan()
}

func (this_ *Parser) IdxOf(offset int) int {
	return this_.Base + offset
}
func (this_ *Parser) Slice(start, end int) string {
	from := start - this_.Base
	to := end - this_.Base
	//if from >= 0 && to <= len(this_.str) {
	return this_.Str[from:to]
	//}

	//return ""
}

func (this_ *Parser) OptionalSemicolon() {
	if this_.Token == token.Semicolon {
		this_.Next()
		return
	}

	if this_.ImplicitSemicolon {
		this_.ImplicitSemicolon = false
		return
	}

	if this_.Token != token.Eof && this_.Token != token.RightBrace {
		this_.Expect("optionalSemicolon", token.Semicolon)
	}
}

func (this_ *Parser) Semicolon(from string) {
	if this_.Token != token.RightParenthesis && this_.Token != token.RightBrace {
		if this_.ImplicitSemicolon {
			this_.ImplicitSemicolon = false
			return
		}

		this_.Expect("semicolon from "+from, token.Semicolon)
	}
}

type State struct {
	Idx                                int
	Tok                                token.Token
	Literal                            string
	ParsedLiteral                      node.String
	ImplicitSemicolon, InsertSemicolon bool
	Chr                                rune
	ChrOffset, Offset                  int
	ErrorCount                         int
}

func (this_ *Parser) Mark(state *State) *State {
	if state == nil {
		state = &State{}
	}
	state.Idx, state.Tok, state.Literal, state.ParsedLiteral, state.ImplicitSemicolon, state.InsertSemicolon, state.Chr, state.ChrOffset, state.Offset =
		this_.Idx, this_.Token, this_.Literal, this_.ParsedLiteral, this_.ImplicitSemicolon, this_.InsertSemicolon, this_.Chr, this_.ChrOffset, this_.Offset

	state.ErrorCount = len(this_.Errors)
	return state
}

func (this_ *Parser) Restore(state *State) {
	this_.Idx, this_.Token, this_.Literal, this_.ParsedLiteral, this_.ImplicitSemicolon, this_.InsertSemicolon, this_.Chr, this_.ChrOffset, this_.Offset =
		state.Idx, state.Tok, state.Literal, state.ParsedLiteral, state.ImplicitSemicolon, state.InsertSemicolon, state.Chr, state.ChrOffset, state.Offset
	this_.Errors = this_.Errors[:state.ErrorCount]
}

func (this_ *Parser) Peek() token.Token {
	implicitSemicolon, insertSemicolon, chr, chrOffset, offset := this_.ImplicitSemicolon, this_.InsertSemicolon, this_.Chr, this_.ChrOffset, this_.Offset
	tok, _, _, _ := this_.Scan()
	this_.ImplicitSemicolon, this_.InsertSemicolon, this_.Chr, this_.ChrOffset, this_.Offset = implicitSemicolon, insertSemicolon, chr, chrOffset, offset
	return tok
}
