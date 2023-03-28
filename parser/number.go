package parser

import "github.com/team-ide/go-interpreter/token"

func (this_ *Parser) ScanNumericLiteral(decimalPoint bool) (token.Token, string) {

	offset := this_.ChrOffset
	tkn := token.Number

	if decimalPoint {
		offset--
		this_.ScanMantissa(10)
	} else {
		if this_.Chr == '0' {
			this_.Read()
			base := 0
			switch this_.Chr {
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
				this_.ScanMantissa(8)
				goto end
			}
			if base > 0 {
				this_.Read()
				if !this_.IsDigit(this_.Chr, base) {
					return token.Illegal, this_.Str[offset:this_.ChrOffset]
				}
				this_.ScanMantissa(base)
				goto end
			}
		} else {
			this_.ScanMantissa(10)
		}
		if this_.Chr == '.' {
			this_.Read()
			this_.ScanMantissa(10)
		}
	}

	if this_.Chr == 'e' || this_.Chr == 'E' {
		this_.Read()
		if this_.Chr == '-' || this_.Chr == '+' {
			this_.Read()
		}
		if this_.IsDecimalDigit(this_.Chr) {
			this_.Read()
			this_.ScanMantissa(10)
		} else {
			return token.Illegal, this_.Str[offset:this_.ChrOffset]
		}
	}
end:
	if this_.IsIdentifierStart(this_.Chr) || this_.IsDecimalDigit(this_.Chr) {
		return token.Illegal, this_.Str[offset:this_.ChrOffset]
	}

	return tkn, this_.Str[offset:this_.ChrOffset]
}
