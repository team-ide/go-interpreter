package golang

import (
	"github.com/team-ide/go-interpreter/parser"
	"github.com/team-ide/go-interpreter/token"
)

var (
	KeywordToken = map[string]parser.Keyword{}

	IdentifierTokens = []token.Token{}

	UnreservedWordTokens = []token.Token{}
)
