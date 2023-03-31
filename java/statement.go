package java

import (
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/token"
)

func (this_ *Parser) parseImportStatement() *node.ImportStatement {

	idx := this_.ExpectAndNext("parseImportStatement", token.Import)

	res := &node.ImportStatement{
		From: idx,
	}

	imp := ""
	for {
		if this_.Token == token.Identifier {
			identifier := this_.ParseIdentifier()
			imp += string(identifier.Name)
		} else if this_.Token == token.Period {
			this_.Next()
			imp += "."
		} else if this_.Token == token.Multiply {
			this_.Next()
			imp += "*"
		} else {
			break
		}
	}
	res.Import = imp
	res.To = this_.Idx
	if token.Semicolon == this_.Token {
		res.To++
	}
	this_.ExpectAndNext("parseImportStatement", token.Semicolon)
	//fmt.Println("parseImportStatement ", "imp:", imp, ",Next token:", this_.Token)

	return res
}
