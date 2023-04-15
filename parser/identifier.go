package parser

import (
	"fmt"
	"unicode/utf8"
)

func (this_ *Parser) ScanIdentifier() (string, string, bool, string) {
	offset := this_.ChrOffset
	hasEscape := false
	isUnicode := false
	length := 0
	for this_.IsIdentifierPart(this_.Chr) {
		r := this_.Chr
		length++
		if r == '\\' {
			hasEscape = true
			distance := this_.ChrOffset - offset
			this_.Read()
			if this_.Chr != 'u' {
				return "", "", false, fmt.Sprintf("Invalid identifier escape character: %c (%s)", this_.Chr, string(this_.Chr))
			}
			var value rune
			if this_.ImplicitRead() == '{' {
				this_.Read()
				value = -1
				for value <= utf8.MaxRune {
					this_.Read()
					if this_.Chr == '}' {
						break
					}
					decimal, ok := hex2decimal(byte(this_.Chr))
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
					this_.Read()
					decimal, ok := hex2decimal(byte(this_.Chr))
					if !ok {
						return "", "", false, fmt.Sprintf("Invalid identifier escape character: %c (%s)", this_.Chr, string(this_.Chr))
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
		this_.Read()
	}

	literal := this_.Str[offset:this_.ChrOffset]
	var parsed string
	if hasEscape || isUnicode {
		var err string
		// TODO strict
		parsed, err = this_.parseStringLiteral(literal, length, isUnicode, false)
		if err != "" {
			return "", "", false, err
		}
	} else {
		parsed = literal
	}

	return literal, parsed, hasEscape, ""
}
