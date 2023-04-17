package java

import (
	"fmt"
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/parser"
	"github.com/team-ide/go-interpreter/token"
)

type Parser struct {
	*parser.Parser
}

func Parse(src string) (tree *node.Tree, err error) {
	p := &Parser{
		Parser: parser.New(src),
	}
	p.ParseStatement = p.parseStatement
	p.KeywordToken = KeywordToken
	p.IdentifierTokens = IdentifierTokens
	p.UnreservedWordTokens = UnreservedWordTokens
	p.ModifierTokens = ModifierTokens
	return p.Parse()
}

func (this_ *Parser) parseStatement() node.Statement {
	fmt.Println("parseStatement this_.Token:", this_.Token)
	if this_.Token == token.Eof {
		_ = this_.ErrorUnexpectedToken("parseStatement this_.Token is token.Eof", this_.Token)
		return &node.BadStatement{From: this_.Idx, To: this_.Idx + 1}
	}

	switch this_.Token {
	case token.Import:
		return this_.parseImportStatement()
	}
	this_.Next()
	return &node.BadStatement{
		//Expression: expression,
	}
}
