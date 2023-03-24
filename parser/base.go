package parser

import (
	"errors"
	"fmt"
	"github.com/dop251/goja/unistring"
	"github.com/team-ide/go-interpreter/token"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

func (this_ *parser) expect(value token.Token) int {
	idx := this_.idx
	if this_.token != value {
		_ = this_.errorUnexpectedToken(this_.token)
	}
	this_.next()
	return idx
}

func (this_ *parser) skipSingleLineComment() {
	for this_.chr != -1 {
		this_.read()
		if this_.IsLineTerminator(this_.chr) {
			return
		}
	}
}

func (this_ *parser) skipMultiLineComment() (hasLineTerminator bool) {
	this_.read()
	for this_.chr >= 0 {
		chr := this_.chr
		if chr == '\r' || chr == '\n' || chr == '\u2028' || chr == '\u2029' {
			hasLineTerminator = true
			break
		}
		this_.read()
		if chr == '*' && this_.chr == '/' {
			this_.read()
			return
		}
	}
	for this_.chr >= 0 {
		chr := this_.chr
		this_.read()
		if chr == '*' && this_.chr == '/' {
			this_.read()
			return
		}
	}

	_ = this_.errorUnexpected(0, this_.chr)
	return
}

func (this_ *parser) skipWhiteSpaceCheckLineTerminator() bool {
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
			return true
		}
		if this_.chr >= utf8.RuneSelf {
			if unicode.IsSpace(this_.chr) {
				this_.read()
				continue
			}
		}
		break
	}
	return false
}

func (this_ *parser) isBindingId(tok token.Token) bool {
	if tok == token.Identifier {
		return true
	}

	if tok == token.Await {
		return !this_.scope.allowAwait
	}
	if tok == token.Yield {
		return !this_.scope.allowYield
	}

	if this_.IsUnreservedWordToken(tok) {
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

func (this_ *parser) parseNumberLiteral(literal string) (value interface{}, err error) {
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

func (this_ *parser) parseStringLiteral(literal string, length int, unicode, strict bool) (unistring.String, string) {
	var sb strings.Builder
	var chars []uint16
	if unicode {
		chars = make([]uint16, 1, length+1)
		chars[0] = unistring.BOM
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
		return unistring.FromUtf16(chars), ""
	}
	if sb.Len() != length {
		panic(fmt.Errorf("unexpected length while parsing '%s'", literal))
	}
	return unistring.String(sb.String()), ""
}
