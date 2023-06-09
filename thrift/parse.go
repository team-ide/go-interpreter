package thrift

import (
	"fmt"
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/parser"
	"github.com/team-ide/go-interpreter/token"
)

type Parser struct {
	*parser.Parser
}

func Parse(filename, src string) (tree *node.Tree, err error) {
	p := &Parser{
		Parser: parser.New(src),
	}
	p.Filename = filename
	p.ParseStatement = p.parseStatement
	p.KeywordToken = KeywordToken
	p.IdentifierTokens = IdentifierTokens
	p.UnreservedWordTokens = UnreservedWordTokens
	p.ModifierTokens = ModifierTokens
	return p.Parse()
}

func (this_ *Parser) parseStatement() node.Statement {
	if this_.Token == token.Eof {
		_ = this_.ErrorUnexpectedToken("parseStatement this_.Token is token.Eof", this_.Token)
		return &node.BadStatement{From: this_.Idx, To: this_.Idx + 1}
	}

	switch this_.Token {
	case token.Include:
		return this_.parseIncludeStatement()
	case token.Namespace:
		return this_.parseNamespaceStatement()
	case token.Exception:
		return this_.parseExceptionStatement()
	case token.Struct:
		return this_.parseStructStatement()
	case token.Enum:
		return this_.parseEnumStatement()
	case token.Service:
		return this_.parseServiceStatement()
	default:
		fmt.Println("parseStatement this_.Token:", this_.Token)
	}
	this_.Next()
	return &node.BadStatement{}
}
