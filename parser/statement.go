package parser

import (
	"fmt"
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/token"
)

func (this_ *parser) parseBlockStatement() *node.BlockStatement {
	res := &node.BlockStatement{}
	res.LeftBrace = this_.expect(token.LeftBrace)
	res.List = this_.parseStatementList()
	res.RightBrace = this_.expect(token.RightBrace)

	return res
}

func (this_ *parser) parseEmptyStatement() node.Statement {
	idx := this_.expect(token.Semicolon)
	return &node.EmptyStatement{Semicolon: idx}
}

func (this_ *parser) parseStatementList() (list []node.Statement) {
	for this_.token != token.RightBrace && this_.token != token.Eof {
		this_.scope.allowLet = true
		list = append(list, this_.parseStatement())
	}

	return
}

func (this_ *parser) parseStatement() node.Statement {

	if this_.token == token.Eof {
		_ = this_.errorUnexpectedToken(this_.token)
		return &node.BadStatement{From: this_.idx, To: this_.idx + 1}
	}

	switch this_.token {
	case token.Semicolon:
		return this_.parseEmptyStatement()
	case token.LeftBrace:
		return this_.parseBlockStatement()
	case token.If:
		return this_.parseIfStatement()
	case token.Do:
		return this_.parseDoWhileStatement()
	case token.While:
		return this_.parseWhileStatement()
	case token.For:
		return this_.parseForOrForInStatement()
	case token.Break:
		return this_.parseBreakStatement()
	case token.Continue:
		return this_.parseContinueStatement()
	case token.Debugger:
		return this_.parseDebuggerStatement()
	case token.With:
		return this_.parseWithStatement()
	case token.Var:
		return this_.parseVariableStatement()
	case token.Let:
		tok := this_.peek()
		if tok == token.LeftBracket || this_.scope.allowLet && (this_.IsIdentifierToken(tok) || tok == token.LeftBrace) {
			return this_.parseLexicalDeclaration(this_.token)
		}
		this_.insertSemicolon = true
	case token.Const:
		return this_.parseLexicalDeclaration(this_.token)
	case token.Async:
		if f := this_.parseMaybeAsyncFunction(true); f != nil {
			return &node.FunctionDeclaration{
				Function: f,
			}
		}
	case token.Function:
		return &node.FunctionDeclaration{
			Function: this_.parseFunction(true, false, this_.idx),
		}
	case token.Class:
		return &node.ClassDeclaration{
			Class: this_.parseClass(true),
		}
	case token.Switch:
		return this_.parseSwitchStatement()
	case token.Return:
		return this_.parseReturnStatement()
	case token.Throw:
		return this_.parseThrowStatement()
	case token.Try:
		return this_.parseTryStatement()
	}

	expression := this_.parseExpression()

	if identifier, isIdentifier := expression.(*node.Identifier); isIdentifier && this_.token == token.Colon {
		// LabelledStatement
		colon := this_.idx
		this_.next() // :
		label := identifier.Name
		for _, value := range this_.scope.labels {
			if label == value {
				_ = this_.error(identifier.Start(), fmt.Sprintf("Label '%s' already exists", label))
			}
		}
		this_.scope.labels = append(this_.scope.labels, label) // Push the label
		this_.scope.allowLet = false
		statement := this_.parseStatement()
		this_.scope.labels = this_.scope.labels[:len(this_.scope.labels)-1] // Pop the label
		return &node.LabelledStatement{
			Label:     identifier,
			Colon:     colon,
			Statement: statement,
		}
	}

	this_.optionalSemicolon()

	return &node.ExpressionStatement{
		Expression: expression,
	}
}

func (this_ *parser) parseTryStatement() node.Statement {

	res := &node.TryStatement{
		Try:  this_.expect(token.Try),
		Body: this_.parseBlockStatement(),
	}

	if this_.token == token.Catch {
		catch := this_.idx
		this_.next()
		var parameter node.BindingTarget
		if this_.token == token.LeftParenthesis {
			this_.next()
			parameter = this_.parseBindingTarget()
			this_.expect(token.RightParenthesis)
		}
		res.Catch = &node.CatchStatement{
			Catch:     catch,
			Parameter: parameter,
			Body:      this_.parseBlockStatement(),
		}
	}

	if this_.token == token.Finally {
		this_.next()
		res.Finally = this_.parseBlockStatement()
	}

	if res.Catch == nil && res.Finally == nil {
		_ = this_.error(res.Try, "Missing catch or finally after try")
		return &node.BadStatement{From: res.Try, To: res.Body.End()}
	}

	return res
}

func (this_ *parser) parseFunctionParameterList() *node.ParameterList {
	opening := this_.expect(token.LeftParenthesis)
	var list []*node.Binding
	var rest node.Expression
	if !this_.scope.inFuncParams {
		this_.scope.inFuncParams = true
		defer func() {
			this_.scope.inFuncParams = false
		}()
	}
	for this_.token != token.RightParenthesis && this_.token != token.Eof {
		if this_.token == token.Ellipsis {
			this_.next()
			rest = this_.reinterpretAsDestructBindingTarget(this_.parseAssignmentExpression())
			break
		}
		this_.parseVariableDeclaration(&list)
		if this_.token != token.RightParenthesis {
			this_.expect(token.Comma)
		}
	}
	closing := this_.expect(token.RightParenthesis)

	return &node.ParameterList{
		Opening: opening,
		List:    list,
		Rest:    rest,
		Closing: closing,
	}
}

func (this_ *parser) parseMaybeAsyncFunction(declaration bool) *node.FunctionLiteral {
	if this_.peek() == token.Function {
		idx := this_.idx
		this_.next()
		return this_.parseFunction(declaration, true, idx)
	}
	return nil
}

func (this_ *parser) parseFunction(declaration, async bool, start int) *node.FunctionLiteral {

	res := &node.FunctionLiteral{
		Function: start,
		Async:    async,
	}
	this_.expect(token.Function)

	if this_.token == token.Multiply {
		res.Generator = true
		this_.next()
	}

	if !declaration {
		if async != this_.scope.allowAwait {
			this_.scope.allowAwait = async
			defer func() {
				this_.scope.allowAwait = !async
			}()
		}
		if res.Generator != this_.scope.allowYield {
			this_.scope.allowYield = res.Generator
			defer func() {
				this_.scope.allowYield = !res.Generator
			}()
		}
	}

	this_.tokenToBindingId()
	var name *node.Identifier
	if this_.token == token.Identifier {
		name = this_.parseIdentifier()
	} else if declaration {
		// Use expect error handling
		this_.expect(token.Identifier)
	}
	res.Name = name

	if declaration {
		if async != this_.scope.allowAwait {
			this_.scope.allowAwait = async
			defer func() {
				this_.scope.allowAwait = !async
			}()
		}
		if res.Generator != this_.scope.allowYield {
			this_.scope.allowYield = res.Generator
			defer func() {
				this_.scope.allowYield = !res.Generator
			}()
		}
	}

	res.ParameterList = this_.parseFunctionParameterList()
	res.Body, res.DeclarationList = this_.parseFunctionBlock(async, async, this_.scope.allowYield)
	res.Source = this_.slice(res.Start(), res.End())

	return res
}

func (this_ *parser) parseFunctionBlock(async, allowAwait, allowYield bool) (body *node.BlockStatement, declarationList []*node.VariableDeclaration) {
	this_.openScope()
	this_.scope.inFunction = true
	this_.scope.inAsync = async
	this_.scope.allowAwait = allowAwait
	this_.scope.allowYield = allowYield
	defer this_.closeScope()
	body = this_.parseBlockStatement()
	declarationList = this_.scope.declarationList
	return
}

func (this_ *parser) parseArrowFunctionBody(async bool) (node.ConciseBody, []*node.VariableDeclaration) {
	if this_.token == token.LeftBrace {
		return this_.parseFunctionBlock(async, async, false)
	}
	if async != this_.scope.inAsync || async != this_.scope.allowAwait {
		inAsync := this_.scope.inAsync
		allowAwait := this_.scope.allowAwait
		this_.scope.inAsync = async
		this_.scope.allowAwait = async
		allowYield := this_.scope.allowYield
		this_.scope.allowYield = false
		defer func() {
			this_.scope.inAsync = inAsync
			this_.scope.allowAwait = allowAwait
			this_.scope.allowYield = allowYield
		}()
	}

	return &node.ExpressionBody{
		Expression: this_.parseAssignmentExpression(),
	}, nil
}

func (this_ *parser) parseClass(declaration bool) *node.ClassLiteral {
	if !this_.scope.allowLet && this_.token == token.Class {
		_ = this_.errorUnexpectedToken(token.Class)
	}

	res := &node.ClassLiteral{
		Class: this_.expect(token.Class),
	}

	this_.tokenToBindingId()
	var name *node.Identifier
	if this_.token == token.Identifier {
		name = this_.parseIdentifier()
	} else if declaration {
		// Use expect error handling
		this_.expect(token.Identifier)
	}

	res.Name = name

	if this_.token != token.LeftBrace {
		this_.expect(token.Extends)
		res.SuperClass = this_.parseLeftHandSideExpressionAllowCall()
	}

	this_.expect(token.LeftBrace)

	for this_.token != token.RightBrace && this_.token != token.Eof {
		if this_.token == token.Semicolon {
			this_.next()
			continue
		}
		start := this_.idx
		static := false
		if this_.token == token.Static {
			switch this_.peek() {
			case token.Assign, token.Semicolon, token.RightBrace, token.LeftParenthesis:
				// treat as identifier
			default:
				this_.next()
				if this_.token == token.LeftBrace {
					b := &node.ClassStaticBlock{
						Static: start,
					}
					b.Block, b.DeclarationList = this_.parseFunctionBlock(false, true, false)
					b.Source = this_.slice(b.Block.LeftBrace, b.Block.End())
					res.Body = append(res.Body, b)
					continue
				}
				static = true
			}
		}

		var kind node.PropertyKind
		var async bool
		methodBodyStart := this_.idx
		if this_.literal == "get" || this_.literal == "set" {
			if tok := this_.peek(); tok != token.Semicolon && tok != token.LeftParenthesis {
				if this_.literal == "get" {
					kind = node.PropertyKindGet
				} else {
					kind = node.PropertyKindSet
				}
				this_.next()
			}
		} else if this_.token == token.Async {
			if tok := this_.peek(); tok != token.Semicolon && tok != token.LeftParenthesis {
				async = true
				kind = node.PropertyKindMethod
				this_.next()
			}
		}
		generator := false
		if this_.token == token.Multiply && (kind == "" || kind == node.PropertyKindMethod) {
			generator = true
			kind = node.PropertyKindMethod
			this_.next()
		}

		_, keyName, value, tkn := this_.parseObjectPropertyKey()
		if value == nil {
			continue
		}
		computed := tkn == token.Illegal
		_, private := value.(*node.PrivateIdentifier)

		if static && !private && keyName == "prototype" {
			_ = this_.error(value.Start(), "Classes may not have a static property named 'prototype'")
		}

		if kind == "" && this_.token == token.LeftParenthesis {
			kind = node.PropertyKindMethod
		}

		if kind != "" {
			// method
			if keyName == "constructor" && !computed {
				if !static {
					if kind != node.PropertyKindMethod {
						_ = this_.error(value.Start(), "Class constructor may not be an accessor")
					} else if async {
						_ = this_.error(value.Start(), "Class constructor may not be an async method")
					} else if generator {
						_ = this_.error(value.Start(), "Class constructor may not be a generator")
					}
				} else if private {
					_ = this_.error(value.Start(), "Class constructor may not be a private method")
				}
			}
			md := &node.MethodDefinition{
				Idx:      start,
				Key:      value,
				Kind:     kind,
				Body:     this_.parseMethodDefinition(methodBodyStart, kind, generator, async),
				Static:   static,
				Computed: computed,
			}
			res.Body = append(res.Body, md)
		} else {
			// field
			isCtor := !computed && keyName == "constructor"
			if !isCtor {
				if name, ok := value.(*node.PrivateIdentifier); ok {
					isCtor = name.Name == "constructor"
				}
			}
			if isCtor {
				_ = this_.error(value.Start(), "Classes may not have a field named 'constructor'")
			}
			var initializer node.Expression
			if this_.token == token.Assign {
				this_.next()
				initializer = this_.parseExpression()
			}

			if !this_.implicitSemicolon && this_.token != token.Semicolon && this_.token != token.RightBrace {
				_ = this_.errorUnexpectedToken(this_.token)
				break
			}
			res.Body = append(res.Body, &node.FieldDefinition{
				Idx:         start,
				Key:         value,
				Initializer: initializer,
				Static:      static,
				Computed:    computed,
			})
		}
	}

	res.RightBrace = this_.expect(token.RightBrace)
	res.Source = this_.slice(res.Class, res.RightBrace+1)

	return res
}

func (this_ *parser) parseDebuggerStatement() node.Statement {
	idx := this_.expect(token.Debugger)

	res := &node.DebuggerStatement{
		Debugger: idx,
	}

	this_.semicolon()

	return res
}

func (this_ *parser) parseReturnStatement() node.Statement {
	idx := this_.expect(token.Return)

	if !this_.scope.inFunction {
		_ = this_.error(idx, "Illegal return statement")
		this_.nextStatement()
		return &node.BadStatement{From: idx, To: this_.idx}
	}

	res := &node.ReturnStatement{
		Return: idx,
	}

	if !this_.implicitSemicolon && this_.token != token.Semicolon && this_.token != token.RightBrace && this_.token != token.Eof {
		res.Argument = this_.parseExpression()
	}

	this_.semicolon()

	return res
}

func (this_ *parser) parseThrowStatement() node.Statement {
	idx := this_.expect(token.Throw)

	if this_.implicitSemicolon {
		if this_.chr == -1 { // Hackish
			_ = this_.error(idx, "Unexpected end of input")
		} else {
			_ = this_.error(idx, "Illegal newline after throw")
		}
		this_.nextStatement()
		return &node.BadStatement{From: idx, To: this_.idx}
	}

	res := &node.ThrowStatement{
		Throw:    idx,
		Argument: this_.parseExpression(),
	}

	this_.semicolon()

	return res
}

func (this_ *parser) parseSwitchStatement() node.Statement {
	this_.expect(token.Switch)
	this_.expect(token.LeftParenthesis)
	res := &node.SwitchStatement{
		Discriminant: this_.parseExpression(),
		Default:      -1,
	}
	this_.expect(token.RightParenthesis)

	this_.expect(token.LeftBrace)

	inSwitch := this_.scope.inSwitch
	this_.scope.inSwitch = true
	defer func() {
		this_.scope.inSwitch = inSwitch
	}()

	for index := 0; this_.token != token.Eof; index++ {
		if this_.token == token.RightBrace {
			this_.next()
			break
		}

		clause := this_.parseCaseStatement()
		if clause.Test == nil {
			if res.Default != -1 {
				_ = this_.error(clause.Case, "Already saw a default in switch")
			}
			res.Default = index
		}
		res.Body = append(res.Body, clause)
	}

	return res
}

func (this_ *parser) parseWithStatement() node.Statement {
	this_.expect(token.With)
	this_.expect(token.LeftParenthesis)
	res := &node.WithStatement{
		Object: this_.parseExpression(),
	}
	this_.expect(token.RightParenthesis)
	this_.scope.allowLet = false
	res.Body = this_.parseStatement()

	return res
}

func (this_ *parser) parseCaseStatement() *node.CaseStatement {

	res := &node.CaseStatement{
		Case: this_.idx,
	}
	if this_.token == token.Default {
		this_.next()
	} else {
		this_.expect(token.Case)
		res.Test = this_.parseExpression()
	}
	this_.expect(token.Colon)

	for {
		if this_.token == token.Eof ||
			this_.token == token.RightBrace ||
			this_.token == token.Case ||
			this_.token == token.Default {
			break
		}
		this_.scope.allowLet = true
		res.Consequent = append(res.Consequent, this_.parseStatement())

	}

	return res
}

func (this_ *parser) parseIterationStatement() node.Statement {
	inIteration := this_.scope.inIteration
	this_.scope.inIteration = true
	defer func() {
		this_.scope.inIteration = inIteration
	}()
	this_.scope.allowLet = false
	return this_.parseStatement()
}

func (this_ *parser) parseForIn(idx int, into node.ForInto) *node.ForInStatement {

	// Already have consumed "<into> in"

	source := this_.parseExpression()
	this_.expect(token.RightParenthesis)

	return &node.ForInStatement{
		For:    idx,
		Into:   into,
		Source: source,
		Body:   this_.parseIterationStatement(),
	}
}

func (this_ *parser) parseForOf(idx int, into node.ForInto) *node.ForOfStatement {

	// Already have consumed "<into> of"

	source := this_.parseAssignmentExpression()
	this_.expect(token.RightParenthesis)

	return &node.ForOfStatement{
		For:    idx,
		Into:   into,
		Source: source,
		Body:   this_.parseIterationStatement(),
	}
}

func (this_ *parser) parseFor(idx int, initializer node.ForLoopInitializer) *node.ForStatement {

	// Already have consumed "<initializer> ;"

	var test, update node.Expression

	if this_.token != token.Semicolon {
		test = this_.parseExpression()
	}
	this_.expect(token.Semicolon)

	if this_.token != token.RightParenthesis {
		update = this_.parseExpression()
	}
	this_.expect(token.RightParenthesis)

	return &node.ForStatement{
		For:         idx,
		Initializer: initializer,
		Test:        test,
		Update:      update,
		Body:        this_.parseIterationStatement(),
	}
}

func (this_ *parser) parseForOrForInStatement() node.Statement {
	idx := this_.expect(token.For)
	this_.expect(token.LeftParenthesis)

	var initializer node.ForLoopInitializer

	forIn := false
	forOf := false
	var into node.ForInto
	if this_.token != token.Semicolon {

		allowIn := this_.scope.allowIn
		this_.scope.allowIn = false
		tok := this_.token
		if tok == token.Let {
			switch this_.peek() {
			case token.Identifier, token.LeftBracket, token.LeftBrace:
			default:
				tok = token.Identifier
			}
		}
		if tok == token.Var || tok == token.Let || tok == token.Const {
			idx := this_.idx
			this_.next()
			var list []*node.Binding
			if tok == token.Var {
				list = this_.parseVarDeclarationList(idx)
			} else {
				list = this_.parseVariableDeclarationList()
			}
			if len(list) == 1 {
				if this_.token == token.In {
					this_.next() // in
					forIn = true
				} else if this_.token == token.Identifier && this_.literal == "of" {
					this_.next()
					forOf = true
				}
			}
			if forIn || forOf {
				if list[0].Initializer != nil {
					_ = this_.error(list[0].Initializer.Start(), "for-in loop variable declaration may not have an initializer")
				}
				if tok == token.Var {
					into = &node.ForIntoVar{
						Binding: list[0],
					}
				} else {
					into = &node.ForDeclaration{
						Idx:     idx,
						IsConst: tok == token.Const,
						Target:  list[0].Target,
					}
				}
			} else {
				this_.ensurePatternInit(list)
				if tok == token.Var {
					initializer = &node.ForLoopInitializerVarDeclList{
						List: list,
					}
				} else {
					initializer = &node.ForLoopInitializerLexicalDecl{
						LexicalDeclaration: node.LexicalDeclaration{
							Idx:   idx,
							Token: tok,
							List:  list,
						},
					}
				}
			}
		} else {
			expr := this_.parseExpression()
			if this_.token == token.In {
				this_.next()
				forIn = true
			} else if this_.token == token.Identifier && this_.literal == "of" {
				this_.next()
				forOf = true
			}
			if forIn || forOf {
				switch e := expr.(type) {
				case *node.Identifier, *node.DotExpression, *node.PrivateDotExpression, *node.BracketExpression, *node.Binding:
					// These are all acceptable
				case *node.ObjectLiteral:
					expr = this_.reinterpretAsObjectAssignmentPattern(e)
				case *node.ArrayLiteral:
					expr = this_.reinterpretAsArrayAssignmentPattern(e)
				default:
					_ = this_.error(idx, "Invalid left-hand side in for-in or for-of")
					this_.nextStatement()
					return &node.BadStatement{From: idx, To: this_.idx}
				}
				into = &node.ForIntoExpression{
					Expression: expr,
				}
			} else {
				initializer = &node.ForLoopInitializerExpression{
					Expression: expr,
				}
			}
		}
		this_.scope.allowIn = allowIn
	}

	if forIn {
		return this_.parseForIn(idx, into)
	}
	if forOf {
		return this_.parseForOf(idx, into)
	}

	this_.expect(token.Semicolon)
	return this_.parseFor(idx, initializer)
}

func (this_ *parser) ensurePatternInit(list []*node.Binding) {
	for _, item := range list {
		if _, ok := item.Target.(node.Pattern); ok {
			if item.Initializer == nil {
				_ = this_.error(item.End(), "Missing initializer in destructuring declaration")
				break
			}
		}
	}
}

func (this_ *parser) parseVariableStatement() *node.VariableStatement {

	idx := this_.expect(token.Var)

	list := this_.parseVarDeclarationList(idx)
	this_.ensurePatternInit(list)
	this_.semicolon()

	return &node.VariableStatement{
		Var:  idx,
		List: list,
	}
}

func (this_ *parser) parseLexicalDeclaration(tok token.Token) *node.LexicalDeclaration {
	idx := this_.expect(tok)
	if !this_.scope.allowLet {
		_ = this_.error(idx, "Lexical declaration cannot appear in a single-statement context")
	}

	list := this_.parseVariableDeclarationList()
	this_.ensurePatternInit(list)
	this_.semicolon()

	return &node.LexicalDeclaration{
		Idx:   idx,
		Token: tok,
		List:  list,
	}
}

func (this_ *parser) parseDoWhileStatement() node.Statement {
	inIteration := this_.scope.inIteration
	this_.scope.inIteration = true
	defer func() {
		this_.scope.inIteration = inIteration
	}()

	this_.expect(token.Do)
	res := &node.DoWhileStatement{}
	if this_.token == token.LeftBrace {
		res.Body = this_.parseBlockStatement()
	} else {
		this_.scope.allowLet = false
		res.Body = this_.parseStatement()
	}

	this_.expect(token.While)
	this_.expect(token.LeftParenthesis)
	res.Test = this_.parseExpression()
	this_.expect(token.RightParenthesis)
	if this_.token == token.Semicolon {
		this_.next()
	}

	return res
}

func (this_ *parser) parseWhileStatement() node.Statement {
	this_.expect(token.While)
	this_.expect(token.LeftParenthesis)
	res := &node.WhileStatement{
		Test: this_.parseExpression(),
	}
	this_.expect(token.RightParenthesis)
	res.Body = this_.parseIterationStatement()

	return res
}

func (this_ *parser) parseIfStatement() node.Statement {
	this_.expect(token.If)
	this_.expect(token.LeftParenthesis)
	res := &node.IfStatement{
		Test: this_.parseExpression(),
	}
	this_.expect(token.RightParenthesis)

	if this_.token == token.LeftBrace {
		res.Consequent = this_.parseBlockStatement()
	} else {
		this_.scope.allowLet = false
		res.Consequent = this_.parseStatement()
	}

	if this_.token == token.Else {
		this_.next()
		this_.scope.allowLet = false
		res.Alternate = this_.parseStatement()
	}

	return res
}

func (this_ *parser) parseBreakStatement() node.Statement {
	idx := this_.expect(token.Break)
	semicolon := this_.implicitSemicolon
	if this_.token == token.Semicolon {
		semicolon = true
		this_.next()
	}

	if semicolon || this_.token == token.RightBrace {
		this_.implicitSemicolon = false
		if !this_.scope.inIteration && !this_.scope.inSwitch {
			goto illegal
		}
		return &node.BranchStatement{
			Idx:   idx,
			Token: token.Break,
		}
	}

	this_.tokenToBindingId()
	if this_.token == token.Identifier {
		identifier := this_.parseIdentifier()
		if !this_.scope.hasLabel(identifier.Name) {
			_ = this_.error(idx, fmt.Sprintf("Undefined label '%s'", identifier.Name))
			return &node.BadStatement{From: idx, To: identifier.End()}
		}
		this_.semicolon()
		return &node.BranchStatement{
			Idx:   idx,
			Token: token.Break,
			Label: identifier,
		}
	}

	this_.expect(token.Identifier)

illegal:
	this_.error(idx, "Illegal break statement")
	this_.nextStatement()
	return &node.BadStatement{From: idx, To: this_.idx}
}

func (this_ *parser) parseContinueStatement() node.Statement {
	idx := this_.expect(token.Continue)
	semicolon := this_.implicitSemicolon
	if this_.token == token.Semicolon {
		semicolon = true
		this_.next()
	}

	if semicolon || this_.token == token.RightBrace {
		this_.implicitSemicolon = false
		if !this_.scope.inIteration {
			goto illegal
		}
		return &node.BranchStatement{
			Idx:   idx,
			Token: token.Continue,
		}
	}

	this_.tokenToBindingId()
	if this_.token == token.Identifier {
		identifier := this_.parseIdentifier()
		if !this_.scope.hasLabel(identifier.Name) {
			_ = this_.error(idx, fmt.Sprintf("Undefined label '%s'", identifier.Name))
			return &node.BadStatement{From: idx, To: identifier.End()}
		}
		if !this_.scope.inIteration {
			goto illegal
		}
		this_.semicolon()
		return &node.BranchStatement{
			Idx:   idx,
			Token: token.Continue,
			Label: identifier,
		}
	}

	this_.expect(token.Identifier)

illegal:
	this_.error(idx, "Illegal continue statement")
	this_.nextStatement()
	return &node.BadStatement{From: idx, To: this_.idx}
}

// Find the next statement after an error (recover)
func (this_ *parser) nextStatement() {
	for {
		switch this_.token {
		case token.Break, token.Continue,
			token.For, token.If, token.Return, token.Switch,
			token.Var, token.Do, token.Try, token.With,
			token.While, token.Throw, token.Catch, token.Finally:
			// Return only if parser made some progress since last
			// sync or if it has not reached 10 next calls without
			// progress. Otherwise, consume at least one token to
			// avoid an endless parser loop
			if this_.idx == this_.recover.idx && this_.recover.count < 10 {
				this_.recover.count++
				return
			}
			if this_.idx > this_.recover.idx {
				this_.recover.idx = this_.idx
				this_.recover.count = 0
				return
			}
			// Reaching here indicates a parser bug, likely an
			// incorrect token list in this function, but it only
			// leads to skipping of possibly correct code if a
			// previous error is present, and thus is preferred
			// over a non-terminating parse.
		case token.Eof:
			return
		}
		this_.next()
	}
}
