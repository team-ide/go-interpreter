package parser

import (
	"github.com/team-ide/go-interpreter/node"
	"unicode/utf8"
)

func (this_ *parser) parseTemplateCharacters() (literal string, parsed node.String, finished bool, parseErr, err string) {
	offset := this_.chrOffset
	var end int
	length := 0
	isUnicode := false
	hasCR := false
	for {
		chr := this_.chr
		if chr < 0 {
			goto unterminated
		}
		this_.read()
		if chr == '`' {
			finished = true
			end = this_.chrOffset - 1
			break
		}
		if chr == '\\' {
			if this_.chr == '\n' || this_.chr == '\r' || this_.chr == '\u2028' || this_.chr == '\u2029' || this_.chr < 0 {
				if this_.chr == '\r' {
					hasCR = true
				}
				this_.scanNewline()
			} else {
				if this_.chr == '8' || this_.chr == '9' {
					if parseErr == "" {
						parseErr = "\\8 and \\9 are not allowed in template strings."
					}
				}
				l, u := this_.scanEscape('`')
				length += l
				if u {
					isUnicode = true
				}
			}
			continue
		}
		if chr == '$' && this_.chr == '{' {
			this_.read()
			end = this_.chrOffset - 2
			break
		}
		if chr >= utf8.RuneSelf {
			isUnicode = true
			if chr > 0xFFFF {
				length++
			}
		} else if chr == '\r' {
			hasCR = true
			if this_.chr == '\n' {
				length--
			}
		}
		length++
	}
	literal = this_.str[offset:end]
	if hasCR {
		literal = normaliseCRLF(literal)
	}
	if parseErr == "" {
		parsed, parseErr = this_.parseStringLiteral(literal, length, isUnicode, true)
	}
	this_.insertSemicolon = true
	return
unterminated:
	err = errUnexpectedEndOfInput
	finished = true
	return
}

func (this_ *parser) scanString(offset int, parse bool) (literal string, parsed node.String, err string) {
	// " ' /
	quote := rune(this_.str[offset])
	length := 0
	isUnicode := false
	for this_.chr != quote {
		chr := this_.chr
		if chr == '\n' || chr == '\r' || chr < 0 {
			goto newline
		}
		if quote == '/' && (this_.chr == '\u2028' || this_.chr == '\u2029') {
			goto newline
		}
		this_.read()
		if chr == '\\' {
			if this_.chr == '\n' || this_.chr == '\r' || this_.chr == '\u2028' || this_.chr == '\u2029' || this_.chr < 0 {
				if quote == '/' {
					goto newline
				}
				this_.scanNewline()
			} else {
				l, u := this_.scanEscape(quote)
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
	this_.read()
	literal = this_.str[offset:this_.chrOffset]
	if parse {
		// TODO strict
		parsed, err = this_.parseStringLiteral(literal[1:len(literal)-1], length, isUnicode, false)
	}
	return

newline:
	this_.scanNewline()
	errStr := "String not terminated"
	if quote == '/' {
		errStr = "Invalid regular expression: missing /"
		_ = this_.error("scanString", this_.idxOf(offset), errStr)
	}
	return "", "", errStr
}

func (this_ *parser) scanEscape(quote rune) (int, bool) {

	var length, base uint32
	chr := this_.chr
	switch chr {
	case '0', '1', '2', '3', '4', '5', '6', '7':
		//    Octal:
		length, base = 3, 8
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '"', '\'':
		this_.read()
		return 1, false
	case '\r':
		this_.read()
		if this_.chr == '\n' {
			this_.read()
			return 2, false
		}
		return 1, false
	case '\n':
		this_.read()
		return 1, false
	case '\u2028', '\u2029':
		this_.read()
		return 1, true
	case 'x':
		this_.read()
		length, base = 2, 16
	case 'u':
		this_.read()
		if this_.chr == '{' {
			this_.read()
			length, base = 0, 16
		} else {
			length, base = 4, 16
		}
	default:
		this_.read() // Always make progress
	}

	if base > 0 {
		var value uint32
		if length > 0 {
			for ; length > 0 && this_.chr != quote && this_.chr >= 0; length-- {
				digit := uint32(this_.DigitValue(this_.chr))
				if digit >= base {
					break
				}
				value = value*base + digit
				this_.read()
			}
		} else {
			for this_.chr != quote && this_.chr >= 0 && value < utf8.MaxRune {
				if this_.chr == '}' {
					this_.read()
					break
				}
				digit := uint32(this_.DigitValue(this_.chr))
				if digit >= base {
					break
				}
				value = value*base + digit
				this_.read()
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

func (this_ *parser) scanNewline() {
	if this_.chr == '\u2028' || this_.chr == '\u2029' {
		this_.read()
		return
	}
	if this_.chr == '\r' {
		this_.read()
		if this_.chr != '\n' {
			return
		}
	}
	this_.read()
}
