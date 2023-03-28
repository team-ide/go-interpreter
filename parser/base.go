package parser

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/token"
	"golang.org/x/text/unicode/rangetable"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

// IsDecimalDigit 是十进制数字
func (this_ *Parser) IsDecimalDigit(chr rune) bool {
	return '0' <= chr && chr <= '9'
}

// DigitValue 数字值
func (this_ *Parser) DigitValue(chr rune) int {
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
func (this_ *Parser) IsDigit(chr rune, base int) bool {
	return this_.DigitValue(chr) < base
}

// IsIdentifier 是标识符
func (this_ *Parser) IsIdentifier(s string) bool {
	if s == "" {
		return false
	}
	r, size := utf8.DecodeRuneInString(s)
	if !this_.IsIdentifierStart(r) {
		return false
	}
	for _, r := range s[size:] {
		if !this_.IsIdentifierPart(r) {
			return false
		}
	}
	return true
}

// IsIdentifierStart 是标识符开始
func (this_ *Parser) IsIdentifierStart(chr rune) bool {
	return chr == '$' || chr == '_' || chr == '\\' ||
		'a' <= chr && chr <= 'z' || 'A' <= chr && chr <= 'Z' ||
		chr >= utf8.RuneSelf && isIdentifierStartUnicode(chr)
}

// IsIdentifierPart 是标识符组成部分
func (this_ *Parser) IsIdentifierPart(chr rune) bool {
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
func (this_ *Parser) IsLineWhiteSpace(chr rune) bool {
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
func (this_ *Parser) IsLineTerminator(chr rune) bool {
	switch chr {
	// 换行符 /n、回车 \r、行分隔符、段落分隔符
	case '\u000a', '\u000d', '\u2028', '\u2029':
		return true
	}
	return false
}

func (this_ *Parser) Expect(from string, value token.Token) int {
	idx := this_.Idx
	if this_.Token != value {
		_ = this_.ErrorUnexpectedToken("expect by "+from+" this_.token:"+this_.Token.String()+",value:"+value.String(), this_.Token)
	}
	this_.Next()
	return idx
}

func (this_ *Parser) IsBindingId(tok token.Token) bool {
	if tok == token.Identifier {
		return true
	}

	if tok == token.Await {
		return !this_.Scope.AllowAwait
	}
	if tok == token.Yield {
		return !this_.Scope.AllowYield
	}

	if this_.IsUnreservedWordToken(tok) {
		return true
	}
	return false
}

func (this_ *Parser) IsLogicalAndExpr(expr node.Expression) bool {
	if exp, ok := expr.(*node.BinaryExpression); ok && exp.Operator == token.LogicalAnd {
		return true
	}
	return false
}

func normaliseCRLF(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\r' {
			buf.WriteByte('\n')
			if i < len(s)-1 && s[i+1] == '\n' {
				i++
			}
		} else {
			buf.WriteByte(s[i])
		}
	}
	return buf.String()
}

func hex2decimal(chr byte) (value rune, ok bool) {
	{
		chr := rune(chr)
		switch {
		case '0' <= chr && chr <= '9':
			return chr - '0', true
		case 'a' <= chr && chr <= 'f':
			return chr - 'a' + 10, true
		case 'A' <= chr && chr <= 'F':
			return chr - 'A' + 10, true
		}
		return
	}
}

func (this_ *Parser) ParseNumberLiteral(literal string) (value interface{}, err error) {
	// TODO Is Uint okay? What about -MAX_UINT
	value, err = strconv.ParseInt(literal, 0, 64)
	if err == nil {
		return
	}

	parseIntErr := err // Save this first error, just in case

	value, err = strconv.ParseFloat(literal, 64)
	if err == nil {
		return
	} else if err.(*strconv.NumError).Err == strconv.ErrRange {
		// Infinity, etc.
		return value, nil
	}

	err = parseIntErr

	if err.(*strconv.NumError).Err == strconv.ErrRange {
		if len(literal) > 2 && literal[0] == '0' && (literal[1] == 'X' || literal[1] == 'x') {
			// Could just be a very large number (e.g. 0x8000000000000000)
			var value float64
			literal = literal[2:]
			for _, chr := range literal {
				digit := this_.DigitValue(chr)
				if digit >= 16 {
					goto error
				}
				value = value*16 + float64(digit)
			}
			return value, nil
		}
	}

error:
	return nil, errors.New("illegal numeric literal")
}

func (this_ *Parser) parseStringLiteral(literal string, length int, unicode, strict bool) (node.String, string) {
	var sb strings.Builder
	var chars []uint16
	if unicode {
		chars = make([]uint16, 1, length+1)
		chars[0] = node.BOM
	} else {
		sb.Grow(length)
	}
	str := literal
	for len(str) > 0 {
		switch chr := str[0]; {
		// We do not explicitly handle the case of the quote
		// value, which can be: " ' /
		// This assumes we're already passed a partially well-formed literal
		case chr >= utf8.RuneSelf:
			chr, size := utf8.DecodeRuneInString(str)
			if chr <= 0xFFFF {
				chars = append(chars, uint16(chr))
			} else {
				first, second := utf16.EncodeRune(chr)
				chars = append(chars, uint16(first), uint16(second))
			}
			str = str[size:]
			continue
		case chr != '\\':
			if unicode {
				chars = append(chars, uint16(chr))
			} else {
				sb.WriteByte(chr)
			}
			str = str[1:]
			continue
		}

		if len(str) <= 1 {
			panic("len(str) <= 1")
		}
		chr := str[1]
		var value rune
		if chr >= utf8.RuneSelf {
			str = str[1:]
			var size int
			value, size = utf8.DecodeRuneInString(str)
			str = str[size:] // \ + <character>
			if value == '\u2028' || value == '\u2029' {
				continue
			}
		} else {
			str = str[2:] // \<character>
			switch chr {
			case 'b':
				value = '\b'
			case 'f':
				value = '\f'
			case 'n':
				value = '\n'
			case 'r':
				value = '\r'
			case 't':
				value = '\t'
			case 'v':
				value = '\v'
			case 'x', 'u':
				size := 0
				switch chr {
				case 'x':
					size = 2
				case 'u':
					if str == "" || str[0] != '{' {
						size = 4
					}
				}
				if size > 0 {
					if len(str) < size {
						return "", fmt.Sprintf("invalid escape: \\%s: len(%q) != %d", string(chr), str, size)
					}
					for j := 0; j < size; j++ {
						decimal, ok := hex2decimal(str[j])
						if !ok {
							return "", fmt.Sprintf("invalid escape: \\%s: %q", string(chr), str[:size])
						}
						value = value<<4 | decimal
					}
				} else {
					str = str[1:]
					var val rune
					value = -1
					for ; size < len(str); size++ {
						if str[size] == '}' {
							if size == 0 {
								return "", fmt.Sprintf("invalid escape: \\%s", string(chr))
							}
							size++
							value = val
							break
						}
						decimal, ok := hex2decimal(str[size])
						if !ok {
							return "", fmt.Sprintf("invalid escape: \\%s: %q", string(chr), str[:size+1])
						}
						val = val<<4 | decimal
						if val > utf8.MaxRune {
							return "", fmt.Sprintf("undefined Unicode code-point: %q", str[:size+1])
						}
					}
					if value == -1 {
						return "", fmt.Sprintf("unterminated \\u{: %q", str)
					}
				}
				str = str[size:]
				if chr == 'x' {
					break
				}
				if value > utf8.MaxRune {
					panic("value > utf8.MaxRune")
				}
			case '0':
				if len(str) == 0 || '0' > str[0] || str[0] > '7' {
					value = 0
					break
				}
				fallthrough
			case '1', '2', '3', '4', '5', '6', '7':
				if strict {
					return "", "Octal escape sequences are not allowed in this context"
				}
				value = rune(chr) - '0'
				j := 0
				for ; j < 2; j++ {
					if len(str) < j+1 {
						break
					}
					chr := str[j]
					if '0' > chr || chr > '7' {
						break
					}
					decimal := rune(str[j]) - '0'
					value = (value << 3) | decimal
				}
				str = str[j:]
			case '\\':
				value = '\\'
			case '\'', '"':
				value = rune(chr)
			case '\r':
				if len(str) > 0 {
					if str[0] == '\n' {
						str = str[1:]
					}
				}
				fallthrough
			case '\n':
				continue
			default:
				value = rune(chr)
			}
		}
		if unicode {
			if value <= 0xFFFF {
				chars = append(chars, uint16(value))
			} else {
				first, second := utf16.EncodeRune(value)
				chars = append(chars, uint16(first), uint16(second))
			}
		} else {
			if value >= utf8.RuneSelf {
				return "", "Unexpected unicode character"
			}
			sb.WriteByte(byte(value))
		}
	}

	if unicode {
		if len(chars) != length+1 {
			panic(fmt.Errorf("unexpected unicode length while parsing '%s'", literal))
		}
		return node.FromUtf16(chars), ""
	}
	if sb.Len() != length {
		panic(fmt.Errorf("unexpected length while parsing '%s'", literal))
	}
	return node.String(sb.String()), ""
}
