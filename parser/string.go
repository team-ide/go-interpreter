package parser

import (
	"github.com/team-ide/go-interpreter/node"
	"unicode/utf8"
)

func (this_ *Parser) ParseTemplateCharacters() (literal string, parsed node.String, finished bool, parseErr, err string) {
	offset := this_.ChrOffset
	var end int
	length := 0
	isUnicode := false
	hasCR := false
	for {
		chr := this_.Chr
		if chr < 0 {
			goto unterminated
		}
		this_.Read()
		if chr == '`' {
			finished = true
			end = this_.ChrOffset - 1
			break
		}
		if chr == '\\' {
			if this_.Chr == '\n' || this_.Chr == '\r' || this_.Chr == '\u2028' || this_.Chr == '\u2029' || this_.Chr < 0 {
				if this_.Chr == '\r' {
					hasCR = true
				}
				this_.ScanNewline()
			} else {
				if this_.Chr == '8' || this_.Chr == '9' {
					if parseErr == "" {
						parseErr = "\\8 and \\9 are not allowed in template strings."
					}
				}
				l, u := this_.ScanEscape('`')
				length += l
				if u {
					isUnicode = true
				}
			}
			continue
		}
		if chr == '$' && this_.Chr == '{' {
			this_.Read()
			end = this_.ChrOffset - 2
			break
		}
		if chr >= utf8.RuneSelf {
			isUnicode = true
			if chr > 0xFFFF {
				length++
			}
		} else if chr == '\r' {
			hasCR = true
			if this_.Chr == '\n' {
				length--
			}
		}
		length++
	}
	literal = this_.Str[offset:end]
	if hasCR {
		literal = normaliseCRLF(literal)
	}
	if parseErr == "" {
		parsed, parseErr = this_.parseStringLiteral(literal, length, isUnicode, true)
	}
	this_.InsertSemicolon = true
	return
unterminated:
	err = errUnexpectedEndOfInput
	finished = true
	return
}

func (this_ *Parser) ScanString(offset int, parse bool) (literal string, parsed node.String, err string) {
	// " ' /
	quote := rune(this_.Str[offset])
	length := 0
	isUnicode := false
	for this_.Chr != quote {
		chr := this_.Chr
		if chr == '\n' || chr == '\r' || chr < 0 {
			goto newline
		}
		if quote == '/' && (this_.Chr == '\u2028' || this_.Chr == '\u2029') {
			goto newline
		}
		this_.Read()
		if chr == '\\' {
			if this_.Chr == '\n' || this_.Chr == '\r' || this_.Chr == '\u2028' || this_.Chr == '\u2029' || this_.Chr < 0 {
				if quote == '/' {
					goto newline
				}
				this_.ScanNewline()
			} else {
				l, u := this_.ScanEscape(quote)
				length += l
				if u {
					isUnicode = true
				}
			}
			continue
		} else if chr == '[' && quote == '/' {
			// Allow a slash (/) in a bracket character class ([...])
			// TODO Fix this, this is hacky...
			quote = -1
		} else if chr == ']' && quote == -1 {
			quote = '/'
		}
		if chr >= utf8.RuneSelf {
			isUnicode = true
			if chr > 0xFFFF {
				length++
			}
		}
		length++
	}

	// " ' /
	this_.Read()
	literal = this_.Str[offset:this_.ChrOffset]
	if parse {
		// TODO strict
		parsed, err = this_.parseStringLiteral(literal[1:len(literal)-1], length, isUnicode, false)
	}
	return

newline:
	this_.ScanNewline()
	errStr := "String not terminated"
	if quote == '/' {
		errStr = "Invalid regular expression: missing /"
		_ = this_.Error("scanString", this_.IdxOf(offset), errStr)
	}
	return "", "", errStr
}

func (this_ *Parser) ScanEscape(quote rune) (int, bool) {

	var length, base uint32
	chr := this_.Chr
	switch chr {
	case '0', '1', '2', '3', '4', '5', '6', '7':
		//    Octal:
		length, base = 3, 8
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '"', '\'':
		this_.Read()
		return 1, false
	case '\r':
		this_.Read()
		if this_.Chr == '\n' {
			this_.Read()
			return 2, false
		}
		return 1, false
	case '\n':
		this_.Read()
		return 1, false
	case '\u2028', '\u2029':
		this_.Read()
		return 1, true
	case 'x':
		this_.Read()
		length, base = 2, 16
	case 'u':
		this_.Read()
		if this_.Chr == '{' {
			this_.Read()
			length, base = 0, 16
		} else {
			length, base = 4, 16
		}
	default:
		this_.Read() // Always make progress
	}

	if base > 0 {
		var value uint32
		if length > 0 {
			for ; length > 0 && this_.Chr != quote && this_.Chr >= 0; length-- {
				digit := uint32(this_.DigitValue(this_.Chr))
				if digit >= base {
					break
				}
				value = value*base + digit
				this_.Read()
			}
		} else {
			for this_.Chr != quote && this_.Chr >= 0 && value < utf8.MaxRune {
				if this_.Chr == '}' {
					this_.Read()
					break
				}
				digit := uint32(this_.DigitValue(this_.Chr))
				if digit >= base {
					break
				}
				value = value*base + digit
				this_.Read()
			}
		}
		chr = rune(value)
	}
	if chr >= utf8.RuneSelf {
		if chr > 0xFFFF {
			return 2, true
		}
		return 1, true
	}
	return 1, false
}

func (this_ *Parser) ScanNewline() {
	if this_.Chr == '\u2028' || this_.Chr == '\u2029' {
		this_.Read()
		return
	}
	if this_.Chr == '\r' {
		this_.Read()
		if this_.Chr != '\n' {
			return
		}
	}
	this_.Read()
}
