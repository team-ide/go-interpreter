package parser

import (
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/token"
)

func (this_ *Parser) ParseChainNameStatement() *node.ChainNameStatement {

	res := &node.ChainNameStatement{
		From: this_.Idx,
	}

	for {
		if this_.Token == token.Period {
			if this_.ImplicitRead() == '.' {
				break
			}
			res.To = this_.Idx + 1
			this_.Next()
			continue
		} else if this_.ParsedLiteral != "" {
			res.Names = append(res.Names, this_.ParsedLiteral)
			res.To = this_.Idx + len(this_.ParsedLiteral)
			this_.Next()
			if this_.Token != token.Period {
				break
			}
		} else {
			break
		}
	}

	return res
}
