package parser

import "github.com/team-ide/go-interpreter/token"

func (this_ *parser) switch2(tkn0, tkn1 token.Token) token.Token {
	if this_.chr == '=' {
		this_.read()
		return tkn1
	}
	return tkn0
}

func (this_ *parser) switch3(tkn0, tkn1 token.Token, chr2 rune, tkn2 token.Token) token.Token {
	if this_.chr == '=' {
		this_.read()
		return tkn1
	}
	if this_.chr == chr2 {
		this_.read()
		return tkn2
	}
	return tkn0
}

func (this_ *parser) switch4(tkn0, tkn1 token.Token, chr2 rune, tkn2, tkn3 token.Token) token.Token {
	if this_.chr == '=' {
		this_.read()
		return tkn1
	}
	if this_.chr == chr2 {
		this_.read()
		if this_.chr == '=' {
			this_.read()
			return tkn3
		}
		return tkn2
	}
	return tkn0
}

func (this_ *parser) switch6(tkn0, tkn1 token.Token, chr2 rune, tkn2, tkn3 token.Token, chr3 rune, tkn4, tkn5 token.Token) token.Token {
	if this_.chr == '=' {
		this_.read()
		return tkn1
	}
	if this_.chr == chr2 {
		this_.read()
		if this_.chr == '=' {
			this_.read()
			return tkn3
		}
		if this_.chr == chr3 {
			this_.read()
			if this_.chr == '=' {
				this_.read()
				return tkn5
			}
			return tkn4
		}
		return tkn2
	}
	return tkn0
}
