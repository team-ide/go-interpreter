package parser

import (
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/token"
	"unicode/utf8"
)

func New(src string) (p *Parser) {
	p = &Parser{
		//Chr:    ' ',
		Str:            src,
		Length:         len(src),
		OffsetPosition: map[int]*node.Position{},
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
func (this_ *Parser) IsModifierToken(tkn token.Token) bool {
	return TokenIndexOf(this_.ModifierTokens, tkn) >= 0
}

type Parser struct {
	Str    string
	Length int

	Chr       rune // 当前 字符
	ChrOffset int  // 当前 字符 偏移量
	Offset    int  // 当前 字符 偏移量

	Idx           int         // The index of token
	Token         token.Token // The token
	Literal       string      // The literal of the token, if any
	ParsedLiteral string

	InsertSemicolon   bool // 如果我们看到一个换行符，那么插入一个隐式分号
	ImplicitSemicolon bool // 存在隐式分号

	Errors ErrorList

	ParseStatement func() node.Statement

	KeywordToken         map[string]Keyword
	IdentifierTokens     []token.Token
	UnreservedWordTokens []token.Token
	ModifierTokens       []token.Token // 修饰符

	BlankSpaceFrom int
	BlankSpaceTo   int

	OnlyReadLess    bool // 只读取 < 不读取 <<、<<<
	OnlyReadGreater bool // 只读取 > 不读取 >>、>>>

	Modifiers []*Modifier // 修饰符 临时存放

	OffsetPosition map[int]*node.Position
	Line           int
	Column         int
}

type Modifier struct {
	Idx  int
	Name token.Token
}

func (this_ *Parser) AddModifier(modifier token.Token) {
	this_.Modifiers = append(this_.Modifiers, &Modifier{
		Idx:  this_.Idx,
		Name: modifier,
	})
}

func (this_ *Parser) GetAndClearModifiers() (modifiers []*Modifier) {
	modifiers = this_.Modifiers
	this_.Modifiers = []*Modifier{}
	return
}

// ImplicitRead 隐式读取下一个
func (this_ *Parser) ImplicitRead() rune {
	if this_.Offset < this_.Length {
		return rune(this_.Str[this_.Offset])
	}
	return -1
}

func (this_ *Parser) GetPosition(offset int) *node.Position {
	return this_.OffsetPosition[offset]
}

// Read 读取下一个 将重新设定偏移量
func (this_ *Parser) Read() {
	if this_.Offset < this_.Length {

		// 当前索引
		position := &node.Position{
			Line:   this_.Line + 1,
			Column: this_.Column + 1,
			Offset: this_.Offset,
		}
		this_.OffsetPosition[this_.Offset] = position

		this_.ChrOffset = this_.Offset
		chr, width := rune(this_.Str[this_.Offset]), 1
		// 检查 编码 是否 是 ASCII
		if chr >= utf8.RuneSelf { // !ASCII
			chr, width = utf8.DecodeRuneInString(this_.Str[this_.Offset:])
			if chr == utf8.RuneError && width == 1 {
				_ = this_.Error("read char utf8.RuneError chr:"+string(chr), this_.ChrOffset, "Invalid UTF-8 character")
			}
		}
		// 如果是换行
		if chr == '\n' {
			this_.Line++
			this_.Column = 0
		} else {
			this_.Column += width
		}
		this_.Offset += width
		this_.Chr = chr
	} else {
		this_.ChrOffset = this_.Length
		this_.Chr = -1 // EOF 读取结束
	}
}
