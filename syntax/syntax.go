package syntax

import (
	"github.com/team-ide/go-interpreter/token"
	"golang.org/x/text/unicode/rangetable"
	"unicode"
	"unicode/utf8"
)

// Syntax 语法 接口 解释器解析各语言需要借助该接口
type Syntax interface {
	// IsDecimalDigit 是十进制数字
	IsDecimalDigit(chr rune) bool
	// DigitValue 数字值
	DigitValue(chr rune) int
	// IsDigit 是数字
	IsDigit(chr rune, base int) bool
	// IsIdentifier 是标识符
	IsIdentifier(s string) bool
	// IsIdentifierStart 是标识符开始
	IsIdentifierStart(chr rune) bool
	// IsIdentifierPart 是标识符组成部分
	IsIdentifierPart(chr rune) bool
	// IsLineWhiteSpace 是空白行
	IsLineWhiteSpace(chr rune) bool
	// IsLineTerminator 是行末尾
	IsLineTerminator(chr rune) bool
	// IsKeyword 是关键字
	IsKeyword(literal string) (token.Token, bool)
}

// IsDecimalDigit 是十进制数字
func IsDecimalDigit(chr rune) bool {
	return '0' <= chr && chr <= '9'
}

// DigitValue 数字值
func DigitValue(chr rune) int {
	switch {
	case '0' <= chr && chr <= '9':
		return int(chr - '0')
	case 'a' <= chr && chr <= 'f':
		return int(chr - 'a' + 10)
	case 'A' <= chr && chr <= 'F':
		return int(chr - 'A' + 10)
	}
	return 16 // Larger than any legal digit value
}

// IsDigit 是数字
func IsDigit(chr rune, base int) bool {
	return DigitValue(chr) < base
}

// IsIdentifier 是标识符
func IsIdentifier(s string) bool {
	if s == "" {
		return false
	}
	r, size := utf8.DecodeRuneInString(s)
	if !IsIdentifierStart(r) {
		return false
	}
	for _, r := range s[size:] {
		if !IsIdentifierPart(r) {
			return false
		}
	}
	return true
}

// IsIdentifierStart 是标识符开始
func IsIdentifierStart(chr rune) bool {
	return chr == '$' || chr == '_' || chr == '\\' ||
		'a' <= chr && chr <= 'z' || 'A' <= chr && chr <= 'Z' ||
		chr >= utf8.RuneSelf && isIdentifierStartUnicode(chr)
}

// IsIdentifierPart 是标识符组成部分
func IsIdentifierPart(chr rune) bool {
	return chr == '$' || chr == '_' || chr == '\\' ||
		'a' <= chr && chr <= 'z' || 'A' <= chr && chr <= 'Z' ||
		'0' <= chr && chr <= '9' ||
		chr >= utf8.RuneSelf && isIdentifierPartUnicode(chr)
}

var (
	unicodeRangeIdentifierNeg      = rangetable.Merge(unicode.Pattern_Syntax, unicode.Pattern_White_Space)
	unicodeRangeIdentifierStartPos = rangetable.Merge(unicode.Letter, unicode.Nl, unicode.Other_ID_Start)
	unicodeRangeIdentifierContPos  = rangetable.Merge(unicodeRangeIdentifierStartPos, unicode.Mn, unicode.Mc, unicode.Nd, unicode.Pc, unicode.Other_ID_Continue)
)

func isIdentifierStartUnicode(r rune) bool {
	return unicode.Is(unicodeRangeIdentifierStartPos, r) && !unicode.Is(unicodeRangeIdentifierNeg, r)
}

func isIdentifierPartUnicode(r rune) bool {
	return unicode.Is(unicodeRangeIdentifierContPos, r) && !unicode.Is(unicodeRangeIdentifierNeg, r) || r == '\u200C' || r == '\u200D'
}

/*
\u200C 是一个 Unicode 字符，它代表零宽度非连接符 它是一个不可见的字符，用于控制字符之间的连接
\u200D 是一个 Unicode 字符，它代表零宽度连接符 它是一个不可见的字符，用于控制字符之间的连接

\u0009 Tab \t
\u000b 垂直制表符 \v
\u000c 换页 \f
\u0020 unicode 半角空格
\u3000 全角空格
\u00a0 不间断空格
\ufeff 字节顺序标记

\u000a 换行符 \n
\u000d 回车 \r
\u2028 行分隔符
\u2029 段落分隔符

\u0085 代表下一行的字符

\u2028 行分隔符
\u2029 段落分隔符
*/

// IsLineWhiteSpace 是空白行
func IsLineWhiteSpace(chr rune) bool {
	switch chr {
	// Tab \t、垂直制表符 \v、换页 \f、unicode 半角空格、全角空格、不间断空格、字节顺序标记
	case '\u0009', '\u000b', '\u000c', '\u0020', '\u3000', '\u00a0', '\ufeff':
		return true
	// 换行符 \n、回车 \r、行分隔符、段落分隔符
	case '\u000a', '\u000d', '\u2028', '\u2029':
		return false
	// 代表下一行的字符
	case '\u0085':
		return false
	}
	return unicode.IsSpace(chr)
}

// IsLineTerminator 是行末尾
func IsLineTerminator(chr rune) bool {
	switch chr {
	// 换行符 /n、回车 \r、行分隔符、段落分隔符
	case '\u000a', '\u000d', '\u2028', '\u2029':
		return true
	}
	return false
}
