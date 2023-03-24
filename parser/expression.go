package parser

import (
	"fmt"
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/token"
	"reflect"
	"strings"
)

func (this_ *parser) parseIdentifier() *node.Identifier {
	literal := this_.parsedLiteral
	idx := this_.idx
	this_.next()
	return &node.Identifier{
		Name: literal,
		Idx:  idx,
	}
}

func (this_ *parser) parsePrimaryExpression() node.Expression {
	literal, parsedLiteral := this_.literal, this_.parsedLiteral
	idx := this_.idx
	switch this_.token {
	case token.Identifier:
		this_.next()
		return &node.Identifier{
			Name: parsedLiteral,
			Idx:  idx,
		}
	case token.Null:
		this_.next()
		return &node.NullLiteral{
			Idx:     idx,
			Literal: literal,
		}
	case token.Boolean:
		this_.next()
		value := false
		switch parsedLiteral {
		case "true":
			value = true
		case "false":
			value = false
		default:
			_ = this_.error("parsePrimaryExpression parsedLiteral:"+string(parsedLiteral), idx, "Illegal boolean literal")
		}
		return &node.BooleanLiteral{
			Idx:     idx,
			Literal: literal,
			Value:   value,
		}
	case token.String:
		this_.next()
		return &node.StringLiteral{
			Idx:     idx,
			Literal: literal,
			Value:   parsedLiteral,
		}
	case token.Number:
		this_.next()
		value, err := this_.parseNumberLiteral(literal)
		if err != nil {
			_ = this_.error("parsePrimaryExpression parseNumberLiteral error literal:"+string(literal), idx, err.Error())
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
		this_.next()
		return &node.ThisExpression{
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

	if this_.isBindingId(this_.token) {
		this_.next()
		return &node.Identifier{
			Name: parsedLiteral,
			Idx:  idx,
		}
	}

	_ = this_.errorUnexpectedToken("parsePrimaryExpression", this_.token)
	this_.nextStatement()
	return &node.BadExpression{From: idx, To: this_.idx}
}

func (this_ *parser) parseSuperProperty() node.Expression {
	idx := this_.idx
	this_.next()
	switch this_.token {
	case token.Period:
		this_.next()
		if !this_.IsIdentifierToken(this_.token) {
			this_.expect(token.Identifier)
			this_.nextStatement()
			return &node.BadExpression{From: idx, To: this_.idx}
		}
		idIdx := this_.idx
		parsedLiteral := this_.parsedLiteral
		this_.next()
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
		_ = this_.error("parseSuperProperty this_.token:"+this_.token.String(), idx, "'super' keyword unexpected here")
		this_.nextStatement()
		return &node.BadExpression{From: idx, To: this_.idx}
	}
}

func (this_ *parser) reinterpretSequenceAsArrowFuncParams(list []node.Expression) *node.ParameterList {
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
			_ = this_.error("reinterpretSequenceAsArrowFuncParams firstRestIdx != -1 firstRestIdx:"+fmt.Sprintf("%d", firstRestIdx), list[firstRestIdx].Start(), "Rest parameter must be last formal parameter")
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

func (this_ *parser) parseParenthesisedExpression() node.Expression {
	opening := this_.idx
	this_.expect(token.LeftParenthesis)
	var list []node.Expression
	if this_.token != token.RightParenthesis {
		for {
			if this_.token == token.Ellipsis {
				start := this_.idx
				_ = this_.errorUnexpectedToken("parseParenthesisedExpression", token.Ellipsis)
				this_.next()
				expr := this_.parseAssignmentExpression()
				list = append(list, &node.BadExpression{
					From: start,
					To:   expr.End(),
				})
			} else {
				list = append(list, this_.parseAssignmentExpression())
			}
			if this_.token != token.Comma {
				break
			}
			this_.next()
			if this_.token == token.RightParenthesis {
				_ = this_.errorUnexpectedToken("parseParenthesisedExpression", token.RightParenthesis)
				break
			}
		}
	}
	this_.expect(token.RightParenthesis)
	if len(list) == 1 && len(this_.errors) == 0 {
		return list[0]
	}
	if len(list) == 0 {
		_ = this_.errorUnexpectedToken("parseParenthesisedExpression", token.RightParenthesis)
		return &node.BadExpression{
			From: opening,
			To:   this_.idx,
		}
	}
	return &node.SequenceExpression{
		Sequence: list,
	}
}

func (this_ *parser) parseRegExpLiteral() *node.RegExpLiteral {

	offset := this_.chrOffset - 1 // Opening slash already gotten
	if this_.token == token.QuotientAssign {
		offset -= 1 // =
	}
	idx := offset

	pattern, _, err := this_.scanString(offset, false)
	endOffset := this_.chrOffset

	if err == "" {
		pattern = pattern[1 : len(pattern)-1]
	}

	flags := ""
	if !this_.IsLineTerminator(this_.chr) && !this_.IsLineWhiteSpace(this_.chr) {
		this_.next()

		if this_.token == token.Identifier { // gim

			flags = this_.literal
			this_.next()
			endOffset = this_.chrOffset - 1
		}
	} else {
		this_.next()
	}

	literal := this_.str[offset:endOffset]

	return &node.RegExpLiteral{
		Idx:     idx,
		Literal: literal,
		Pattern: pattern,
		Flags:   flags,
	}
}

func (this_ *parser) tokenToBindingId() {
	if this_.isBindingId(this_.token) {
		this_.token = token.Identifier
	}
}

func (this_ *parser) parseBindingTarget() (target node.BindingTarget) {
	this_.tokenToBindingId()
	switch this_.token {
	case token.Identifier:
		target = &node.Identifier{
			Name: this_.parsedLiteral,
			Idx:  this_.idx,
		}
		this_.next()
	case token.LeftBracket:
		target = this_.parseArrayBindingPattern()
	case token.LeftBrace:
		target = this_.parseObjectBindingPattern()
	default:
		idx := this_.expect(token.Identifier)
		this_.nextStatement()
		target = &node.BadExpression{From: idx, To: this_.idx}
	}

	return
}

func (this_ *parser) parseVariableDeclaration(declarationList *[]*node.Binding) *node.Binding {
	res := &node.Binding{
		Target: this_.parseBindingTarget(),
	}

	if declarationList != nil {
		*declarationList = append(*declarationList, res)
	}

	if this_.token == token.Assign {
		this_.next()
		res.Initializer = this_.parseAssignmentExpression()
	}

	return res
}

func (this_ *parser) parseVariableDeclarationList() (declarationList []*node.Binding) {
	for {
		this_.parseVariableDeclaration(&declarationList)
		if this_.token != token.Comma {
			break
		}
		this_.next()
	}
	return
}

func (this_ *parser) parseVarDeclarationList(var_ int) []*node.Binding {
	declarationList := this_.parseVariableDeclarationList()

	this_.scope.declare(&node.VariableDeclaration{
		Var:  var_,
		List: declarationList,
	})

	return declarationList
}

func (this_ *parser) parseObjectPropertyKey() (string, node.String, node.Expression, token.Token) {
	if this_.token == token.LeftBracket {
		this_.next()
		expr := this_.parseAssignmentExpression()
		this_.expect(token.RightBracket)
		return "", "", expr, token.Illegal
	}
	idx, tkn, literal, parsedLiteral := this_.idx, this_.token, this_.literal, this_.parsedLiteral
	var value node.Expression
	this_.next()
	switch tkn {
	case token.Identifier, token.String, token.Keyword, token.EscapedReservedWord:
		value = &node.StringLiteral{
			Idx:     idx,
			Literal: literal,
			Value:   parsedLiteral,
		}
	case token.Number:
		num, err := this_.parseNumberLiteral(literal)
		if err != nil {
			_ = this_.error("parseObjectPropertyKey parseNumberLiteral literal:"+string(literal), idx, err.Error())
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
				Value:   node.String(literal),
			}
		} else {
			_ = this_.errorUnexpectedToken("parseObjectPropertyKey not IsIdentifierToken:"+tkn.String(), tkn)
		}
	}
	return literal, parsedLiteral, value, tkn
}

func (this_ *parser) parseObjectProperty() node.Property {
	if this_.token == token.Ellipsis {
		this_.next()
		return &node.SpreadElement{
			Expression: this_.parseAssignmentExpression(),
		}
	}
	keyStartIdx := this_.idx
	generator := false
	if this_.token == token.Multiply {
		generator = true
		this_.next()
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
		case this_.token == token.LeftParenthesis:
			return &node.PropertyKeyed{
				Key:      value,
				Kind:     node.PropertyKindMethod,
				Value:    this_.parseMethodDefinition(keyStartIdx, node.PropertyKindMethod, false, false),
				Computed: tkn == token.Illegal,
			}
		case this_.token == token.Comma || this_.token == token.RightBrace || this_.token == token.Assign: // shorthand property
			if this_.isBindingId(tkn) {
				var initializer node.Expression
				if this_.token == token.Assign {
					// allow the initializer syntax here in case the object literal
					// needs to be reinterpreted as an assignment pattern, enforce later if it doesn't.
					this_.next()
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
				_ = this_.errorUnexpectedToken("parseObjectProperty not this_.isBindingId:"+tkn.String(), this_.token)
			}
		case (literal == "get" || literal == "set" || tkn == token.Async) && this_.token != token.Colon:
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

	this_.expect(token.Colon)
	return &node.PropertyKeyed{
		Key:      value,
		Kind:     node.PropertyKindValue,
		Value:    this_.parseAssignmentExpression(),
		Computed: tkn == token.Illegal,
	}
}

func (this_ *parser) parseMethodDefinition(keyStartIdx int, kind node.PropertyKind, generator, async bool) *node.FunctionLiteral {
	idx1 := this_.idx
	if generator != this_.scope.allowYield {
		this_.scope.allowYield = generator
		defer func() {
			this_.scope.allowYield = !generator
		}()
	}
	if async != this_.scope.allowAwait {
		this_.scope.allowAwait = async
		defer func() {
			this_.scope.allowAwait = !async
		}()
	}
	parameterList := this_.parseFunctionParameterList()
	switch kind {
	case node.PropertyKindGet:
		if len(parameterList.List) > 0 || parameterList.Rest != nil {
			_ = this_.error("parseMethodDefinition node.PropertyKindGet", idx1, "Getter must not have any formal parameters.")
		}
	case node.PropertyKindSet:
		if len(parameterList.List) != 1 || parameterList.Rest != nil {
			_ = this_.error("parseMethodDefinition node.PropertyKindSet", idx1, "Setter must have exactly one formal parameter.")
		}
	}
	res := &node.FunctionLiteral{
		Function:      keyStartIdx,
		ParameterList: parameterList,
		Generator:     generator,
		Async:         async,
	}
	res.Body, res.DeclarationList = this_.parseFunctionBlock(async, async, generator)
	res.Source = this_.slice(keyStartIdx, res.Body.End())
	return res
}

func (this_ *parser) parseObjectLiteral() *node.ObjectLiteral {
	var value []node.Property
	idx0 := this_.expect(token.LeftBrace)
	for this_.token != token.RightBrace && this_.token != token.Eof {
		property := this_.parseObjectProperty()
		if property != nil {
			value = append(value, property)
		}
		if this_.token != token.RightBrace {
			this_.expect(token.Comma)
		} else {
			break
		}
	}
	idx1 := this_.expect(token.RightBrace)

	return &node.ObjectLiteral{
		LeftBrace:  idx0,
		RightBrace: idx1,
		Value:      value,
	}
}

func (this_ *parser) parseArrayLiteral() *node.ArrayLiteral {

	idx0 := this_.expect(token.LeftBracket)
	var value []node.Expression
	for this_.token != token.RightBracket && this_.token != token.Eof {
		if this_.token == token.Comma {
			this_.next()
			value = append(value, nil)
			continue
		}
		if this_.token == token.Ellipsis {
			this_.next()
			value = append(value, &node.SpreadElement{
				Expression: this_.parseAssignmentExpression(),
			})
		} else {
			value = append(value, this_.parseAssignmentExpression())
		}
		if this_.token != token.RightBracket {
			this_.expect(token.Comma)
		}
	}
	idx1 := this_.expect(token.RightBracket)

	return &node.ArrayLiteral{
		LeftBracket:  idx0,
		RightBracket: idx1,
		Value:        value,
	}
}

func (this_ *parser) parseTemplateLiteral(tagged bool) *node.TemplateLiteral {
	res := &node.TemplateLiteral{
		OpenQuote: this_.idx,
	}
	for {
		start := this_.offset
		literal, parsed, finished, parseErr, err := this_.parseTemplateCharacters()
		if err != "" {
			_ = this_.error("parseTemplateLiteral parseTemplateCharacters err", this_.offset, err)
		}
		res.Elements = append(res.Elements, &node.TemplateElement{
			Idx:     start,
			Literal: literal,
			Parsed:  parsed,
			Valid:   parseErr == "",
		})
		if !tagged && parseErr != "" {
			_ = this_.error("parseTemplateLiteral parseTemplateCharacters parseErr", this_.offset, parseErr)
		}
		end := this_.chrOffset - 1
		this_.next()
		if finished {
			res.CloseQuote = end
			break
		}
		expr := this_.parseExpression()
		res.Expressions = append(res.Expressions, expr)
		if this_.token != token.RightBrace {
			_ = this_.errorUnexpectedToken("parseTemplateLiteral this_.token:"+this_.token.String()+" is not token.RightBrace:"+token.RightBrace.String(), this_.token)
		}
	}
	return res
}

func (this_ *parser) parseTaggedTemplateLiteral(tag node.Expression) *node.TemplateLiteral {
	l := this_.parseTemplateLiteral(true)
	l.Tag = tag
	return l
}

func (this_ *parser) parseArgumentList() (argumentList []node.Expression, idx0, idx1 int) {
	idx0 = this_.expect(token.LeftParenthesis)
	for this_.token != token.RightParenthesis {
		var item node.Expression
		if this_.token == token.Ellipsis {
			this_.next()
			item = &node.SpreadElement{
				Expression: this_.parseAssignmentExpression(),
			}
		} else {
			item = this_.parseAssignmentExpression()
		}
		argumentList = append(argumentList, item)
		if this_.token != token.Comma {
			break
		}
		this_.next()
	}
	idx1 = this_.expect(token.RightParenthesis)
	return
}

func (this_ *parser) parseCallExpression(left node.Expression) node.Expression {
	argumentList, idx0, idx1 := this_.parseArgumentList()
	return &node.CallExpression{
		Callee:           left,
		LeftParenthesis:  idx0,
		ArgumentList:     argumentList,
		RightParenthesis: idx1,
	}
}

func (this_ *parser) parseDotMember(left node.Expression) node.Expression {
	period := this_.idx
	this_.next()

	literal := this_.parsedLiteral
	idx := this_.idx

	if this_.token == token.PrivateIdentifier {
		this_.next()
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

	if !this_.IsIdentifierToken(this_.token) {
		this_.expect(token.Identifier)
		this_.nextStatement()
		return &node.BadExpression{From: period, To: this_.idx}
	}

	this_.next()

	return &node.DotExpression{
		Left: left,
		Identifier: node.Identifier{
			Idx:  idx,
			Name: literal,
		},
	}
}

func (this_ *parser) parseBracketMember(left node.Expression) node.Expression {
	idx0 := this_.expect(token.LeftBracket)
	member := this_.parseExpression()
	idx1 := this_.expect(token.RightBracket)
	return &node.BracketExpression{
		LeftBracket:  idx0,
		Left:         left,
		Member:       member,
		RightBracket: idx1,
	}
}

func (this_ *parser) parseNewExpression() node.Expression {
	idx := this_.expect(token.New)
	if this_.token == token.Period {
		this_.next()
		if this_.literal == "target" {
			return &node.MetaProperty{
				Meta: &node.Identifier{
					Name: node.String(token.New.String()),
					Idx:  idx,
				},
				Property: this_.parseIdentifier(),
			}
		}
		_ = this_.errorUnexpectedToken("parseNewExpression", token.Identifier)
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
	if this_.token == token.LeftParenthesis {
		argumentList, idx0, idx1 := this_.parseArgumentList()
		res.ArgumentList = argumentList
		res.LeftParenthesis = idx0
		res.RightParenthesis = idx1
	}
	return res
}

func (this_ *parser) parseLeftHandSideExpression() node.Expression {

	var left node.Expression
	if this_.token == token.New {
		left = this_.parseNewExpression()
	} else {
		left = this_.parsePrimaryExpression()
	}
L:
	for {
		switch this_.token {
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

func (this_ *parser) parseLeftHandSideExpressionAllowCall() node.Expression {

	allowIn := this_.scope.allowIn
	this_.scope.allowIn = true
	defer func() {
		this_.scope.allowIn = allowIn
	}()

	var left node.Expression
	start := this_.idx
	if this_.token == token.New {
		left = this_.parseNewExpression()
	} else {
		left = this_.parsePrimaryExpression()
	}

	optionalChain := false
L:
	for {
		switch this_.token {
		case token.Period:
			left = this_.parseDotMember(left)
		case token.LeftBracket:
			left = this_.parseBracketMember(left)
		case token.LeftParenthesis:
			left = this_.parseCallExpression(left)
		case token.Backtick:
			if optionalChain {
				_ = this_.error("parseLeftHandSideExpressionAllowCall token.Backtick optionalChain:true", this_.idx, "Invalid template literal on optional chain")
				this_.nextStatement()
				return &node.BadExpression{From: start, To: this_.idx}
			}
			left = this_.parseTaggedTemplateLiteral(left)
		case token.QuestionDot:
			optionalChain = true
			left = &node.Optional{Expression: left}

			switch this_.peek() {
			case token.LeftBracket, token.LeftParenthesis, token.Backtick:
				this_.next()
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

func (this_ *parser) parsePostfixExpression() node.Expression {
	operand := this_.parseLeftHandSideExpressionAllowCall()

	switch this_.token {
	case token.Increment, token.Decrement:
		// Make sure there is no line terminator here
		if this_.implicitSemicolon {
			break
		}
		tkn := this_.token
		idx := this_.idx
		this_.next()
		switch operand.(type) {
		case *node.Identifier, *node.DotExpression, *node.PrivateDotExpression, *node.BracketExpression:
		default:
			_ = this_.error("parsePostfixExpression operand type:"+reflect.TypeOf(operand).String(), idx, "Invalid left-hand side in assignment")
			this_.nextStatement()
			return &node.BadExpression{From: idx, To: this_.idx}
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

func (this_ *parser) parseUnaryExpression() node.Expression {

	switch this_.token {
	case token.Plus, token.Minus, token.Not, token.BitwiseNot:
		fallthrough
	case token.Delete, token.Void, token.Typeof:
		tkn := this_.token
		idx := this_.idx
		this_.next()
		return &node.UnaryExpression{
			Operator: tkn,
			Idx:      idx,
			Operand:  this_.parseUnaryExpression(),
		}
	case token.Increment, token.Decrement:
		tkn := this_.token
		idx := this_.idx
		this_.next()
		operand := this_.parseUnaryExpression()
		switch operand.(type) {
		case *node.Identifier, *node.DotExpression, *node.PrivateDotExpression, *node.BracketExpression:
		default:
			_ = this_.error("parseUnaryExpression operand type:"+reflect.TypeOf(operand).String(), idx, "Invalid left-hand side in assignment")
			this_.nextStatement()
			return &node.BadExpression{From: idx, To: this_.idx}
		}
		return &node.UnaryExpression{
			Operator: tkn,
			Idx:      idx,
			Operand:  operand,
		}
	case token.Await:
		if this_.scope.allowAwait {
			idx := this_.idx
			this_.next()
			if !this_.scope.inAsync {
				_ = this_.errorUnexpectedToken("parseUnaryExpression is not this_.scope.inAsync", token.Await)
				return &node.BadExpression{
					From: idx,
					To:   this_.idx,
				}
			}
			if this_.scope.inFuncParams {
				_ = this_.error("parseUnaryExpression this_.scope.inFuncParams", idx, "Illegal await-expression in formal parameters of async function")
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

func (this_ *parser) parseExponentiationExpression() node.Expression {
	left := this_.parseUnaryExpression()

	for this_.token == token.Exponent && isUpdateExpression(left) {
		this_.next()
		left = &node.BinaryExpression{
			Operator: token.Exponent,
			Left:     left,
			Right:    this_.parseExponentiationExpression(),
		}
	}

	return left
}

func (this_ *parser) parseMultiplicativeExpression() node.Expression {
	left := this_.parseExponentiationExpression()

	for this_.token == token.Multiply || this_.token == token.Slash ||
		this_.token == token.Remainder {
		tkn := this_.token
		this_.next()
		left = &node.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    this_.parseExponentiationExpression(),
		}
	}

	return left
}

func (this_ *parser) parseAdditiveExpression() node.Expression {
	left := this_.parseMultiplicativeExpression()

	for this_.token == token.Plus || this_.token == token.Minus {
		tkn := this_.token
		this_.next()
		left = &node.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    this_.parseMultiplicativeExpression(),
		}
	}

	return left
}

func (this_ *parser) parseShiftExpression() node.Expression {
	left := this_.parseAdditiveExpression()

	for this_.token == token.ShiftLeft || this_.token == token.ShiftRight ||
		this_.token == token.UnsignedShiftRight {
		tkn := this_.token
		this_.next()
		left = &node.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    this_.parseAdditiveExpression(),
		}
	}

	return left
}

func (this_ *parser) parseRelationalExpression() node.Expression {
	if this_.scope.allowIn && this_.token == token.PrivateIdentifier {
		left := &node.PrivateIdentifier{
			Identifier: node.Identifier{
				Idx:  this_.idx,
				Name: this_.parsedLiteral,
			},
		}
		this_.next()
		if this_.token == token.In {
			this_.next()
			return &node.BinaryExpression{
				Operator: this_.token,
				Left:     left,
				Right:    this_.parseShiftExpression(),
			}
		}
		return left
	}
	left := this_.parseShiftExpression()

	allowIn := this_.scope.allowIn
	this_.scope.allowIn = true
	defer func() {
		this_.scope.allowIn = allowIn
	}()

	switch this_.token {
	case token.Less, token.LessOrEqual, token.Greater, token.GreaterOrEqual:
		tkn := this_.token
		this_.next()
		return &node.BinaryExpression{
			Operator:   tkn,
			Left:       left,
			Right:      this_.parseRelationalExpression(),
			Comparison: true,
		}
	case token.Instanceof:
		tkn := this_.token
		this_.next()
		return &node.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    this_.parseRelationalExpression(),
		}
	case token.In:
		if !allowIn {
			return left
		}
		tkn := this_.token
		this_.next()
		return &node.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    this_.parseRelationalExpression(),
		}
	}

	return left
}

func (this_ *parser) parseEqualityExpression() node.Expression {
	left := this_.parseRelationalExpression()

	for this_.token == token.Equal || this_.token == token.NotEqual ||
		this_.token == token.StrictEqual || this_.token == token.StrictNotEqual {
		tkn := this_.token
		this_.next()
		left = &node.BinaryExpression{
			Operator:   tkn,
			Left:       left,
			Right:      this_.parseRelationalExpression(),
			Comparison: true,
		}
	}

	return left
}

func (this_ *parser) parseBitwiseAndExpression() node.Expression {
	left := this_.parseEqualityExpression()

	for this_.token == token.And {
		tkn := this_.token
		this_.next()
		left = &node.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    this_.parseEqualityExpression(),
		}
	}

	return left
}

func (this_ *parser) parseBitwiseExclusiveOrExpression() node.Expression {
	left := this_.parseBitwiseAndExpression()

	for this_.token == token.ExclusiveOr {
		tkn := this_.token
		this_.next()
		left = &node.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    this_.parseBitwiseAndExpression(),
		}
	}

	return left
}

func (this_ *parser) parseBitwiseOrExpression() node.Expression {
	left := this_.parseBitwiseExclusiveOrExpression()

	for this_.token == token.Or {
		tkn := this_.token
		this_.next()
		left = &node.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    this_.parseBitwiseExclusiveOrExpression(),
		}
	}

	return left
}

func (this_ *parser) parseLogicalAndExpression() node.Expression {
	left := this_.parseBitwiseOrExpression()

	for this_.token == token.LogicalAnd {
		tkn := this_.token
		this_.next()
		left = &node.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    this_.parseBitwiseOrExpression(),
		}
	}

	return left
}

func isLogicalAndExpr(expr node.Expression) bool {
	if exp, ok := expr.(*node.BinaryExpression); ok && exp.Operator == token.LogicalAnd {
		return true
	}
	return false
}

func (this_ *parser) parseLogicalOrExpression() node.Expression {
	var idx int
	parenthesis := this_.token == token.LeftParenthesis
	left := this_.parseLogicalAndExpression()

	if this_.token == token.LogicalOr || !parenthesis && isLogicalAndExpr(left) {
		for {
			switch this_.token {
			case token.LogicalOr:
				this_.next()
				left = &node.BinaryExpression{
					Operator: token.LogicalOr,
					Left:     left,
					Right:    this_.parseLogicalAndExpression(),
				}
			case token.Coalesce:
				idx = this_.idx
				goto mixed
			default:
				return left
			}
		}
	} else {
		for {
			switch this_.token {
			case token.Coalesce:
				idx = this_.idx
				this_.next()

				parenthesis := this_.token == token.LeftParenthesis
				right := this_.parseLogicalAndExpression()
				if !parenthesis && isLogicalAndExpr(right) {
					goto mixed
				}

				left = &node.BinaryExpression{
					Operator: token.Coalesce,
					Left:     left,
					Right:    right,
				}
			case token.LogicalOr:
				idx = this_.idx
				goto mixed
			default:
				return left
			}
		}
	}

mixed:
	_ = this_.error("parseLogicalOrExpression", idx, "Logical expressions and coalesce expressions cannot be mixed. Wrap either by parentheses")
	return left
}

func (this_ *parser) parseConditionalExpression() node.Expression {
	left := this_.parseLogicalOrExpression()

	if this_.token == token.QuestionMark {
		this_.next()
		allowIn := this_.scope.allowIn
		this_.scope.allowIn = true
		consequent := this_.parseAssignmentExpression()
		this_.scope.allowIn = allowIn
		this_.expect(token.Colon)
		return &node.ConditionalExpression{
			Test:       left,
			Consequent: consequent,
			Alternate:  this_.parseAssignmentExpression(),
		}
	}

	return left
}

func (this_ *parser) parseArrowFunction(start int, paramList *node.ParameterList, async bool) node.Expression {
	this_.expect(token.Arrow)
	res := &node.ArrowFunctionLiteral{
		Start_:        start,
		ParameterList: paramList,
		Async:         async,
	}
	res.Body, res.DeclarationList = this_.parseArrowFunctionBody(async)
	res.Source = this_.slice(start, res.Body.End())
	return res
}

func (this_ *parser) parseSingleArgArrowFunction(start int, async bool) node.Expression {
	if async != this_.scope.allowAwait {
		this_.scope.allowAwait = async
		defer func() {
			this_.scope.allowAwait = !async
		}()
	}
	this_.tokenToBindingId()
	if this_.token != token.Identifier {
		_ = this_.errorUnexpectedToken("parseSingleArgArrowFunction this_.token:"+this_.token.String()+" not token.Identifier:"+token.Identifier.String(), this_.token)
		this_.next()
		return &node.BadExpression{
			From: start,
			To:   this_.idx,
		}
	}

	id := this_.parseIdentifier()

	paramList := &node.ParameterList{
		Opening: id.Idx,
		Closing: id.End(),
		List: []*node.Binding{{
			Target: id,
		}},
	}

	return this_.parseArrowFunction(start, paramList, async)
}

func (this_ *parser) parseAssignmentExpression() node.Expression {
	start := this_.idx
	parenthesis := false
	async := false
	var state parserState
	switch this_.token {
	case token.LeftParenthesis:
		this_.mark(&state)
		parenthesis = true
	case token.Async:
		tok := this_.peek()
		if this_.isBindingId(tok) {
			// async x => ...
			this_.next()
			return this_.parseSingleArgArrowFunction(start, true)
		} else if tok == token.LeftParenthesis {
			this_.mark(&state)
			async = true
		}
	case token.Yield:
		if this_.scope.allowYield {
			return this_.parseYieldExpression()
		}
		fallthrough
	default:
		this_.tokenToBindingId()
	}
	left := this_.parseConditionalExpression()
	var operator token.Token
	switch this_.token {
	case token.Assign:
		operator = this_.token
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
			if seq, ok := left.(*node.SequenceExpression); ok && len(this_.errors) == 0 {
				paramList = this_.reinterpretSequenceAsArrowFuncParams(seq.Sequence)
			} else {
				this_.restore(&state)
				paramList = this_.parseFunctionParameterList()
			}
		} else if async {
			// async (x, y) => ...
			if !this_.scope.allowAwait {
				this_.scope.allowAwait = true
				defer func() {
					this_.scope.allowAwait = false
				}()
			}
			if _, ok := left.(*node.CallExpression); ok {
				this_.restore(&state)
				this_.next() // skip "async"
				paramList = this_.parseFunctionParameterList()
			}
		}
		if paramList == nil {
			_ = this_.error("parseAssignmentExpression paramList is empty ", left.Start(), "Malformed arrow function parameter list")
			return &node.BadExpression{From: left.Start(), To: left.End()}
		}
		return this_.parseArrowFunction(start, paramList, async)
	}

	if operator != "" {
		idx := this_.idx
		this_.next()
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
		_ = this_.error("parseAssignmentExpression", left.Start(), "Invalid left-hand side in assignment")
		this_.nextStatement()
		return &node.BadExpression{From: idx, To: this_.idx}
	}

	return left
}

func (this_ *parser) parseYieldExpression() node.Expression {
	idx := this_.expect(token.Yield)

	if this_.scope.inFuncParams {
		_ = this_.error("parseYieldExpression this_.scope.inFuncParams:true", idx, "Yield expression not allowed in formal parameter")
	}

	res := &node.YieldExpression{
		Yield: idx,
	}

	if !this_.implicitSemicolon && this_.token == token.Multiply {
		res.Delegate = true
		this_.next()
	}

	if !this_.implicitSemicolon && this_.token != token.Semicolon && this_.token != token.RightBrace && this_.token != token.Eof {
		var state parserState
		this_.mark(&state)
		expr := this_.parseAssignmentExpression()
		if _, bad := expr.(*node.BadExpression); bad {
			expr = nil
			this_.restore(&state)
		}
		res.Argument = expr
	}

	return res
}

func (this_ *parser) parseExpression() node.Expression {
	left := this_.parseAssignmentExpression()

	if this_.token == token.Comma {
		sequence := []node.Expression{left}
		for {
			if this_.token != token.Comma {
				break
			}
			this_.next()
			sequence = append(sequence, this_.parseAssignmentExpression())
		}
		return &node.SequenceExpression{
			Sequence: sequence,
		}
	}

	return left
}

func (this_ *parser) checkComma(from, to int) {
	if pos := strings.IndexByte(this_.str[(from):(to)], ','); pos >= 0 {
		_ = this_.error("checkComma", from+(pos), "Comma is not allowed here")
	}
}

func (this_ *parser) reinterpretAsArrayAssignmentPattern(left *node.ArrayLiteral) node.Expression {
	value := left.Value
	var rest node.Expression
	for i, item := range value {
		if spread, ok := item.(*node.SpreadElement); ok {
			if i != len(value)-1 {
				_ = this_.error("reinterpretAsArrayAssignmentPattern", item.Start(), "Rest element must be last element")
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

func (this_ *parser) reinterpretArrayAssignPatternAsBinding(pattern *node.ArrayPattern) *node.ArrayPattern {
	for i, item := range pattern.Elements {
		pattern.Elements[i] = this_.reinterpretAsDestructBindingTarget(item)
	}
	if pattern.Rest != nil {
		pattern.Rest = this_.reinterpretAsDestructBindingTarget(pattern.Rest)
	}
	return pattern
}

func (this_ *parser) reinterpretAsArrayBindingPattern(left *node.ArrayLiteral) node.BindingTarget {
	value := left.Value
	var rest node.Expression
	for i, item := range value {
		if spread, ok := item.(*node.SpreadElement); ok {
			if i != len(value)-1 {
				_ = this_.error("reinterpretAsArrayBindingPattern", item.Start(), "Rest element must be last element")
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

func (this_ *parser) parseArrayBindingPattern() node.BindingTarget {
	return this_.reinterpretAsArrayBindingPattern(this_.parseArrayLiteral())
}

func (this_ *parser) parseObjectBindingPattern() node.BindingTarget {
	return this_.reinterpretAsObjectBindingPattern(this_.parseObjectLiteral())
}

func (this_ *parser) reinterpretArrayObjectPatternAsBinding(pattern *node.ObjectPattern) *node.ObjectPattern {
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

func (this_ *parser) reinterpretAsObjectBindingPattern(expr *node.ObjectLiteral) node.BindingTarget {
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
				_ = this_.error("reinterpretAsObjectBindingPattern", prop.Start(), "Rest element must be last element")
				return &node.BadExpression{From: expr.Start(), To: expr.End()}
			}
			// TODO make sure there is no trailing Comma
			rest = this_.reinterpretAsBindingRestElement(prop.Expression)
			value = value[:i]
			ok = true
		}
		if !ok {
			_ = this_.error("reinterpretAsObjectBindingPattern", prop.Start(), "Invalid destructuring binding target")
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

func (this_ *parser) reinterpretAsObjectAssignmentPattern(l *node.ObjectLiteral) node.Expression {
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
				_ = this_.error("reinterpretAsObjectAssignmentPattern", prop.Start(), "Rest element must be last element")
				return &node.BadExpression{From: l.Start(), To: l.End()}
			}
			// TODO make sure there is no trailing Comma
			rest = prop.Expression
			value = value[:i]
			ok = true
		}
		if !ok {
			_ = this_.error("reinterpretAsObjectAssignmentPattern", prop.Start(), "Invalid destructuring assignment target")
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

func (this_ *parser) reinterpretAsAssignmentElement(expr node.Expression) node.Expression {
	switch expr := expr.(type) {
	case *node.AssignExpression:
		if expr.Operator == token.Assign {
			expr.Left = this_.reinterpretAsDestructAssignTarget(expr.Left)
			return expr
		} else {
			_ = this_.error("reinterpretAsAssignmentElement", expr.Start(), "Invalid destructuring assignment target")
			return &node.BadExpression{From: expr.Start(), To: expr.End()}
		}
	default:
		return this_.reinterpretAsDestructAssignTarget(expr)
	}
}

func (this_ *parser) reinterpretAsBindingElement(expr node.Expression) node.Expression {
	switch expr := expr.(type) {
	case *node.AssignExpression:
		if expr.Operator == token.Assign {
			expr.Left = this_.reinterpretAsDestructBindingTarget(expr.Left)
			return expr
		} else {
			_ = this_.error("reinterpretAsBindingElement", expr.Start(), "Invalid destructuring assignment target")
			return &node.BadExpression{From: expr.Start(), To: expr.End()}
		}
	default:
		return this_.reinterpretAsDestructBindingTarget(expr)
	}
}

func (this_ *parser) reinterpretAsBinding(expr node.Expression) *node.Binding {
	switch expr := expr.(type) {
	case *node.AssignExpression:
		if expr.Operator == token.Assign {
			return &node.Binding{
				Target:      this_.reinterpretAsDestructBindingTarget(expr.Left),
				Initializer: expr.Right,
			}
		} else {
			_ = this_.error("reinterpretAsBinding", expr.Start(), "Invalid destructuring assignment target")
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

func (this_ *parser) reinterpretAsDestructAssignTarget(item node.Expression) node.Expression {
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
	_ = this_.error("reinterpretAsDestructAssignTarget", item.Start(), "Invalid destructuring assignment target")
	return &node.BadExpression{From: item.Start(), To: item.End()}
}

func (this_ *parser) reinterpretAsDestructBindingTarget(item node.Expression) node.BindingTarget {
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
		if !this_.scope.allowAwait || item.Name != "await" {
			return item
		}
	}
	_ = this_.error("reinterpretAsDestructBindingTarget", item.Start(), "Invalid destructuring binding target")
	return &node.BadExpression{From: item.Start(), To: item.End()}
}

func (this_ *parser) reinterpretAsBindingRestElement(expr node.Expression) node.Expression {
	if _, ok := expr.(*node.Identifier); ok {
		return expr
	}
	_ = this_.error("reinterpretAsBindingRestElement", expr.Start(), "Invalid binding rest")
	return &node.BadExpression{From: expr.Start(), To: expr.End()}
}
