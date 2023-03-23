package parser

import "github.com/team-ide/go-interpreter/token"

func (this_ *parser) scanNumericLiteral(decimalPoint bool) (token.Token, string) {

	offset := this_.chrOffset
	tkn := token.Number

	if decimalPoint {
		offset--
		this_.scanMantissa(10)
	} else {
		if this_.chr == '0' {
			this_.read()
			base := 0
			switch this_.chr {
			case 'x', 'X':
				base = 16
			case 'o', 'O':
				base = 8
			case 'b', 'B':
				base = 2
			case '.', 'e', 'E':
				// no-op
			default:
				// legacy octal
				this_.scanMantissa(8)
				goto end
			}
			if base > 0 {
				this_.read()
				if !this_.IsDigit(this_.chr, base) {
					return token.Illegal, this_.str[offset:this_.chrOffset]
				}
				this_.scanMantissa(base)
				goto end
			}
		} else {
			this_.scanMantissa(10)
		}
		if this_.chr == '.' {
			this_.read()
			this_.scanMantissa(10)
		}
	}

	if this_.chr == 'e' || this_.chr == 'E' {
		this_.read()
		if this_.chr == '-' || this_.chr == '+' {
			this_.read()
		}
		if this_.IsDecimalDigit(this_.chr) {
			this_.read()
			this_.scanMantissa(10)
		} else {
			return token.Illegal, this_.str[offset:this_.chrOffset]
		}
	}
end:
	if this_.IsIdentifierStart(this_.chr) || this_.IsDecimalDigit(this_.chr) {
		return token.Illegal, this_.str[offset:this_.chrOffset]
	}

	return token.Token(tkn), this_.str[offset:this_.chrOffset]
}

func (this_ *parser) scanMantissa(base int) {
	for this_.DigitValue(this_.chr) < base {
		this_.read()
	}
}
