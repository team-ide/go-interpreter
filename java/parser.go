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
	return p.Parse()
}

func (this_ *Parser) parseStatement() node.Statement {
	fmt.Println("parseStatement this_.Token:", this_.Token)
	if this_.Token == token.Eof {
		_ = this_.ErrorUnexpectedToken("parseStatement this_.Token is token.Eof", this_.Token)
		return &node.BadStatement{From: this_.Idx, To: this_.Idx + 1}
	}

	switch this_.Token {
	case token.Semicolon:
		// 解析 分号
		return this_.ParseSemicolonStatement()
	case token.LeftBrace:
		// 解析 { } 子语句
		return this_.ParseBlockStatement()
	case token.Import:
		return this_.parseImportStatement()
	case token.If:
		//return this_.parseIfStatement()
	case token.Do:
		//return this_.parseDoWhileStatement()
	case token.While:
		//return this_.parseWhileStatement()
	case token.For:
		//return this_.parseForOrForInStatement()
	case token.Break:
		//return this_.parseBreakStatement()
	case token.Continue:
		//return this_.parseContinueStatement()
	case token.Debugger:
		//return this_.parseDebuggerStatement()
	case token.With:
		//return this_.parseWithStatement()
	case token.Function:
		//return &node.FunctionDeclaration{
		//	Function: this_.parseFunction(true, false, this_.Idx),
		//}
	case token.Class:
		//return &node.ClassDeclaration{
		//	Class: this_.parseClass(true),
		//}
	case token.Switch:
		//return this_.parseSwitchStatement()
	case token.Return:
		//return this_.parseReturnStatement()
	case token.Throw:
		//return this_.parseThrowStatement()
	case token.Try:
		//return this_.parseTryStatement()
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
