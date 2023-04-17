package java

import (
	"github.com/team-ide/go-interpreter/token"
)

func (this_ *Parser) parseImportStatement() *ImportStatement {

	idx := this_.ExpectAndNext("parseImportStatement", token.Import)

	res := &ImportStatement{
		From: idx,
	}

	//res.Import = this_.ParseChainNameStatement()
	//res.To = res.Import.To

	this_.ExpectAndNext("parseImportStatement", token.Semicolon)
	//fmt.Println("parseImportStatement ", "imp:", imp, ",Next token:", this_.Token)

	return res
}
