package parser

import (
	"fmt"
	"github.com/team-ide/go-interpreter/node"
	"unicode/utf8"
)

func (this_ *parser) scanIdentifier() (string, node.String, bool, string) {
	offset := this_.chrOffset
	hasEscape := false
	isUnicode := false
	length := 0
	for this_.IsIdentifierPart(this_.chr) {
		r := this_.chr
		length++
		if r == '\\' {
			hasEscape = true
			distance := this_.chrOffset - offset
			this_.read()
			if this_.chr != 'u' {
				return "", "", false, fmt.Sprintf("Invalid identifier escape character: %c (%s)", this_.chr, string(this_.chr))
			}
			var value rune
			if this_.implicitRead() == '{' {
				this_.read()
				value = -1
				for value <= utf8.MaxRune {
					this_.read()
					if this_.chr == '}' {
						break
					}
					decimal, ok := hex2decimal(byte(this_.chr))
					if !ok {
						return "", "", false, "Invalid Unicode escape sequence"
					}
					if value == -1 {
						value = decimal
					} else {
						value = value<<4 | decimal
					}
				}
				if value == -1 {
					return "", "", false, "Invalid Unicode escape sequence"
				}
			} else {
				for j := 0; j < 4; j++ {
					this_.read()
					decimal, ok := hex2decimal(byte(this_.chr))
					if !ok {
						return "", "", false, fmt.Sprintf("Invalid identifier escape character: %c (%s)", this_.chr, string(this_.chr))
					}
					value = value<<4 | decimal
				}
			}
			if value == '\\' {
				return "", "", false, fmt.Sprintf("Invalid identifier escape value: %c (%s)", value, string(value))
			} else if distance == 0 {
				if !this_.IsIdentifierStart(value) {
					return "", "", false, fmt.Sprintf("Invalid identifier escape value: %c (%s)", value, string(value))
				}
			} else if distance > 0 {
				if !this_.IsIdentifierPart(value) {
					return "", "", false, fmt.Sprintf("Invalid identifier escape value: %c (%s)", value, string(value))
				}
			}
			r = value
		}
		if r >= utf8.RuneSelf {
			isUnicode = true
			if r > 0xFFFF {
				length++
			}
		}
		this_.read()
	}

	literal := this_.str[offset:this_.chrOffset]
	var parsed node.String
	if hasEscape || isUnicode {
		var err string
		// TODO strict
		parsed, err = this_.parseStringLiteral(literal, length, isUnicode, false)
		if err != "" {
			return "", "", false, err
		}
	} else {
		parsed = node.String(literal)
	}

	return literal, parsed, hasEscape, ""
}
