package golang

import (
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/parser"
)

type Parser struct {
	*parser.Parser
}

func Parse(src string) (tree *node.Tree, err error) {
	p := &Parser{
		Parser: parser.New(src),
	}
	p.KeywordToken = KeywordToken
	p.IdentifierTokens = IdentifierTokens
	p.UnreservedWordTokens = UnreservedWordTokens
	return p.Parse()
}
