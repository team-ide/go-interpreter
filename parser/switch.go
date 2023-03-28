package parser

import "github.com/team-ide/go-interpreter/token"

func (this_ *Parser) Switch2(tkn0, tkn1 token.Token) token.Token {
	if this_.Chr == '=' {
		this_.Read()
		return tkn1
	}
	return tkn0
}

func (this_ *Parser) Switch3(tkn0, tkn1 token.Token, chr2 rune, tkn2 token.Token) token.Token {
	if this_.Chr == '=' {
		this_.Read()
		return tkn1
	}
	if this_.Chr == chr2 {
		this_.Read()
		return tkn2
	}
	return tkn0
}

func (this_ *Parser) Switch4(tkn0, tkn1 token.Token, chr2 rune, tkn2, tkn3 token.Token) token.Token {
	if this_.Chr == '=' {
		this_.Read()
		return tkn1
	}
	if this_.Chr == chr2 {
		this_.Read()
		if this_.Chr == '=' {
			this_.Read()
			return tkn3
		}
		return tkn2
	}
	return tkn0
}

func (this_ *Parser) Switch6(tkn0, tkn1 token.Token, chr2 rune, tkn2, tkn3 token.Token, chr3 rune, tkn4, tkn5 token.Token) token.Token {
	if this_.Chr == '=' {
		this_.Read()
		return tkn1
	}
	if this_.Chr == chr2 {
		this_.Read()
		if this_.Chr == '=' {
			this_.Read()
			return tkn3
		}
		if this_.Chr == chr3 {
			this_.Read()
			if this_.Chr == '=' {
				this_.Read()
				return tkn5
			}
			return tkn4
		}
		return tkn2
	}
	return tkn0
}
