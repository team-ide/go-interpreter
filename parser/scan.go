package parser

import (
	"github.com/team-ide/go-interpreter/token"
)

func (this_ *Parser) Scan() (tkn token.Token, literal string, parsedLiteral string, idx int) {

	this_.ImplicitSemicolon = false

	for {
		this_.SkipWhiteSpace()

		idx = this_.IdxOf(this_.ChrOffset)
		InsertSemicolon := false

		switch chr := this_.Chr; {
		case this_.IsIdentifierStart(chr):
			var err string
			var hasEscape bool
			literal, parsedLiteral, hasEscape, err = this_.ScanIdentifier()
			if err != "" {
				tkn = token.Illegal
				break
			}
			if len(parsedLiteral) > 1 {
				// Keywords are longer than 1 character, avoid lookup otherwise
				var strict bool
				tkn, strict = this_.IsKeyword(string(parsedLiteral))
				if hasEscape {
					this_.InsertSemicolon = true
					if tkn == "" || this_.IsBindingIdentifier(tkn) {
						tkn = token.Identifier
					} else {
						tkn = token.EscapedReservedWord
					}
					return
				}
				if this_.IsModifierToken(tkn) {
					this_.AddModifier(tkn)
					continue
				}
				switch tkn {
				case "": // Not a keyword
					// no-op
				case token.Keyword:
					if strict {
						// TODO If strict and in strict mode, then this is not a break
						break
					}
					return

				case
					token.Boolean,
					token.Nil,
					token.Null,
					token.This,
					token.Break,
					token.Throw, // A newline after a throw is not allowed, but we need to detect it
					token.Return,
					token.Continue,
					token.Debugger:
					this_.InsertSemicolon = true
					return

				case token.Async:
					// async only has special meaning if not followed by a LineTerminator
					if this_.SkipWhiteSpaceCheckLineTerminator() {
						this_.InsertSemicolon = true
						tkn = token.Identifier
					}
					return
				default:
					return

				}
			}
			this_.InsertSemicolon = true
			tkn = token.Identifier
			return
		case '0' <= chr && chr <= '9':
			this_.InsertSemicolon = true
			tkn, literal = this_.ScanNumericLiteral(false)
			return
		//case chr == '\n' || chr == '\r' || chr == '\t' || chr == ' ':
		//	tkn = token.BlankSpace
		//	return
		default:
			this_.Read()
			switch chr {
			case -1:
				if this_.InsertSemicolon {
					this_.InsertSemicolon = false
					this_.ImplicitSemicolon = true
				}
				tkn = token.Eof
			case '\r', '\n', '\u2028', '\u2029':
				this_.InsertSemicolon = false
				this_.ImplicitSemicolon = true
				continue
			case ':':
				tkn = token.Colon
			case '.':
				if this_.DigitValue(this_.Chr) < 10 {
					InsertSemicolon = true
					tkn, literal = this_.ScanNumericLiteral(true)
				} else {
					if this_.Chr == '.' {
						this_.Read()
						if this_.Chr == '.' {
							this_.Read()
							tkn = token.Ellipsis
						} else {
							tkn = token.Illegal
						}
					} else {
						tkn = token.Period
					}
				}
			case ',':
				tkn = token.Comma
			case ';':
				tkn = token.Semicolon
			case '(':
				tkn = token.LeftParenthesis
			case ')':
				tkn = token.RightParenthesis
				InsertSemicolon = true
			case '[':
				tkn = token.LeftBracket
			case ']':
				tkn = token.RightBracket
				InsertSemicolon = true
			case '{':
				tkn = token.LeftBrace
			case '}':
				tkn = token.RightBrace
				InsertSemicolon = true
			case '+':
				tkn = this_.Switch3(token.Plus, token.AddAssign, '+', token.Increment)
				if tkn == token.Increment {
					InsertSemicolon = true
				}
			case '-':
				tkn = this_.Switch3(token.Minus, token.SubtractAssign, '-', token.Decrement)
				if tkn == token.Decrement {
					InsertSemicolon = true
				}
			case '*':
				if this_.Chr == '*' {
					this_.Read()
					tkn = this_.Switch2(token.Exponent, token.ExponentAssign)
				} else {
					tkn = this_.Switch2(token.Multiply, token.MultiplyAssign)
				}
			case '/':
				if this_.Chr == '/' {
					this_.SkipSingleLineComment()
					continue
				} else if this_.Chr == '*' {
					if this_.SkipMultiLineComment() {
						this_.InsertSemicolon = false
						this_.ImplicitSemicolon = true
					}
					continue
				} else {
					// Could be division, could be RegExp literal
					tkn = this_.Switch2(token.Slash, token.QuotientAssign)
					InsertSemicolon = true
				}
			case '%':
				tkn = this_.Switch2(token.Remainder, token.RemainderAssign)
			case '^':
				tkn = this_.Switch2(token.ExclusiveOr, token.ExclusiveOrAssign)
			case '<':
				if this_.OnlyReadLess {
					tkn = token.Less
				} else {
					tkn = this_.Switch4(token.Less, token.LessOrEqual, '<', token.ShiftLeft, token.ShiftLeftAssign)
				}
			case '>':
				if this_.OnlyReadGreater {
					tkn = token.Greater
				} else {
					tkn = this_.Switch6(token.Greater, token.GreaterOrEqual, '>', token.ShiftRight, token.ShiftRightAssign, '>', token.UnsignedShiftRight, token.UnsignedShiftRightAssign)
				}
			case '=':
				if this_.Chr == '>' {
					this_.Read()
					if this_.ImplicitSemicolon {
						tkn = token.Illegal
					} else {
						tkn = token.Arrow
					}
				} else {
					tkn = this_.Switch2(token.Assign, token.Equal)
					if tkn == token.Equal && this_.Chr == '=' {
						this_.Read()
						tkn = token.StrictEqual
					}
				}
			case '!':
				tkn = this_.Switch2(token.Not, token.NotEqual)
				if tkn == token.NotEqual && this_.Chr == '=' {
					this_.Read()
					tkn = token.StrictNotEqual
				}
			case '&':
				tkn = this_.Switch3(token.And, token.AndAssign, '&', token.LogicalAnd)
			case '|':
				tkn = this_.Switch3(token.Or, token.OrAssign, '|', token.LogicalOr)
			case '~':
				tkn = token.BitwiseNot
			case '?':
				if this_.Chr == '.' && !this_.IsDecimalDigit(this_.ImplicitRead()) {
					this_.Read()
					tkn = token.QuestionDot
				} else if this_.Chr == '?' {
					this_.Read()
					tkn = token.Coalesce
				} else {
					tkn = token.QuestionMark
				}
			case '"', '\'':
				InsertSemicolon = true
				tkn = token.String
				var err string
				literal, parsedLiteral, err = this_.ScanString(this_.ChrOffset-1, true)
				if err != "" {
					tkn = token.Illegal
				}
			case '`':
				tkn = token.Backtick
			case '#':
				if this_.ChrOffset == 1 && this_.Chr == '!' {
					this_.SkipSingleLineComment()
					continue
				}

				var err string
				literal, parsedLiteral, _, err = this_.ScanIdentifier()
				if err != "" || literal == "" {
					tkn = token.Illegal
					break
				}
				this_.InsertSemicolon = true
				tkn = token.PrivateIdentifier
				return
			default:
				_ = this_.ErrorUnexpected("scan chr:"+string(chr), idx, chr)
				tkn = token.Illegal
			}
		}
		this_.InsertSemicolon = InsertSemicolon
		return
	}
}

func (this_ *Parser) ScanMantissa(base int) {
	for this_.DigitValue(this_.Chr) < base {
		this_.Read()
	}
}
