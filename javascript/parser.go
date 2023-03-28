package javascript

import (
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
	p.KeywordToken = KeywordToken
	p.IdentifierTokens = IdentifierTokens
	p.UnreservedWordTokens = UnreservedWordTokens
	return p.Parse()
}

// Parse 解析
func (this_ *Parser) Parse() (tree *node.Tree, err error) {
	tree = this_.parseTree()
	//this_.errors.Sort()
	err = this_.Errors.Err()
	return
}

func (this_ *Parser) parseTree() (tree *node.Tree) {
	this_.OpenScope()
	defer this_.CloseScope()

	var statements []node.Statement
	this_.Next()
	for this_.Token != token.Eof {
		this_.Scope.AllowLet = true
		statements = append(statements, this_.parseStatement())
	}

	tree = &node.Tree{
		Children:        statements,
		DeclarationList: this_.Scope.DeclarationList,
	}
	//this_.file.SetSourceMap(this_.parseSourceMap())
	return
}
