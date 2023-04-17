package java

import (
	"github.com/team-ide/go-interpreter/parser"
	"github.com/team-ide/go-interpreter/token"
)

var (
	KeywordToken = map[string]parser.Keyword{
		"import": {
			Token: token.Import,
		},
	}

	IdentifierTokens = []token.Token{}

	UnreservedWordTokens = []token.Token{}

	ModifierTokens = []token.Token{}
)
