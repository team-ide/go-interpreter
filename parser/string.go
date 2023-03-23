package parser

import (
	"github.com/dop251/goja/unistring"
	"unicode/utf8"
)

func (this_ *parser) scanString(offset int, parse bool) (literal string, parsed unistring.String, err string) {
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
				l, u := self.scanEscape(quote)
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
		parsed, err = parseStringLiteral(literal[1:len(literal)-1], length, isUnicode, false)
	}
	return

newline:
	this_.scanNewline()
	errStr := "String not terminated"
	if quote == '/' {
		errStr = "Invalid regular expression: missing /"
		this_.error(this_.idxOf(offset), errStr)
	}
	return "", "", errStr
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
