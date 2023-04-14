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
	if this_.Token == token.Eof {
		_ = this_.ErrorUnexpectedToken("parseStatement this_.Token is token.Eof", this_.Token)
		return &node.BadStatement{From: this_.Idx, To: this_.Idx + 1}
	}

	switch this_.Token {
	case token.Semicolon:
		// 解析 分号
		return this_.ParseSemicolonStatement()
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
		var identifier *node.Identifier
		if this_.Token == token.Identifier {
			identifier = this_.ParseIdentifier()
		}
		fmt.Println("parseStatement this_.Token:", this_.Token, ",identifier:", identifier)
	}
	this_.Next()
	//expression := this_.parseExpression()
	//
	//if identifier, isIdentifier := expression.(*node.Identifier); isIdentifier && this_.Token == token.Colon {
	//	// LabelledStatement
	//	colon := this_.Idx
	//	this_.Next() // :
	//	label := identifier.Name
	//	for _, value := range this_.Scope.Labels {
	//		if label == value {
	//			_ = this_.Error("parseStatement", identifier.Start(), fmt.Sprintf("Label '%s' already exists", label))
	//		}
	//	}
	//	this_.Scope.Labels = append(this_.Scope.Labels, label) // Push the label
	//	this_.Scope.AllowLet = false
	//	statement := this_.parseStatement()
	//	this_.Scope.Labels = this_.Scope.Labels[:len(this_.Scope.Labels)-1] // Pop the label
	//	return &node.LabelledStatement{
	//		Label:     identifier,
	//		Colon:     colon,
	//		Statement: statement,
	//	}
	//}
	//
	//this_.OptionalSemicolon()

	//bs, _ := json.Marshal(expression)
	//fmt.Println("expression type:", reflect.TypeOf(expression).String(), ",value:", this_.Slice(expression.Start(), expression.End()), ",data:", string(bs))
	return &node.ExpressionStatement{
		//Expression: expression,
	}
}
