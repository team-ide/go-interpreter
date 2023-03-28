package parser

import (
	"unicode"
	"unicode/utf8"
)

// SkipWhiteSpace 跳过 空白
func (this_ *Parser) SkipWhiteSpace() {
	for {
		switch this_.Chr {
		case ' ', '\t', '\f', '\v', '\u00a0', '\ufeff':
			this_.Read()
			continue
		case '\r':
			if this_.ImplicitRead() == '\n' {
				this_.Read()
			}
			fallthrough
		case '\u2028', '\u2029', '\n':
			if this_.InsertSemicolon {
				return
			}
			this_.Read()
			continue
		}
		if this_.Chr >= utf8.RuneSelf {
			if unicode.IsSpace(this_.Chr) {
				this_.Read()
				continue
			}
		}
		break
	}
}

// SkipSingleLineComment 跳过单行注释
func (this_ *Parser) SkipSingleLineComment() {
	for this_.Chr != -1 {
		this_.Read()
		if this_.IsLineTerminator(this_.Chr) {
			return
		}
	}
}

// SkipMultiLineComment 跳过多行注释
func (this_ *Parser) SkipMultiLineComment() (hasLineTerminator bool) {
	this_.Read()
	for this_.Chr >= 0 {
		chr := this_.Chr
		// 换行符 /n、回车 \r、行分隔符、段落分隔符
		if chr == '\r' || chr == '\n' || chr == '\u2028' || chr == '\u2029' {
			hasLineTerminator = true
			break
		}
		this_.Read()
		if chr == '*' && this_.Chr == '/' {
			this_.Read()
			return
		}
	}
	for this_.Chr >= 0 {
		chr := this_.Chr
		this_.Read()
		if chr == '*' && this_.Chr == '/' {
			this_.Read()
			return
		}
	}

	_ = this_.ErrorUnexpected("skipMultiLineComment", 0, this_.Chr)
	return
}

// SkipWhiteSpaceCheckLineTerminator 跳过空白检查线路终止器
func (this_ *Parser) SkipWhiteSpaceCheckLineTerminator() bool {
	for {
		switch this_.Chr {
		case ' ', '\t', '\f', '\v', '\u00a0', '\ufeff':
			this_.Read()
			continue
		case '\r':
			if this_.ImplicitRead() == '\n' {
				this_.Read()
			}
			fallthrough
		case '\u2028', '\u2029', '\n':
			return true
		}
		if this_.Chr >= utf8.RuneSelf {
			if unicode.IsSpace(this_.Chr) {
				this_.Read()
				continue
			}
		}
		break
	}
	return false
}
