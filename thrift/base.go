package thrift

import (
	"github.com/team-ide/go-interpreter/parser"
	"github.com/team-ide/go-interpreter/token"
)

var (
	KeywordToken = map[string]parser.Keyword{
		"include": {
			Token: token.Include,
		},

		"namespace": {
			Token: token.Namespace,
		},
		"struct": {
			Token: token.Struct,
		},
		"exception": {
			Token: token.Exception,
		},
		"enum": {
			Token: token.Enum,
		},
		"service": {
			Token: token.Service,
		},
		"optional": {
			Token: token.Optional,
		},
	}

	IdentifierTokens = []token.Token{
		token.Identifier,
		token.Keyword,
		token.Include,
		token.Namespace,
		token.Struct,
		token.Exception,
		token.Enum,
		token.Service,
		token.Optional,
	}

	UnreservedWordTokens = []token.Token{}

	ModifierTokens = []token.Token{}
)
