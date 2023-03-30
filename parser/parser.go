package parser

import (
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/token"
	"unicode/utf8"
)

func New(src string) (p *Parser) {
	p = &Parser{
		//Chr:    ' ',
		Str:    src,
		Length: len(src),
		Base:   0,
	}
	return
}

// TokenIndexOf 返回 某个值 在数组中的索引位置，未找到返回 -1
func TokenIndexOf(array []token.Token, v token.Token) (index int) {
	index = -1
	size := len(array)
	for i := 0; i < size; i++ {
		if array[i] == v {
			index = i
			return
		}
	}
	return
}

type Keyword struct {
	Token token.Token
	// 未来关键字
	FutureKeyword bool
	// 严格的
	Strict bool
}

func (this_ *Parser) IsKeyword(literal string) (token.Token, bool) {
	KeywordToken := this_.KeywordToken
	if KeywordToken != nil {
		if keyword, exists := KeywordToken[literal]; exists {
			if keyword.FutureKeyword {
				return token.Keyword, keyword.Strict
			}
			return keyword.Token, false
		}
	}
	return "", false
}
func (this_ *Parser) IsIdentifierToken(tkn token.Token) bool {
	return TokenIndexOf(this_.IdentifierTokens, tkn) >= 0
}
func (this_ *Parser) IsUnreservedWordToken(tkn token.Token) bool {
	return TokenIndexOf(this_.UnreservedWordTokens, tkn) >= 0
}

type Parser struct {
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
	InsertSemicolon   bool // 如果我们看到一个换行符，那么插入一个隐式分号
	ImplicitSemicolon bool // 存在隐式分号

	Recover struct {
		// Scratch when trying to seek to the next statement, etc.
		Idx   int
		Count int
	}

	Errors ErrorList

	ParseStatement func() node.Statement

	KeywordToken         map[string]Keyword
	IdentifierTokens     []token.Token
	UnreservedWordTokens []token.Token

	BlankSpaceFrom int
	BlankSpaceTo   int
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
