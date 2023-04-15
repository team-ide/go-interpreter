package javascript

import (
	"fmt"
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/parser"
	"github.com/team-ide/go-interpreter/token"
	"reflect"
	"strings"
)

func (this_ *Parser) parsePrimaryExpression() node.Expression {
	literal, parsedLiteral := this_.Literal, this_.ParsedLiteral
	idx := this_.Idx
	switch this_.Token {
	case token.Identifier:
		this_.Next()
		return &node.Identifier{
			Name: parsedLiteral,
			Idx:  idx,
		}
	case token.Null:
		this_.Next()
		return &node.NullLiteral{
			Idx:     idx,
			Literal: literal,
		}
	case token.Boolean:
		this_.Next()
		value := false
		switch parsedLiteral {
		case "true":
			value = true
		case "false":
			value = false
		default:
			_ = this_.Error("parsePrimaryExpression parsedLiteral:"+string(parsedLiteral), idx, "Illegal boolean literal")
		}
		return &node.BooleanLiteral{
			Idx:     idx,
			Literal: literal,
			Value:   value,
		}
	case token.String:
		this_.Next()
		return &node.StringLiteral{
			Idx:     idx,
			Literal: literal,
			Value:   parsedLiteral,
		}
	case token.Number:
		this_.Next()
		value, err := this_.ParseNumberLiteral(literal)
		if err != nil {
			_ = this_.Error("parsePrimaryExpression parseNumberLiteral error literal:"+string(literal), idx, err.Error())
			value = 0
		}
		return &node.NumberLiteral{
			Idx:     idx,
			Literal: literal,
			Value:   value,
		}
	case token.Slash, token.QuotientAssign:
		return this_.parseRegExpLiteral()
	case token.LeftBrace:
		return this_.parseObjectLiteral()
	case token.LeftBracket:
		return this_.parseArrayLiteral()
	case token.LeftParenthesis:
		return this_.parseParenthesisedExpression()
	case token.Backtick:
		return this_.parseTemplateLiteral(false)
	case token.This:
		this_.Next()
		return &node.ThIsExpressionNode{
			Idx: idx,
		}
	case token.Super:
		return this_.parseSuperProperty()
	case token.Async:
		if f := this_.parseMaybeAsyncFunction(false); f != nil {
			return f
		}
	case token.Function:
		return this_.parseFunction(false, false, idx)
	case token.Class:
		return this_.parseClass(false)
	}

	if this_.IsBindingIdentifier(this_.Token) {
		this_.Next()
		return &node.Identifier{
			Name: parsedLiteral,
			Idx:  idx,
		}
	}

	_ = this_.ErrorUnexpectedToken("parsePrimaryExpression", this_.Token)
	this_.nextStatement()
	return &node.BadExpression{From: idx, To: this_.Idx}
}

func (this_ *Parser) parseSuperProperty() node.Expression {
	idx := this_.Idx
	this_.Next()
	switch this_.Token {
	case token.Period:
		this_.Next()
		if !this_.IsIdentifierToken(this_.Token) {
			this_.ExpectAndNext("parseSuperProperty", token.Identifier)
			this_.nextStatement()
			return &node.BadExpression{From: idx, To: this_.Idx}
		}
		idIdx := this_.Idx
		parsedLiteral := this_.ParsedLiteral
		this_.Next()
		return &node.DotExpression{
			Left: &node.SuperExpression{
				Idx: idx,
			},
			Identifier: node.Identifier{
				Name: parsedLiteral,
				Idx:  idIdx,
			},
		}
	case token.LeftBracket:
		return this_.parseBracketMember(&node.SuperExpression{
			Idx: idx,
		})
	case token.LeftParenthesis:
		return this_.parseCallExpression(&node.SuperExpression{
			Idx: idx,
		})
	default:
		_ = this_.Error("parseSuperProperty this_.Token:"+this_.Token.String(), idx, "'super' keyword unexpected here")
		this_.nextStatement()
		return &node.BadExpression{From: idx, To: this_.Idx}
	}
}

func (this_ *Parser) reinterpretSequenceAsArrowFuncParams(list []node.Expression) *node.ParameterList {
	firstRestIdx := -1
	params := make([]*node.Binding, 0, len(list))
	for i, item := range list {
		if _, ok := item.(*node.SpreadElement); ok {
			if firstRestIdx == -1 {
				firstRestIdx = i
				continue
			}
		}
		if firstRestIdx != -1 {
			_ = this_.Error("reinterpretSequenceAsArrowFuncParams firstRestIdx != -1 firstRestIdx:"+fmt.Sprintf("%d", firstRestIdx), list[firstRestIdx].Start(), "Rest parameter must be last formal parameter")
			return &node.ParameterList{}
		}
		params = append(params, this_.reinterpretAsBinding(item))
	}
	var rest node.Expression
	if firstRestIdx != -1 {
		rest = this_.reinterpretAsBindingRestElement(list[firstRestIdx])
	}
	return &node.ParameterList{
		List: params,
		Rest: rest,
	}
}

func (this_ *Parser) parseParenthesisedExpression() node.Expression {
	opening := this_.Idx
	this_.ExpectAndNext("parseParenthesisedExpression", token.LeftParenthesis)
	var list []node.Expression
	if this_.Token != token.RightParenthesis {
		for {
			if this_.Token == token.Ellipsis {
				start := this_.Idx
				_ = this_.ErrorUnexpectedToken("parseParenthesisedExpression", token.Ellipsis)
				this_.Next()
				expr := this_.parseAssignmentExpression()
				list = append(list, &node.BadExpression{
					From: start,
					To:   expr.End(),
				})
			} else {
				list = append(list, this_.parseAssignmentExpression())
			}
			if this_.Token != token.Comma {
				break
			}
			this_.Next()
			if this_.Token == token.RightParenthesis {
				_ = this_.ErrorUnexpectedToken("parseParenthesisedExpression", token.RightParenthesis)
				break
			}
		}
	}
	this_.ExpectAndNext("parseParenthesisedExpression", token.RightParenthesis)
	if len(list) == 1 && len(this_.Errors) == 0 {
		return list[0]
	}
	if len(list) == 0 {
		_ = this_.ErrorUnexpectedToken("parseParenthesisedExpression", token.RightParenthesis)
		return &node.BadExpression{
			From: opening,
			To:   this_.Idx,
		}
	}
	return &node.SequenceExpression{
		Sequence: list,
	}
}

func (this_ *Parser) parseRegExpLiteral() *node.RegExpLiteral {

	offset := this_.ChrOffset - 1 // Opening slash already gotten
	if this_.Token == token.QuotientAssign {
		offset -= 1 // =
	}
	idx := this_.IdxOf(offset)

	pattern, _, err := this_.ScanString(offset, false)
	endOffset := this_.ChrOffset

	if err == "" {
		pattern = pattern[1 : len(pattern)-1]
	}

	flags := ""
	if !this_.IsLineTerminator(this_.Chr) && !this_.IsLineWhiteSpace(this_.Chr) {
		this_.Next()

		if this_.Token == token.Identifier { // gim

			flags = this_.Literal
			this_.Next()
			endOffset = this_.ChrOffset - 1
		}
	} else {
		this_.Next()
	}

	literal := this_.Str[offset:endOffset]

	return &node.RegExpLiteral{
		Idx:     idx,
		Literal: literal,
		Pattern: pattern,
		Flags:   flags,
	}
}

// TokenToBindingIdentifier 如果当前 Token 是 BindingIdentifier 则 将当前 Token 设置为 Identifier
func (this_ *Parser) TokenToBindingIdentifier() {
	if this_.IsBindingIdentifier(this_.Token) {
		this_.Token = token.Identifier
	}
}

func (this_ *Parser) parseBindingTarget() (target node.BindingTarget) {
	this_.TokenToBindingIdentifier()
	switch this_.Token {
	case token.Identifier:
		target = &node.Identifier{
			Name: this_.ParsedLiteral,
			Idx:  this_.Idx,
		}
		this_.Next()
	case token.LeftBracket:
		target = this_.parseArrayBindingPattern()
	case token.LeftBrace:
		target = this_.parseObjectBindingPattern()
	default:
		idx := this_.ExpectAndNext("parseBindingTarget", token.Identifier)
		this_.nextStatement()
		target = &node.BadExpression{From: idx, To: this_.Idx}
	}

	return
}

func (this_ *Parser) parseVariableDeclaration(declarationList *[]*node.Binding) *node.Binding {
	res := &node.Binding{
		Target: this_.parseBindingTarget(),
	}

	if declarationList != nil {
		*declarationList = append(*declarationList, res)
	}

	if this_.Token == token.Assign {
		this_.Next()
		res.Initializer = this_.parseAssignmentExpression()
	}

	return res
}

func (this_ *Parser) parseVariableDeclarationList() (declarationList []*node.Binding) {
	for {
		this_.parseVariableDeclaration(&declarationList)
		if this_.Token != token.Comma {
			break
		}
		this_.Next()
	}
	return
}

func (this_ *Parser) parseVarDeclarationList(var_ int) []*node.Binding {
	declarationList := this_.parseVariableDeclarationList()

	this_.Scope.Declare(&node.VariableDeclaration{
		Var:  var_,
		List: declarationList,
	})

	return declarationList
}

func (this_ *Parser) parseObjectPropertyKey() (string, string, node.Expression, token.Token) {
	if this_.Token == token.LeftBracket {
		this_.Next()
		expr := this_.parseAssignmentExpression()
		this_.ExpectAndNext("parseObjectPropertyKey", token.RightBracket)
		return "", "", expr, token.Illegal
	}
	idx, tkn, literal, parsedLiteral := this_.Idx, this_.Token, this_.Literal, this_.ParsedLiteral
	var value node.Expression
	this_.Next()
	switch tkn {
	case token.Identifier, token.String, token.Keyword, token.EscapedReservedWord:
		value = &node.StringLiteral{
			Idx:     idx,
			Literal: literal,
			Value:   parsedLiteral,
		}
	case token.Number:
		num, err := this_.ParseNumberLiteral(literal)
		if err != nil {
			_ = this_.Error("parseObjectPropertyKey parseNumberLiteral literal:"+string(literal), idx, err.Error())
		} else {
			value = &node.NumberLiteral{
				Idx:     idx,
				Literal: literal,
				Value:   num,
			}
		}
	case token.PrivateIdentifier:
		value = &node.PrivateIdentifier{
			Identifier: node.Identifier{
				Idx:  idx,
				Name: parsedLiteral,
			},
		}
	default:
		// null, false, class, etc.
		if this_.IsIdentifierToken(tkn) {
			value = &node.StringLiteral{
				Idx:     idx,
				Literal: literal,
				Value:   literal,
			}
		} else {
			_ = this_.ErrorUnexpectedToken("parseObjectPropertyKey not IsIdentifierToken:"+tkn.String(), tkn)
		}
	}
	return literal, parsedLiteral, value, tkn
}

func (this_ *Parser) parseObjectProperty() node.Property {
	if this_.Token == token.Ellipsis {
		this_.Next()
		return &node.SpreadElement{
			Expression: this_.parseAssignmentExpression(),
		}
	}
	keyStartIdx := this_.Idx
	generator := false
	if this_.Token == token.Multiply {
		generator = true
		this_.Next()
	}
	literal, parsedLiteral, value, tkn := this_.parseObjectPropertyKey()
	if value == nil {
		return nil
	}
	if this_.IsIdentifierToken(tkn) || tkn == token.String || tkn == token.Number || tkn == token.Illegal {
		if generator {
			return &node.PropertyKeyed{
				Key:      value,
				Kind:     node.PropertyKindMethod,
				Value:    this_.parseMethodDefinition(keyStartIdx, node.PropertyKindMethod, true, false),
				Computed: tkn == token.Illegal,
			}
		}
		switch {
		case this_.Token == token.LeftParenthesis:
			return &node.PropertyKeyed{
				Key:      value,
				Kind:     node.PropertyKindMethod,
				Value:    this_.parseMethodDefinition(keyStartIdx, node.PropertyKindMethod, false, false),
				Computed: tkn == token.Illegal,
			}
		case this_.Token == token.Comma || this_.Token == token.RightBrace || this_.Token == token.Assign: // shorthand property
			if this_.IsBindingIdentifier(tkn) {
				var initializer node.Expression
				if this_.Token == token.Assign {
					// allow the initializer syntax here in case the object literal
					// needs to be reinterpreted as an assignment pattern, enforce later if it doesn't.
					this_.Next()
					initializer = this_.parseAssignmentExpression()
				}
				return &node.PropertyShort{
					Name: node.Identifier{
						Name: parsedLiteral,
						Idx:  value.Start(),
					},
					Initializer: initializer,
				}
			} else {
				_ = this_.ErrorUnexpectedToken("parseObjectProperty not this_.isBindingId:"+tkn.String(), this_.Token)
			}
		case (literal == "get" || literal == "set" || tkn == token.Async) && this_.Token != token.Colon:
			_, _, keyValue, tkn1 := this_.parseObjectPropertyKey()
			if keyValue == nil {
				return nil
			}

			var kind node.PropertyKind
			var async bool
			if tkn == token.Async {
				async = true
				kind = node.PropertyKindMethod
			} else if literal == "get" {
				kind = node.PropertyKindGet
			} else {
				kind = node.PropertyKindSet
			}

			return &node.PropertyKeyed{
				Key:      keyValue,
				Kind:     kind,
				Value:    this_.parseMethodDefinition(keyStartIdx, kind, false, async),
				Computed: tkn1 == token.Illegal,
			}
		}
	}

	this_.ExpectAndNext("parseObjectProperty", token.Colon)
	return &node.PropertyKeyed{
		Key:      value,
		Kind:     node.PropertyKindValue,
		Value:    this_.parseAssignmentExpression(),
		Computed: tkn == token.Illegal,
	}
}

func (this_ *Parser) parseMethodDefinition(keyStartIdx int, kind node.PropertyKind, generator, async bool) *node.FunctionLiteral {
	idx1 := this_.Idx
	if generator != this_.Scope.AllowYield {
		this_.Scope.AllowYield = generator
		defer func() {
			this_.Scope.AllowYield = !generator
		}()
	}
	if async != this_.Scope.AllowAwait {
		this_.Scope.AllowAwait = async
		defer func() {
			this_.Scope.AllowAwait = !async
		}()
	}
	parameterList := this_.parseFunctionParameterList()
	switch kind {
	case node.PropertyKindGet:
		if len(parameterList.List) > 0 || parameterList.Rest != nil {
			_ = this_.Error("parseMethodDefinition node.PropertyKindGet", idx1, "Getter must not have any formal parameters.")
		}
	case node.PropertyKindSet:
		if len(parameterList.List) != 1 || parameterList.Rest != nil {
			_ = this_.Error("parseMethodDefinition node.PropertyKindSet", idx1, "Setter must have exactly one formal parameter.")
		}
	}
	res := &node.FunctionLiteral{
		Function:      keyStartIdx,
		ParameterList: parameterList,
		Generator:     generator,
		Async:         async,
	}
	res.Body, res.DeclarationList = this_.parseFunctionBlock(async, async, generator)
	res.Source = this_.Slice(keyStartIdx, res.Body.End())
	return res
}

func (this_ *Parser) parseObjectLiteral() *node.ObjectLiteral {
	var value []node.Property
	idx0 := this_.ExpectAndNext("parseObjectLiteral", token.LeftBrace)
	for this_.Token != token.RightBrace && this_.Token != token.Eof {
		property := this_.parseObjectProperty()
		if property != nil {
			value = append(value, property)
		}
		if this_.Token != token.RightBrace {
			this_.ExpectAndNext("parseObjectLiteral", token.Comma)
		} else {
			break
		}
	}
	idx1 := this_.ExpectAndNext("parseObjectLiteral", token.RightBrace)

	return &node.ObjectLiteral{
		LeftBrace:  idx0,
		RightBrace: idx1,
		Value:      value,
	}
}

func (this_ *Parser) parseArrayLiteral() *node.ArrayLiteral {

	idx0 := this_.ExpectAndNext("parseArrayLiteral", token.LeftBracket)
	var value []node.Expression
	for this_.Token != token.RightBracket && this_.Token != token.Eof {
		if this_.Token == token.Comma {
			this_.Next()
			value = append(value, nil)
			continue
		}
		if this_.Token == token.Ellipsis {
			this_.Next()
			value = append(value, &node.SpreadElement{
				Expression: this_.parseAssignmentExpression(),
			})
		} else {
			value = append(value, this_.parseAssignmentExpression())
		}
		if this_.Token != token.RightBracket {
			this_.ExpectAndNext("parseArrayLiteral", token.Comma)
		}
	}
	idx1 := this_.ExpectAndNext("parseArrayLiteral", token.RightBracket)

	return &node.ArrayLiteral{
		LeftBracket:  idx0,
		RightBracket: idx1,
		Value:        value,
	}
}

func (this_ *Parser) parseTemplateLiteral(tagged bool) *node.TemplateLiteral {
	res := &node.TemplateLiteral{
		OpenQuote: this_.Idx,
	}
	for {
		start := this_.Offset
		literal, parsed, finished, parseErr, err := this_.ParseTemplateCharacters()
		if err != "" {
			_ = this_.Error("parseTemplateLiteral parseTemplateCharacters err", this_.Offset, err)
		}
		res.Elements = append(res.Elements, &node.TemplateElement{
			Idx:     this_.IdxOf(start),
			Literal: literal,
			Parsed:  parsed,
			Valid:   parseErr == "",
		})
		if !tagged && parseErr != "" {
			_ = this_.Error("parseTemplateLiteral parseTemplateCharacters parseErr", this_.Offset, parseErr)
		}
		end := this_.ChrOffset - 1
		this_.Next()
		if finished {
			res.CloseQuote = this_.IdxOf(end)
			break
		}
		expr := this_.parseExpression()
		res.Expressions = append(res.Expressions, expr)
		if this_.Token != token.RightBrace {
			_ = this_.ErrorUnexpectedToken("parseTemplateLiteral this_.Token:"+this_.Token.String()+" is not token.RightBrace:"+token.RightBrace.String(), this_.Token)
		}
	}
	return res
}

func (this_ *Parser) parseTaggedTemplateLiteral(tag node.Expression) *node.TemplateLiteral {
	l := this_.parseTemplateLiteral(true)
	l.Tag = tag
	return l
}

func (this_ *Parser) parseArgumentList() (argumentList []node.Expression, idx0, idx1 int) {
	idx0 = this_.ExpectAndNext("parseArgumentList", token.LeftParenthesis)
	for this_.Token != token.RightParenthesis {
		var item node.Expression
		if this_.Token == token.Ellipsis {
			this_.Next()
			item = &node.SpreadElement{
				Expression: this_.parseAssignmentExpression(),
			}
		} else {
			item = this_.parseAssignmentExpression()
		}
		argumentList = append(argumentList, item)
		if this_.Token != token.Comma {
			break
		}
		this_.Next()
	}
	idx1 = this_.ExpectAndNext("parseArgumentList", token.RightParenthesis)
	return
}

func (this_ *Parser) parseCallExpression(left node.Expression) node.Expression {
	argumentList, idx0, idx1 := this_.parseArgumentList()
	return &node.CallExpression{
		Callee:           left,
		LeftParenthesis:  idx0,
		ArgumentList:     argumentList,
		RightParenthesis: idx1,
	}
}

func (this_ *Parser) parseDotMember(left node.Expression) node.Expression {
	period := this_.Idx
	this_.Next()

	literal := this_.ParsedLiteral
	idx := this_.Idx

	if this_.Token == token.PrivateIdentifier {
		this_.Next()
		return &node.PrivateDotExpression{
			Left: left,
			Identifier: node.PrivateIdentifier{
				Identifier: node.Identifier{
					Idx:  idx,
					Name: literal,
				},
			},
		}
	}

	if !this_.IsIdentifierToken(this_.Token) {
		this_.ExpectAndNext("parseDotMember", token.Identifier)
		this_.nextStatement()
		return &node.BadExpression{From: period, To: this_.Idx}
	}

	this_.Next()

	return &node.DotExpression{
		Left: left,
		Identifier: node.Identifier{
			Idx:  idx,
			Name: literal,
		},
	}
}

func (this_ *Parser) parseBracketMember(left node.Expression) node.Expression {
	idx0 := this_.ExpectAndNext("parseBracketMember", token.LeftBracket)
	member := this_.parseExpression()
	idx1 := this_.ExpectAndNext("parseBracketMember", token.RightBracket)
	return &node.BracketExpression{
		LeftBracket:  idx0,
		Left:         left,
		Member:       member,
		RightBracket: idx1,
	}
}

func (this_ *Parser) parseNewExpression() node.Expression {
	idx := this_.ExpectAndNext("parseNewExpression", token.New)
	if this_.Token == token.Period {
		this_.Next()
		if this_.Literal == "target" {
			return &node.MetaProperty{
				Meta: &node.Identifier{
					Name: token.New.String(),
					Idx:  idx,
				},
				Property: this_.ParseIdentifier(),
			}
		}
		_ = this_.ErrorUnexpectedToken("parseNewExpression", token.Identifier)
	}
	callee := this_.parseLeftHandSideExpression()
	if bad, ok := callee.(*node.BadExpression); ok {
		bad.From = idx
		return bad
	}
	res := &node.NewExpression{
		New:    idx,
		Callee: callee,
	}
	if this_.Token == token.LeftParenthesis {
		argumentList, idx0, idx1 := this_.parseArgumentList()
		res.ArgumentList = argumentList
		res.LeftParenthesis = idx0
		res.RightParenthesis = idx1
	}
	return res
}

func (this_ *Parser) parseLeftHandSideExpression() node.Expression {

	var left node.Expression
	if this_.Token == token.New {
		left = this_.parseNewExpression()
	} else {
		left = this_.parsePrimaryExpression()
	}
L:
	for {
		switch this_.Token {
		case token.Period:
			left = this_.parseDotMember(left)
		case token.LeftBracket:
			left = this_.parseBracketMember(left)
		case token.Backtick:
			left = this_.parseTaggedTemplateLiteral(left)
		default:
			break L
		}
	}

	return left
}

func (this_ *Parser) parseLeftHandSideExpressionAllowCall() node.Expression {

	allowIn := this_.Scope.AllowIn
	this_.Scope.AllowIn = true
	defer func() {
		this_.Scope.AllowIn = allowIn
	}()

	var left node.Expression
	start := this_.Idx
	if this_.Token == token.New {
		left = this_.parseNewExpression()
	} else {
		left = this_.parsePrimaryExpression()
	}

	optionalChain := false
L:
	for {
		switch this_.Token {
		case token.Period:
			left = this_.parseDotMember(left)
		case token.LeftBracket:
			left = this_.parseBracketMember(left)
		case token.LeftParenthesis:
			left = this_.parseCallExpression(left)
		case token.Backtick:
			if optionalChain {
				_ = this_.Error("parseLeftHandSideExpressionAllowCall token.Backtick optionalChain:true", this_.Idx, "Invalid template literal on optional chain")
				this_.nextStatement()
				return &node.BadExpression{From: start, To: this_.Idx}
			}
			left = this_.parseTaggedTemplateLiteral(left)
		case token.QuestionDot:
			optionalChain = true
			left = &node.Optional{Expression: left}

			switch this_.Peek() {
			case token.LeftBracket, token.LeftParenthesis, token.Backtick:
				this_.Next()
			default:
				left = this_.parseDotMember(left)
			}
		default:
			break L
		}
	}

	if optionalChain {
		left = &node.OptionalChain{Expression: left}
	}
	return left
}

func (this_ *Parser) parsePostfixExpression() node.Expression {
	operand := this_.parseLeftHandSideExpressionAllowCall()

	//fmt.Println("parsePostfixExpression start:", operand.Start(), ",end:", operand.End(), ",operand:", this_.slice(operand.Start(), operand.End()))
	switch this_.Token {
	case token.Increment, token.Decrement:
		// Make sure there is no line terminator here
		if this_.ImplicitSemicolon {
			break
		}
		tkn := this_.Token
		idx := this_.Idx
		this_.Next()
		switch operand.(type) {
		case *node.Identifier, *node.DotExpression, *node.PrivateDotExpression, *node.BracketExpression:
		default:
			_ = this_.Error("parsePostfixExpression operand type:"+reflect.TypeOf(operand).String(), idx, "Invalid left-hand side in assignment")
			this_.nextStatement()
			return &node.BadExpression{From: idx, To: this_.Idx}
		}
		return &node.UnaryExpression{
			Operator: tkn,
			Idx:      idx,
			Operand:  operand,
			Postfix:  true,
		}
	}

	return operand
}

func (this_ *Parser) parseUnaryExpression() node.Expression {

	switch this_.Token {
	case token.Plus, token.Minus, token.Not, token.BitwiseNot:
		fallthrough
	case token.Delete, token.Void, token.Typeof:
		tkn := this_.Token
		idx := this_.Idx
		this_.Next()
		return &node.UnaryExpression{
			Operator: tkn,
			Idx:      idx,
			Operand:  this_.parseUnaryExpression(),
		}
	case token.Increment, token.Decrement:
		tkn := this_.Token
		idx := this_.Idx
		this_.Next()
		operand := this_.parseUnaryExpression()
		switch operand.(type) {
		case *node.Identifier, *node.DotExpression, *node.PrivateDotExpression, *node.BracketExpression:
		default:
			_ = this_.Error("parseUnaryExpression operand type:"+reflect.TypeOf(operand).String(), idx, "Invalid left-hand side in assignment")
			this_.nextStatement()
			return &node.BadExpression{From: idx, To: this_.Idx}
		}
		return &node.UnaryExpression{
			Operator: tkn,
			Idx:      idx,
			Operand:  operand,
		}
	case token.Await:
		if this_.Scope.AllowAwait {
			idx := this_.Idx
			this_.Next()
			if !this_.Scope.InAsync {
				_ = this_.ErrorUnexpectedToken("parseUnaryExpression is not this_.scope.inAsync", token.Await)
				return &node.BadExpression{
					From: idx,
					To:   this_.Idx,
				}
			}
			if this_.Scope.InFuncParams {
				_ = this_.Error("parseUnaryExpression this_.scope.inFuncParams", idx, "Illegal await-expression in formal parameters of async function")
			}
			return &node.AwaitExpression{
				Await:    idx,
				Argument: this_.parseUnaryExpression(),
			}
		}
	}

	return this_.parsePostfixExpression()
}

func isUpdateExpression(expr node.Expression) bool {
	if ux, ok := expr.(*node.UnaryExpression); ok {
		return ux.Operator == token.Increment || ux.Operator == token.Decrement
	}
	return true
}

func (this_ *Parser) parseExponentiationExpression() node.Expression {
	left := this_.parseUnaryExpression()

	for this_.Token == token.Exponent && isUpdateExpression(left) {
		this_.Next()
		left = &node.BinaryExpression{
			Operator: token.Exponent,
			Left:     left,
			Right:    this_.parseExponentiationExpression(),
		}
	}

	return left
}

func (this_ *Parser) parseMultiplicativeExpression() node.Expression {
	left := this_.parseExponentiationExpression()

	for this_.Token == token.Multiply || this_.Token == token.Slash ||
		this_.Token == token.Remainder {
		tkn := this_.Token
		this_.Next()
		left = &node.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    this_.parseExponentiationExpression(),
		}
	}

	return left
}

func (this_ *Parser) parseAdditiveExpression() node.Expression {
	left := this_.parseMultiplicativeExpression()

	for this_.Token == token.Plus || this_.Token == token.Minus {
		tkn := this_.Token
		this_.Next()
		left = &node.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    this_.parseMultiplicativeExpression(),
		}
	}

	return left
}

func (this_ *Parser) parseShiftExpression() node.Expression {
	left := this_.parseAdditiveExpression()

	for this_.Token == token.ShiftLeft || this_.Token == token.ShiftRight ||
		this_.Token == token.UnsignedShiftRight {
		tkn := this_.Token
		this_.Next()
		left = &node.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    this_.parseAdditiveExpression(),
		}
	}

	return left
}

func (this_ *Parser) parseRelationalExpression() node.Expression {
	if this_.Scope.AllowIn && this_.Token == token.PrivateIdentifier {
		left := &node.PrivateIdentifier{
			Identifier: node.Identifier{
				Idx:  this_.Idx,
				Name: this_.ParsedLiteral,
			},
		}
		this_.Next()
		if this_.Token == token.In {
			this_.Next()
			return &node.BinaryExpression{
				Operator: this_.Token,
				Left:     left,
				Right:    this_.parseShiftExpression(),
			}
		}
		return left
	}
	left := this_.parseShiftExpression()

	allowIn := this_.Scope.AllowIn
	this_.Scope.AllowIn = true
	defer func() {
		this_.Scope.AllowIn = allowIn
	}()

	switch this_.Token {
	case token.Less, token.LessOrEqual, token.Greater, token.GreaterOrEqual:
		tkn := this_.Token
		this_.Next()
		return &node.BinaryExpression{
			Operator:   tkn,
			Left:       left,
			Right:      this_.parseRelationalExpression(),
			Comparison: true,
		}
	case token.Instanceof:
		tkn := this_.Token
		this_.Next()
		return &node.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    this_.parseRelationalExpression(),
		}
	case token.In:
		if !allowIn {
			return left
		}
		tkn := this_.Token
		this_.Next()
		return &node.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    this_.parseRelationalExpression(),
		}
	}

	return left
}

func (this_ *Parser) parseEqualityExpression() node.Expression {
	left := this_.parseRelationalExpression()

	for this_.Token == token.Equal || this_.Token == token.NotEqual ||
		this_.Token == token.StrictEqual || this_.Token == token.StrictNotEqual {
		tkn := this_.Token
		this_.Next()
		left = &node.BinaryExpression{
			Operator:   tkn,
			Left:       left,
			Right:      this_.parseRelationalExpression(),
			Comparison: true,
		}
	}

	return left
}

func (this_ *Parser) parseBitwiseAndExpression() node.Expression {
	left := this_.parseEqualityExpression()

	for this_.Token == token.And {
		tkn := this_.Token
		this_.Next()
		left = &node.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    this_.parseEqualityExpression(),
		}
	}

	return left
}

func (this_ *Parser) parseBitwiseExclusiveOrExpression() node.Expression {
	left := this_.parseBitwiseAndExpression()

	for this_.Token == token.ExclusiveOr {
		tkn := this_.Token
		this_.Next()
		left = &node.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    this_.parseBitwiseAndExpression(),
		}
	}

	return left
}

func (this_ *Parser) parseBitwiseOrExpression() node.Expression {
	left := this_.parseBitwiseExclusiveOrExpression()

	for this_.Token == token.Or {
		tkn := this_.Token
		this_.Next()
		left = &node.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    this_.parseBitwiseExclusiveOrExpression(),
		}
	}

	return left
}

func (this_ *Parser) parseLogicalAndExpression() node.Expression {
	left := this_.parseBitwiseOrExpression()

	for this_.Token == token.LogicalAnd {
		tkn := this_.Token
		this_.Next()
		left = &node.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    this_.parseBitwiseOrExpression(),
		}
	}

	return left
}

func (this_ *Parser) parseLogicalOrExpression() node.Expression {
	var idx int
	parenthesis := this_.Token == token.LeftParenthesis
	left := this_.parseLogicalAndExpression()

	if this_.Token == token.LogicalOr || !parenthesis && this_.IsLogicalAndExpr(left) {
		for {
			switch this_.Token {
			case token.LogicalOr:
				this_.Next()
				left = &node.BinaryExpression{
					Operator: token.LogicalOr,
					Left:     left,
					Right:    this_.parseLogicalAndExpression(),
				}
			case token.Coalesce:
				idx = this_.Idx
				goto mixed
			default:
				return left
			}
		}
	} else {
		for {
			switch this_.Token {
			case token.Coalesce:
				idx = this_.Idx
				this_.Next()

				parenthesis := this_.Token == token.LeftParenthesis
				right := this_.parseLogicalAndExpression()
				if !parenthesis && this_.IsLogicalAndExpr(right) {
					goto mixed
				}

				left = &node.BinaryExpression{
					Operator: token.Coalesce,
					Left:     left,
					Right:    right,
				}
			case token.LogicalOr:
				idx = this_.Idx
				goto mixed
			default:
				return left
			}
		}
	}

mixed:
	_ = this_.Error("parseLogicalOrExpression", idx, "Logical expressions and coalesce expressions cannot be mixed. Wrap either by parentheses")
	return left
}

func (this_ *Parser) parseConditionalExpression() node.Expression {
	left := this_.parseLogicalOrExpression()

	if this_.Token == token.QuestionMark {
		this_.Next()
		allowIn := this_.Scope.AllowIn
		this_.Scope.AllowIn = true
		consequent := this_.parseAssignmentExpression()
		this_.Scope.AllowIn = allowIn
		this_.ExpectAndNext("parseConditionalExpression", token.Colon)
		return &node.ConditionalExpression{
			Test:       left,
			Consequent: consequent,
			Alternate:  this_.parseAssignmentExpression(),
		}
	}

	return left
}

func (this_ *Parser) parseArrowFunction(start int, paramList *node.ParameterList, async bool) node.Expression {
	this_.ExpectAndNext("parseArrowFunction", token.Arrow)
	res := &node.ArrowFunctionLiteral{
		Idx:           start,
		ParameterList: paramList,
		Async:         async,
	}
	res.Body, res.DeclarationList = this_.parseArrowFunctionBody(async)
	res.Source = this_.Slice(start, res.Body.End())
	return res
}

func (this_ *Parser) parseSingleArgArrowFunction(start int, async bool) node.Expression {
	if async != this_.Scope.AllowAwait {
		this_.Scope.AllowAwait = async
		defer func() {
			this_.Scope.AllowAwait = !async
		}()
	}
	this_.TokenToBindingIdentifier()
	if this_.Token != token.Identifier {
		_ = this_.ErrorUnexpectedToken("parseSingleArgArrowFunction this_.Token:"+this_.Token.String()+" not token.Identifier:"+token.Identifier.String(), this_.Token)
		this_.Next()
		return &node.BadExpression{
			From: start,
			To:   this_.Idx,
		}
	}

	id := this_.ParseIdentifier()

	paramList := &node.ParameterList{
		Opening: id.Idx,
		Closing: id.End(),
		List: []*node.Binding{{
			Target: id,
		}},
	}

	return this_.parseArrowFunction(start, paramList, async)
}

func (this_ *Parser) parseAssignmentExpression() node.Expression {
	start := this_.Idx
	parenthesis := false
	async := false
	var state parser.State
	switch this_.Token {
	case token.LeftParenthesis:
		this_.Mark(&state)
		parenthesis = true
	case token.Async:
		tok := this_.Peek()
		if this_.IsBindingIdentifier(tok) {
			// async x => ...
			this_.Next()
			return this_.parseSingleArgArrowFunction(start, true)
		} else if tok == token.LeftParenthesis {
			this_.Mark(&state)
			async = true
		}
	case token.Yield:
		if this_.Scope.AllowYield {
			return this_.parseYieldExpression()
		}
		fallthrough
	default:
		this_.TokenToBindingIdentifier()
	}
	left := this_.parseConditionalExpression()
	var operator token.Token
	switch this_.Token {
	case token.Assign:
		operator = this_.Token
	case token.AddAssign:
		operator = token.Plus
	case token.SubtractAssign:
		operator = token.Minus
	case token.MultiplyAssign:
		operator = token.Multiply
	case token.ExponentAssign:
		operator = token.Exponent
	case token.QuotientAssign:
		operator = token.Slash
	case token.RemainderAssign:
		operator = token.Remainder
	case token.AndAssign:
		operator = token.And
	case token.OrAssign:
		operator = token.Or
	case token.ExclusiveOrAssign:
		operator = token.ExclusiveOr
	case token.ShiftLeftAssign:
		operator = token.ShiftLeft
	case token.ShiftRightAssign:
		operator = token.ShiftRight
	case token.UnsignedShiftRightAssign:
		operator = token.UnsignedShiftRight
	case token.Arrow:
		var paramList *node.ParameterList
		if id, ok := left.(*node.Identifier); ok {
			paramList = &node.ParameterList{
				Opening: id.Idx,
				Closing: id.End() - 1,
				List: []*node.Binding{{
					Target: id,
				}},
			}
		} else if parenthesis {
			if seq, ok := left.(*node.SequenceExpression); ok && len(this_.Errors) == 0 {
				paramList = this_.reinterpretSequenceAsArrowFuncParams(seq.Sequence)
			} else {
				this_.Restore(&state)
				paramList = this_.parseFunctionParameterList()
			}
		} else if async {
			// async (x, y) => ...
			if !this_.Scope.AllowAwait {
				this_.Scope.AllowAwait = true
				defer func() {
					this_.Scope.AllowAwait = false
				}()
			}
			if _, ok := left.(*node.CallExpression); ok {
				this_.Restore(&state)
				this_.Next() // skip "async"
				paramList = this_.parseFunctionParameterList()
			}
		}
		if paramList == nil {
			_ = this_.Error("parseAssignmentExpression paramList is empty ", left.Start(), "Malformed arrow function parameter list")
			return &node.BadExpression{From: left.Start(), To: left.End()}
		}
		return this_.parseArrowFunction(start, paramList, async)
	}

	if operator != "" {
		idx := this_.Idx
		this_.Next()
		ok := false
		switch l := left.(type) {
		case *node.Identifier, *node.DotExpression, *node.PrivateDotExpression, *node.BracketExpression:
			ok = true
		case *node.ArrayLiteral:
			if !parenthesis && operator == token.Assign {
				left = this_.reinterpretAsArrayAssignmentPattern(l)
				ok = true
			}
		case *node.ObjectLiteral:
			if !parenthesis && operator == token.Assign {
				left = this_.reinterpretAsObjectAssignmentPattern(l)
				ok = true
			}
		}
		if ok {
			return &node.AssignExpression{
				Left:     left,
				Operator: operator,
				Right:    this_.parseAssignmentExpression(),
			}
		}
		_ = this_.Error("parseAssignmentExpression", left.Start(), "Invalid left-hand side in assignment")
		this_.nextStatement()
		return &node.BadExpression{From: idx, To: this_.Idx}
	}

	return left
}

func (this_ *Parser) parseYieldExpression() node.Expression {
	idx := this_.ExpectAndNext("parseYieldExpression", token.Yield)

	if this_.Scope.InFuncParams {
		_ = this_.Error("parseYieldExpression this_.scope.inFuncParams:true", idx, "Yield expression not allowed in formal parameter")
	}

	res := &node.YieldExpression{
		Yield: idx,
	}

	if !this_.ImplicitSemicolon && this_.Token == token.Multiply {
		res.Delegate = true
		this_.Next()
	}

	if !this_.ImplicitSemicolon && this_.Token != token.Semicolon && this_.Token != token.RightBrace && this_.Token != token.Eof {
		var state parser.State
		this_.Mark(&state)
		expr := this_.parseAssignmentExpression()
		if _, bad := expr.(*node.BadExpression); bad {
			expr = nil
			this_.Restore(&state)
		}
		res.Argument = expr
	}

	return res
}

func (this_ *Parser) parseExpression() node.Expression {
	left := this_.parseAssignmentExpression()

	if this_.Token == token.Comma {
		sequence := []node.Expression{left}
		for {
			if this_.Token != token.Comma {
				break
			}
			this_.Next()
			sequence = append(sequence, this_.parseAssignmentExpression())
		}
		return &node.SequenceExpression{
			Sequence: sequence,
		}
	}

	return left
}

func (this_ *Parser) checkComma(from, to int) {
	if pos := strings.IndexByte(this_.Str[(from)-this_.Base:(to)-this_.Base], ','); pos >= 0 {
		_ = this_.Error("checkComma", from+(pos), "Comma is not allowed here")
	}
}

func (this_ *Parser) reinterpretAsArrayAssignmentPattern(left *node.ArrayLiteral) node.Expression {
	value := left.Value
	var rest node.Expression
	for i, item := range value {
		if spread, ok := item.(*node.SpreadElement); ok {
			if i != len(value)-1 {
				_ = this_.Error("reinterpretAsArrayAssignmentPattern", item.Start(), "Rest element must be last element")
				return &node.BadExpression{From: left.Start(), To: left.End()}
			}
			this_.checkComma(spread.Expression.End(), left.RightBracket)
			rest = this_.reinterpretAsDestructAssignTarget(spread.Expression)
			value = value[:len(value)-1]
		} else {
			value[i] = this_.reinterpretAsAssignmentElement(item)
		}
	}
	return &node.ArrayPattern{
		LeftBracket:  left.LeftBracket,
		RightBracket: left.RightBracket,
		Elements:     value,
		Rest:         rest,
	}
}

func (this_ *Parser) reinterpretArrayAssignPatternAsBinding(pattern *node.ArrayPattern) *node.ArrayPattern {
	for i, item := range pattern.Elements {
		pattern.Elements[i] = this_.reinterpretAsDestructBindingTarget(item)
	}
	if pattern.Rest != nil {
		pattern.Rest = this_.reinterpretAsDestructBindingTarget(pattern.Rest)
	}
	return pattern
}

func (this_ *Parser) reinterpretAsArrayBindingPattern(left *node.ArrayLiteral) node.BindingTarget {
	value := left.Value
	var rest node.Expression
	for i, item := range value {
		if spread, ok := item.(*node.SpreadElement); ok {
			if i != len(value)-1 {
				_ = this_.Error("reinterpretAsArrayBindingPattern", item.Start(), "Rest element must be last element")
				return &node.BadExpression{From: left.Start(), To: left.End()}
			}
			this_.checkComma(spread.Expression.End(), left.RightBracket)
			rest = this_.reinterpretAsDestructBindingTarget(spread.Expression)
			value = value[:len(value)-1]
		} else {
			value[i] = this_.reinterpretAsBindingElement(item)
		}
	}
	return &node.ArrayPattern{
		LeftBracket:  left.LeftBracket,
		RightBracket: left.RightBracket,
		Elements:     value,
		Rest:         rest,
	}
}

func (this_ *Parser) parseArrayBindingPattern() node.BindingTarget {
	return this_.reinterpretAsArrayBindingPattern(this_.parseArrayLiteral())
}

func (this_ *Parser) parseObjectBindingPattern() node.BindingTarget {
	return this_.reinterpretAsObjectBindingPattern(this_.parseObjectLiteral())
}

func (this_ *Parser) reinterpretArrayObjectPatternAsBinding(pattern *node.ObjectPattern) *node.ObjectPattern {
	for _, prop := range pattern.Properties {
		if keyed, ok := prop.(*node.PropertyKeyed); ok {
			keyed.Value = this_.reinterpretAsBindingElement(keyed.Value)
		}
	}
	if pattern.Rest != nil {
		pattern.Rest = this_.reinterpretAsBindingRestElement(pattern.Rest)
	}
	return pattern
}

func (this_ *Parser) reinterpretAsObjectBindingPattern(expr *node.ObjectLiteral) node.BindingTarget {
	var rest node.Expression
	value := expr.Value
	for i, prop := range value {
		ok := false
		switch prop := prop.(type) {
		case *node.PropertyKeyed:
			if prop.Kind == node.PropertyKindValue {
				prop.Value = this_.reinterpretAsBindingElement(prop.Value)
				ok = true
			}
		case *node.PropertyShort:
			ok = true
		case *node.SpreadElement:
			if i != len(expr.Value)-1 {
				_ = this_.Error("reinterpretAsObjectBindingPattern", prop.Start(), "Rest element must be last element")
				return &node.BadExpression{From: expr.Start(), To: expr.End()}
			}
			// TODO make sure there is no trailing Comma
			rest = this_.reinterpretAsBindingRestElement(prop.Expression)
			value = value[:i]
			ok = true
		}
		if !ok {
			_ = this_.Error("reinterpretAsObjectBindingPattern", prop.Start(), "Invalid destructuring binding target")
			return &node.BadExpression{From: expr.Start(), To: expr.End()}
		}
	}
	return &node.ObjectPattern{
		LeftBrace:  expr.LeftBrace,
		RightBrace: expr.RightBrace,
		Properties: value,
		Rest:       rest,
	}
}

func (this_ *Parser) reinterpretAsObjectAssignmentPattern(l *node.ObjectLiteral) node.Expression {
	var rest node.Expression
	value := l.Value
	for i, prop := range value {
		ok := false
		switch prop := prop.(type) {
		case *node.PropertyKeyed:
			if prop.Kind == node.PropertyKindValue {
				prop.Value = this_.reinterpretAsAssignmentElement(prop.Value)
				ok = true
			}
		case *node.PropertyShort:
			ok = true
		case *node.SpreadElement:
			if i != len(l.Value)-1 {
				_ = this_.Error("reinterpretAsObjectAssignmentPattern", prop.Start(), "Rest element must be last element")
				return &node.BadExpression{From: l.Start(), To: l.End()}
			}
			// TODO make sure there is no trailing Comma
			rest = prop.Expression
			value = value[:i]
			ok = true
		}
		if !ok {
			_ = this_.Error("reinterpretAsObjectAssignmentPattern", prop.Start(), "Invalid destructuring assignment target")
			return &node.BadExpression{From: l.Start(), To: l.End()}
		}
	}
	return &node.ObjectPattern{
		LeftBrace:  l.LeftBrace,
		RightBrace: l.RightBrace,
		Properties: value,
		Rest:       rest,
	}
}

func (this_ *Parser) reinterpretAsAssignmentElement(expr node.Expression) node.Expression {
	switch expr := expr.(type) {
	case *node.AssignExpression:
		if expr.Operator == token.Assign {
			expr.Left = this_.reinterpretAsDestructAssignTarget(expr.Left)
			return expr
		} else {
			_ = this_.Error("reinterpretAsAssignmentElement", expr.Start(), "Invalid destructuring assignment target")
			return &node.BadExpression{From: expr.Start(), To: expr.End()}
		}
	default:
		return this_.reinterpretAsDestructAssignTarget(expr)
	}
}

func (this_ *Parser) reinterpretAsBindingElement(expr node.Expression) node.Expression {
	switch expr := expr.(type) {
	case *node.AssignExpression:
		if expr.Operator == token.Assign {
			expr.Left = this_.reinterpretAsDestructBindingTarget(expr.Left)
			return expr
		} else {
			_ = this_.Error("reinterpretAsBindingElement", expr.Start(), "Invalid destructuring assignment target")
			return &node.BadExpression{From: expr.Start(), To: expr.End()}
		}
	default:
		return this_.reinterpretAsDestructBindingTarget(expr)
	}
}

func (this_ *Parser) reinterpretAsBinding(expr node.Expression) *node.Binding {
	switch expr := expr.(type) {
	case *node.AssignExpression:
		if expr.Operator == token.Assign {
			return &node.Binding{
				Target:      this_.reinterpretAsDestructBindingTarget(expr.Left),
				Initializer: expr.Right,
			}
		} else {
			_ = this_.Error("reinterpretAsBinding", expr.Start(), "Invalid destructuring assignment target")
			return &node.Binding{
				Target: &node.BadExpression{From: expr.Start(), To: expr.End()},
			}
		}
	default:
		return &node.Binding{
			Target: this_.reinterpretAsDestructBindingTarget(expr),
		}
	}
}

func (this_ *Parser) reinterpretAsDestructAssignTarget(item node.Expression) node.Expression {
	switch item := item.(type) {
	case nil:
		return nil
	case *node.ArrayLiteral:
		return this_.reinterpretAsArrayAssignmentPattern(item)
	case *node.ObjectLiteral:
		return this_.reinterpretAsObjectAssignmentPattern(item)
	case node.Pattern, *node.Identifier, *node.DotExpression, *node.PrivateDotExpression, *node.BracketExpression:
		return item
	}
	_ = this_.Error("reinterpretAsDestructAssignTarget", item.Start(), "Invalid destructuring assignment target")
	return &node.BadExpression{From: item.Start(), To: item.End()}
}

func (this_ *Parser) reinterpretAsDestructBindingTarget(item node.Expression) node.BindingTarget {
	switch item := item.(type) {
	case nil:
		return nil
	case *node.ArrayPattern:
		return this_.reinterpretArrayAssignPatternAsBinding(item)
	case *node.ObjectPattern:
		return this_.reinterpretArrayObjectPatternAsBinding(item)
	case *node.ArrayLiteral:
		return this_.reinterpretAsArrayBindingPattern(item)
	case *node.ObjectLiteral:
		return this_.reinterpretAsObjectBindingPattern(item)
	case *node.Identifier:
		if !this_.Scope.AllowAwait || item.Name != "await" {
			return item
		}
	}
	_ = this_.Error("reinterpretAsDestructBindingTarget", item.Start(), "Invalid destructuring binding target")
	return &node.BadExpression{From: item.Start(), To: item.End()}
}

func (this_ *Parser) reinterpretAsBindingRestElement(expr node.Expression) node.Expression {
	if _, ok := expr.(*node.Identifier); ok {
		return expr
	}
	_ = this_.Error("reinterpretAsBindingRestElement", expr.Start(), "Invalid binding rest")
	return &node.BadExpression{From: expr.Start(), To: expr.End()}
}
