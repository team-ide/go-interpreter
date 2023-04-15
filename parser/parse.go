package parser

import (
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/token"
)

// Parse 解析
func (this_ *Parser) Parse() (tree *node.Tree, err error) {
	tree = this_.parseTree()
	//this_.errors.Sort()
	err = this_.Errors.Err()
	return
}

func (this_ *Parser) parseTree() (tree *node.Tree) {
	this_.OpenScope()
	defer this_.CloseScope()

	var statements []node.Statement
	this_.Read()
	this_.Next()
	for this_.Token != token.Eof {
		this_.Scope.AllowLet = true
		statements = append(statements, this_.ParseStatement())
	}

	tree = &node.Tree{
		Children:        statements,
		DeclarationList: this_.Scope.DeclarationList,
		OffsetPosition:  this_.OffsetPosition,
	}
	//this_.file.SetSourceMap(this_.parseSourceMap())
	return
}

// ParseBlockStatement 解析 {} 子 语句
func (this_ *Parser) ParseBlockStatement() *node.BlockStatement {
	res := &node.BlockStatement{}
	res.LeftBrace = this_.ExpectAndNext("ParseBlockStatement", token.LeftBrace)
	res.List = this_.ParseStatementList()
	res.RightBrace = this_.ExpectAndNext("ParseBlockStatement", token.RightBrace)

	return res
}

// ParseSemicolonStatement 分号 ; 语句
func (this_ *Parser) ParseSemicolonStatement() node.Statement {
	idx := this_.ExpectAndNext("ParseSemicolonStatement", token.Semicolon)
	return &node.SemicolonStatement{Semicolon: idx}
}

// ParseBlankSpaceStatement 空白 语句
func (this_ *Parser) ParseBlankSpaceStatement() node.Statement {
	this_.ExpectAndNext("ParseBlankSpaceStatement", token.BlankSpace)
	return &node.BlankSpaceStatement{From: this_.BlankSpaceFrom, To: this_.BlankSpaceTo}
}

// ParseStatementList 解析 子语句
func (this_ *Parser) ParseStatementList() (list []node.Statement) {
	for this_.Token != token.RightBrace && this_.Token != token.Eof {
		this_.Scope.AllowLet = true
		list = append(list, this_.ParseStatement())
	}

	return
}

// ParseIdentifier 解析 标识符
func (this_ *Parser) ParseIdentifier() *node.Identifier {
	literal := this_.ParsedLiteral
	idx := this_.Idx
	this_.Next()
	return &node.Identifier{
		Name: literal,
		Idx:  idx,
	}
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
		this_.ExpectAndNext("optionalSemicolon", token.Semicolon)
	}
}

func (this_ *Parser) Semicolon(from string) {
	if this_.Token != token.RightParenthesis && this_.Token != token.RightBrace {
		if this_.ImplicitSemicolon {
			this_.ImplicitSemicolon = false
			return
		}

		this_.ExpectAndNext("semicolon from "+from, token.Semicolon)
	}
}

type State struct {
	Idx                                int
	Tok                                token.Token
	Literal                            string
	ParsedLiteral                      string
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
