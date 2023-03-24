package parser

import (
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/token"
)

func (this_ *parser) scan() (tkn token.Token, literal string, parsedLiteral node.String, idx int) {

	this_.implicitSemicolon = false

	for {
		//this_.skipWhiteSpace()

		idx = this_.chrOffset
		insertSemicolon := false

		switch chr := this_.chr; {
		case this_.IsIdentifierStart(chr):
			var err string
			var hasEscape bool
			literal, parsedLiteral, hasEscape, err = this_.scanIdentifier()
			if err != "" {
				tkn = token.Illegal
				break
			}
			if len(parsedLiteral) > 1 {
				// Keywords are longer than 1 character, avoid lookup otherwise
				var strict bool
				tkn, strict = this_.IsKeyword(string(parsedLiteral))
				if hasEscape {
					this_.insertSemicolon = true
					if tkn == "" || this_.isBindingId(tkn) {
						tkn = token.Identifier
					} else {
						tkn = token.EscapedReservedWord
					}
					return
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
					token.Yield,
					token.Return,
					token.Continue,
					token.Debugger:
					this_.insertSemicolon = true
					return

				case token.Async:
					// async only has special meaning if not followed by a LineTerminator
					if this_.skipWhiteSpaceCheckLineTerminator() {
						this_.insertSemicolon = true
						tkn = token.Identifier
					}
					return
				default:
					return

				}
			}
			this_.insertSemicolon = true
			tkn = token.Identifier
			return
		case '0' <= chr && chr <= '9':
			this_.insertSemicolon = true
			tkn, literal = this_.scanNumericLiteral(false)
			return
		default:
			this_.read()
			switch chr {
			case -1:
				if this_.insertSemicolon {
					this_.insertSemicolon = false
					this_.implicitSemicolon = true
				}
				tkn = token.Eof
			case '\r', '\n', '\u2028', '\u2029':
				this_.insertSemicolon = false
				this_.implicitSemicolon = true
				continue
			case ':':
				tkn = token.Colon
			case '.':
				if this_.DigitValue(this_.chr) < 10 {
					insertSemicolon = true
					tkn, literal = this_.scanNumericLiteral(true)
				} else {
					if this_.chr == '.' {
						this_.read()
						if this_.chr == '.' {
							this_.read()
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
				insertSemicolon = true
			case '[':
				tkn = token.LeftBracket
			case ']':
				tkn = token.RightBracket
				insertSemicolon = true
			case '{':
				tkn = token.LeftBrace
			case '}':
				tkn = token.RightBrace
				insertSemicolon = true
			case '+':
				tkn = this_.switch3(token.Plus, token.AddAssign, '+', token.Increment)
				if tkn == token.Increment {
					insertSemicolon = true
				}
			case '-':
				tkn = this_.switch3(token.Minus, token.SubtractAssign, '-', token.Decrement)
				if tkn == token.Decrement {
					insertSemicolon = true
				}
			case '*':
				if this_.chr == '*' {
					this_.read()
					tkn = this_.switch2(token.Exponent, token.ExponentAssign)
				} else {
					tkn = this_.switch2(token.Multiply, token.MultiplyAssign)
				}
			case '/':
				if this_.chr == '/' {
					this_.skipSingleLineComment()
					continue
				} else if this_.chr == '*' {
					if this_.skipMultiLineComment() {
						this_.insertSemicolon = false
						this_.implicitSemicolon = true
					}
					continue
				} else {
					// Could be division, could be RegExp literal
					tkn = this_.switch2(token.Slash, token.QuotientAssign)
					insertSemicolon = true
				}
			case '%':
				tkn = this_.switch2(token.Remainder, token.RemainderAssign)
			case '^':
				tkn = this_.switch2(token.ExclusiveOr, token.ExclusiveOrAssign)
			case '<':
				tkn = this_.switch4(token.Less, token.LessOrEqual, '<', token.ShiftLeft, token.ShiftLeftAssign)
			case '>':
				tkn = this_.switch6(token.Greater, token.GreaterOrEqual, '>', token.ShiftRight, token.ShiftRightAssign, '>', token.UnsignedShiftRight, token.UnsignedShiftRightAssign)
			case '=':
				if this_.chr == '>' {
					this_.read()
					if this_.implicitSemicolon {
						tkn = token.Illegal
					} else {
						tkn = token.Arrow
					}
				} else {
					tkn = this_.switch2(token.Assign, token.Equal)
					if tkn == token.Equal && this_.chr == '=' {
						this_.read()
						tkn = token.StrictEqual
					}
				}
			case '!':
				tkn = this_.switch2(token.Not, token.NotEqual)
				if tkn == token.NotEqual && this_.chr == '=' {
					this_.read()
					tkn = token.StrictNotEqual
				}
			case '&':
				tkn = this_.switch3(token.And, token.AndAssign, '&', token.LogicalAnd)
			case '|':
				tkn = this_.switch3(token.Or, token.OrAssign, '|', token.LogicalOr)
			case '~':
				tkn = token.BitwiseNot
			case '?':
				if this_.chr == '.' && !this_.IsDecimalDigit(this_.implicitRead()) {
					this_.read()
					tkn = token.QuestionDot
				} else if this_.chr == '?' {
					this_.read()
					tkn = token.Coalesce
				} else {
					tkn = token.QuestionMark
				}
			case '"', '\'':
				insertSemicolon = true
				tkn = token.String
				var err string
				literal, parsedLiteral, err = this_.scanString(this_.chrOffset-1, true)
				if err != "" {
					tkn = token.Illegal
				}
			case '`':
				tkn = token.Backtick
			case '#':
				if this_.chrOffset == 1 && this_.chr == '!' {
					this_.skipSingleLineComment()
					continue
				}

				var err string
				literal, parsedLiteral, _, err = this_.scanIdentifier()
				if err != "" || literal == "" {
					tkn = token.Illegal
					break
				}
				this_.insertSemicolon = true
				tkn = token.PrivateIdentifier
				return
			default:
				_ = this_.errorUnexpected(idx, chr)
				tkn = token.Illegal
			}
		}
		this_.insertSemicolon = insertSemicolon
		return
	}
}
