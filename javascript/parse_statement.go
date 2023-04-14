package javascript

import (
	"fmt"
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/token"
)

func (this_ *Parser) parseTryStatement() node.Statement {

	res := &node.TryStatement{
		Try:  this_.ExpectAndNext("parseTryStatement", token.Try),
		Body: this_.ParseBlockStatement(),
	}

	if this_.Token == token.Catch {
		catch := this_.Idx
		this_.Next()
		var parameter node.BindingTarget
		if this_.Token == token.LeftParenthesis {
			this_.Next()
			parameter = this_.parseBindingTarget()
			this_.ExpectAndNext("parseTryStatement", token.RightParenthesis)
		}
		res.Catch = &node.CatchStatement{
			Catch:     catch,
			Parameter: parameter,
			Body:      this_.ParseBlockStatement(),
		}
	}

	if this_.Token == token.Finally {
		this_.Next()
		res.Finally = this_.ParseBlockStatement()
	}

	if res.Catch == nil && res.Finally == nil {
		_ = this_.Error("parseTryStatement", res.Try, "Missing catch or finally after try")
		return &node.BadStatement{From: res.Try, To: res.Body.End()}
	}

	return res
}

func (this_ *Parser) parseFunctionParameterList() *node.ParameterList {
	opening := this_.ExpectAndNext("parseFunctionParameterList", token.LeftParenthesis)
	var list []*node.Binding
	var rest node.Expression
	if !this_.Scope.InFuncParams {
		this_.Scope.InFuncParams = true
		defer func() {
			this_.Scope.InFuncParams = false
		}()
	}
	for this_.Token != token.RightParenthesis && this_.Token != token.Eof {
		if this_.Token == token.Ellipsis {
			this_.Next()
			rest = this_.reinterpretAsDestructBindingTarget(this_.parseAssignmentExpression())
			break
		}
		this_.parseVariableDeclaration(&list)
		if this_.Token != token.RightParenthesis {
			this_.ExpectAndNext("parseFunctionParameterList", token.Comma)
		}
	}
	closing := this_.ExpectAndNext("parseFunctionParameterList", token.RightParenthesis)

	return &node.ParameterList{
		Opening: opening,
		List:    list,
		Rest:    rest,
		Closing: closing,
	}
}

func (this_ *Parser) parseMaybeAsyncFunction(declaration bool) *node.FunctionLiteral {
	if this_.Peek() == token.Function {
		idx := this_.Idx
		this_.Next()
		return this_.parseFunction(declaration, true, idx)
	}
	return nil
}

func (this_ *Parser) parseFunction(declaration, async bool, start int) *node.FunctionLiteral {

	res := &node.FunctionLiteral{
		Function: start,
		Async:    async,
	}
	this_.ExpectAndNext("parseFunction", token.Function)

	if this_.Token == token.Multiply {
		res.Generator = true
		this_.Next()
	}

	if !declaration {
		if async != this_.Scope.AllowAwait {
			this_.Scope.AllowAwait = async
			defer func() {
				this_.Scope.AllowAwait = !async
			}()
		}
		if res.Generator != this_.Scope.AllowYield {
			this_.Scope.AllowYield = res.Generator
			defer func() {
				this_.Scope.AllowYield = !res.Generator
			}()
		}
	}

	this_.TokenToBindingIdentifier()
	var name *node.Identifier
	if this_.Token == token.Identifier {
		name = this_.ParseIdentifier()
	} else if declaration {
		// Use expect error handling
		this_.ExpectAndNext("parseFunction", token.Identifier)
	}
	res.Name = name

	if declaration {
		if async != this_.Scope.AllowAwait {
			this_.Scope.AllowAwait = async
			defer func() {
				this_.Scope.AllowAwait = !async
			}()
		}
		if res.Generator != this_.Scope.AllowYield {
			this_.Scope.AllowYield = res.Generator
			defer func() {
				this_.Scope.AllowYield = !res.Generator
			}()
		}
	}

	res.ParameterList = this_.parseFunctionParameterList()
	res.Body, res.DeclarationList = this_.parseFunctionBlock(async, async, this_.Scope.AllowYield)
	res.Source = this_.Slice(res.Start(), res.End())

	return res
}

func (this_ *Parser) parseFunctionBlock(async, allowAwait, allowYield bool) (body *node.BlockStatement, declarationList []*node.VariableDeclaration) {
	this_.OpenScope()
	defer this_.CloseScope()
	this_.Scope.InFunction = true
	this_.Scope.InAsync = async
	this_.Scope.AllowAwait = allowAwait
	this_.Scope.AllowYield = allowYield
	body = this_.ParseBlockStatement()
	declarationList = this_.Scope.DeclarationList
	return
}

func (this_ *Parser) parseArrowFunctionBody(async bool) (node.ConciseBody, []*node.VariableDeclaration) {
	if this_.Token == token.LeftBrace {
		return this_.parseFunctionBlock(async, async, false)
	}
	if async != this_.Scope.InAsync || async != this_.Scope.AllowAwait {
		inAsync := this_.Scope.InAsync
		allowAwait := this_.Scope.AllowAwait
		this_.Scope.InAsync = async
		this_.Scope.AllowAwait = async
		allowYield := this_.Scope.AllowYield
		this_.Scope.AllowYield = false
		defer func() {
			this_.Scope.InAsync = inAsync
			this_.Scope.AllowAwait = allowAwait
			this_.Scope.AllowYield = allowYield
		}()
	}

	return &node.ExpressionBody{
		Expression: this_.parseAssignmentExpression(),
	}, nil
}

func (this_ *Parser) parseClass(declaration bool) *node.ClassLiteral {
	if !this_.Scope.AllowLet && this_.Token == token.Class {
		_ = this_.ErrorUnexpectedToken("parseClass", token.Class)
	}

	res := &node.ClassLiteral{
		Class: this_.ExpectAndNext("parseClass", token.Class),
	}

	this_.TokenToBindingIdentifier()
	var name *node.Identifier
	if this_.Token == token.Identifier {
		name = this_.ParseIdentifier()
	} else if declaration {
		// Use expect error handling
		this_.ExpectAndNext("parseClass", token.Identifier)
	}

	res.Name = name

	if this_.Token != token.LeftBrace {
		this_.ExpectAndNext("parseClass", token.Extends)
		res.Extend = this_.parseLeftHandSideExpressionAllowCall()
	}

	this_.ExpectAndNext("parseClass", token.LeftBrace)

	for this_.Token != token.RightBrace && this_.Token != token.Eof {
		if this_.Token == token.Semicolon {
			this_.Next()
			continue
		}
		start := this_.Idx
		static := false
		if this_.Token == token.Static {
			switch this_.Peek() {
			case token.Assign, token.Semicolon, token.RightBrace, token.LeftParenthesis:
				// treat as identifier
			default:
				this_.Next()
				if this_.Token == token.LeftBrace {
					b := &node.ClassStaticBlock{
						Static: start,
					}
					b.Block, b.DeclarationList = this_.parseFunctionBlock(false, true, false)
					b.Source = this_.Slice(b.Block.LeftBrace, b.Block.End())
					res.Body = append(res.Body, b)
					continue
				}
				static = true
			}
		}

		var kind node.PropertyKind
		var async bool
		methodBodyStart := this_.Idx
		if this_.Literal == "get" || this_.Literal == "set" {
			if tok := this_.Peek(); tok != token.Semicolon && tok != token.LeftParenthesis {
				if this_.Literal == "get" {
					kind = node.PropertyKindGet
				} else {
					kind = node.PropertyKindSet
				}
				this_.Next()
			}
		} else if this_.Token == token.Async {
			if tok := this_.Peek(); tok != token.Semicolon && tok != token.LeftParenthesis {
				async = true
				kind = node.PropertyKindMethod
				this_.Next()
			}
		}
		generator := false
		if this_.Token == token.Multiply && (kind == "" || kind == node.PropertyKindMethod) {
			generator = true
			kind = node.PropertyKindMethod
			this_.Next()
		}

		_, keyName, value, tkn := this_.parseObjectPropertyKey()
		if value == nil {
			continue
		}
		computed := tkn == token.Illegal
		_, private := value.(*node.PrivateIdentifier)

		if static && !private && keyName == "prototype" {
			_ = this_.Error("parseClass", value.Start(), "Classes may not have a static property named 'prototype'")
		}

		if kind == "" && this_.Token == token.LeftParenthesis {
			kind = node.PropertyKindMethod
		}

		if kind != "" {
			// method
			if keyName == "constructor" && !computed {
				if !static {
					if kind != node.PropertyKindMethod {
						_ = this_.Error("parseClass", value.Start(), "Class constructor may not be an accessor")
					} else if async {
						_ = this_.Error("parseClass", value.Start(), "Class constructor may not be an async method")
					} else if generator {
						_ = this_.Error("parseClass", value.Start(), "Class constructor may not be a generator")
					}
				} else if private {
					_ = this_.Error("parseClass", value.Start(), "Class constructor may not be a private method")
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
				_ = this_.Error("parseClass", value.Start(), "Classes may not have a field named 'constructor'")
			}
			var initializer node.Expression
			if this_.Token == token.Assign {
				this_.Next()
				initializer = this_.parseExpression()
			}

			if !this_.ImplicitSemicolon && this_.Token != token.Semicolon && this_.Token != token.RightBrace {
				_ = this_.ErrorUnexpectedToken("parseClass", this_.Token)
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

	res.RightBrace = this_.ExpectAndNext("parseClass", token.RightBrace)
	res.Source = this_.Slice(res.Class, res.RightBrace+1)

	return res
}

func (this_ *Parser) parseDebuggerStatement() node.Statement {
	idx := this_.ExpectAndNext("parseDebuggerStatement", token.Debugger)

	res := &DebuggerStatement{
		Debugger: idx,
	}

	this_.Semicolon("parseDebuggerStatement")

	return res
}

func (this_ *Parser) parseReturnStatement() node.Statement {
	idx := this_.ExpectAndNext("parseReturnStatement", token.Return)

	if !this_.Scope.InFunction {
		_ = this_.Error("parseReturnStatement", idx, "Illegal return statement")
		this_.nextStatement()
		return &node.BadStatement{From: idx, To: this_.Idx}
	}

	res := &node.ReturnStatement{
		Return: idx,
	}

	if !this_.ImplicitSemicolon && this_.Token != token.Semicolon && this_.Token != token.RightBrace && this_.Token != token.Eof {
		res.Argument = this_.parseExpression()
	}

	this_.Semicolon("parseReturnStatement")

	return res
}

func (this_ *Parser) parseThrowStatement() node.Statement {
	idx := this_.ExpectAndNext("parseThrowStatement", token.Throw)

	if this_.ImplicitSemicolon {
		if this_.Chr == -1 { // Hackish
			_ = this_.Error("parseThrowStatement", idx, "Unexpected end of input")
		} else {
			_ = this_.Error("parseThrowStatement", idx, "Illegal newline after throw")
		}
		this_.nextStatement()
		return &node.BadStatement{From: idx, To: this_.Idx}
	}

	res := &node.ThrowStatement{
		Throw:    idx,
		Argument: this_.parseExpression(),
	}

	this_.Semicolon("parseThrowStatement")

	return res
}

func (this_ *Parser) parseSwitchStatement() node.Statement {
	idx := this_.ExpectAndNext("parseSwitchStatement", token.Switch)
	this_.ExpectAndNext("parseSwitchStatement", token.LeftParenthesis)
	res := &node.SwitchStatement{
		Switch:       idx,
		Discriminant: this_.parseExpression(),
		Default:      -1,
	}
	this_.ExpectAndNext("parseSwitchStatement", token.RightParenthesis)

	this_.ExpectAndNext("parseSwitchStatement", token.LeftBrace)

	inSwitch := this_.Scope.InSwitch
	this_.Scope.InSwitch = true
	defer func() {
		this_.Scope.InSwitch = inSwitch
	}()

	for index := 0; this_.Token != token.Eof; index++ {
		if this_.Token == token.RightBrace {
			this_.Next()
			break
		}

		clause := this_.parseCaseStatement()
		if clause.Test == nil {
			if res.Default != -1 {
				_ = this_.Error("parseSwitchStatement", clause.Case, "Already saw a default in switch")
			}
			res.Default = index
		}
		res.Body = append(res.Body, clause)
	}
	//bs, _ := json.Marshal(res)
	//fmt.Println("parseSwitchStatement res:", string(bs))

	return res
}

func (this_ *Parser) parseWithStatement() node.Statement {
	idx := this_.ExpectAndNext("parseWithStatement", token.With)
	this_.ExpectAndNext("parseWithStatement", token.LeftParenthesis)
	res := &node.WithStatement{
		With:   idx,
		Object: this_.parseExpression(),
	}
	this_.ExpectAndNext("parseWithStatement", token.RightParenthesis)
	this_.Scope.AllowLet = false
	res.Body = this_.parseStatement()

	return res
}

func (this_ *Parser) parseCaseStatement() *node.CaseStatement {

	res := &node.CaseStatement{
		Case: this_.Idx,
	}
	if this_.Token == token.Default {
		this_.Next()
	} else {
		this_.ExpectAndNext("parseCaseStatement", token.Case)
		res.Test = this_.parseExpression()
	}
	this_.ExpectAndNext("parseCaseStatement", token.Colon)

	for {
		if this_.Token == token.Eof ||
			this_.Token == token.RightBrace ||
			this_.Token == token.Case ||
			this_.Token == token.Default {
			break
		}
		//fmt.Println("parseCaseStatement token:", this_.token)
		this_.Scope.AllowLet = true
		state := this_.parseStatement()
		//fmt.Println("parseCaseStatement state:", reflect.TypeOf(state).String())
		res.Consequent = append(res.Consequent, state)

	}
	//fmt.Println("parseCaseStatement res.Consequent:", res.Consequent)

	return res
}

func (this_ *Parser) parseIterationStatement() node.Statement {
	inIteration := this_.Scope.InIteration
	this_.Scope.InIteration = true
	defer func() {
		this_.Scope.InIteration = inIteration
	}()
	this_.Scope.AllowLet = false
	return this_.parseStatement()
}

func (this_ *Parser) parseForIn(idx int, into node.ForInto) *node.ForInStatement {

	// Already have consumed "<into> in"

	source := this_.parseExpression()
	this_.ExpectAndNext("parseForIn", token.RightParenthesis)

	return &node.ForInStatement{
		For:    idx,
		Into:   into,
		Source: source,
		Body:   this_.parseIterationStatement(),
	}
}

func (this_ *Parser) parseForOf(idx int, into node.ForInto) *node.ForOfStatement {

	// Already have consumed "<into> of"

	source := this_.parseAssignmentExpression()
	this_.ExpectAndNext("parseForOf", token.RightParenthesis)

	return &node.ForOfStatement{
		For:    idx,
		Into:   into,
		Source: source,
		Body:   this_.parseIterationStatement(),
	}
}

func (this_ *Parser) parseFor(idx int, initializer node.ForLoopInitializer) *node.ForStatement {

	// Already have consumed "<initializer> ;"

	var test, update node.Expression

	if this_.Token != token.Semicolon {
		test = this_.parseExpression()
	}
	this_.ExpectAndNext("parseFor", token.Semicolon)

	if this_.Token != token.RightParenthesis {
		update = this_.parseExpression()
	}
	this_.ExpectAndNext("parseFor", token.RightParenthesis)

	return &node.ForStatement{
		For:         idx,
		Initializer: initializer,
		Test:        test,
		Update:      update,
		Body:        this_.parseIterationStatement(),
	}
}

func (this_ *Parser) parseForOrForInStatement() node.Statement {
	idx := this_.ExpectAndNext("parseForOrForInStatement", token.For)
	this_.ExpectAndNext("parseForOrForInStatement", token.LeftParenthesis)

	var initializer node.ForLoopInitializer

	forIn := false
	forOf := false
	var into node.ForInto
	if this_.Token != token.Semicolon {

		allowIn := this_.Scope.AllowIn
		this_.Scope.AllowIn = false
		tok := this_.Token
		if tok == token.Let {
			switch this_.Peek() {
			case token.Identifier, token.LeftBracket, token.LeftBrace:
			default:
				tok = token.Identifier
			}
		}
		if tok == token.Var || tok == token.Let || tok == token.Const {
			idx := this_.Idx
			this_.Next()
			var list []*node.Binding
			if tok == token.Var {
				list = this_.parseVarDeclarationList(idx)
			} else {
				list = this_.parseVariableDeclarationList()
			}
			if len(list) == 1 {
				if this_.Token == token.In {
					this_.Next() // in
					forIn = true
				} else if this_.Token == token.Identifier && this_.Literal == "of" {
					this_.Next()
					forOf = true
				}
			}
			if forIn || forOf {
				if list[0].Initializer != nil {
					_ = this_.Error("parseForOrForInStatement", list[0].Initializer.Start(), "for-in loop variable declaration may not have an initializer")
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
			if this_.Token == token.In {
				this_.Next()
				forIn = true
			} else if this_.Token == token.Identifier && this_.Literal == "of" {
				this_.Next()
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
					_ = this_.Error("parseForOrForInStatement", idx, "Invalid left-hand side in for-in or for-of")
					this_.nextStatement()
					return &node.BadStatement{From: idx, To: this_.Idx}
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
		this_.Scope.AllowIn = allowIn
	}

	if forIn {
		return this_.parseForIn(idx, into)
	}
	if forOf {
		return this_.parseForOf(idx, into)
	}

	this_.ExpectAndNext("parseForOrForInStatement", token.Semicolon)
	return this_.parseFor(idx, initializer)
}

// 确保模式初始化
func (this_ *Parser) ensurePatternInit(list []*node.Binding) {
	for _, item := range list {
		if _, ok := item.Target.(node.Pattern); ok {
			if item.Initializer == nil {
				_ = this_.Error("ensurePatternInit", item.End(), "Missing initializer in destructuring declaration")
				break
			}
		}
	}
}

func (this_ *Parser) parseVariableStatement() *node.VariableStatement {

	idx := this_.ExpectAndNext("parseVariableStatement", token.Var)

	list := this_.parseVarDeclarationList(idx)
	this_.ensurePatternInit(list)
	this_.Semicolon("parseVariableStatement")

	return &node.VariableStatement{
		Var:  idx,
		List: list,
	}
}

func (this_ *Parser) parseLexicalDeclaration(tok token.Token) *node.LexicalDeclaration {
	idx := this_.ExpectAndNext("parseLexicalDeclaration", tok)
	if !this_.Scope.AllowLet {
		_ = this_.Error("parseLexicalDeclaration", idx, "Lexical declaration cannot appear in a single-statement context")
	}

	list := this_.parseVariableDeclarationList()
	this_.ensurePatternInit(list)
	this_.Semicolon("parseLexicalDeclaration")

	return &node.LexicalDeclaration{
		Idx:   idx,
		Token: tok,
		List:  list,
	}
}

func (this_ *Parser) parseDoWhileStatement() node.Statement {
	inIteration := this_.Scope.InIteration
	this_.Scope.InIteration = true
	defer func() {
		this_.Scope.InIteration = inIteration
	}()

	idx := this_.ExpectAndNext("parseDoWhileStatement", token.Do)
	res := &node.DoWhileStatement{
		Do: idx,
	}
	if this_.Token == token.LeftBrace {
		res.Body = this_.ParseBlockStatement()
	} else {
		this_.Scope.AllowLet = false
		res.Body = this_.parseStatement()
	}

	this_.ExpectAndNext("parseDoWhileStatement", token.While)
	this_.ExpectAndNext("parseDoWhileStatement", token.LeftParenthesis)
	res.Test = this_.parseExpression()
	endIdx := this_.ExpectAndNext("parseDoWhileStatement", token.RightParenthesis)
	res.EndIdx = endIdx + 1
	if this_.Token == token.Semicolon {
		this_.Next()
	}

	return res
}

func (this_ *Parser) parseWhileStatement() node.Statement {
	idx := this_.ExpectAndNext("parseWhileStatement", token.While)
	this_.ExpectAndNext("parseWhileStatement", token.LeftParenthesis)
	res := &node.WhileStatement{
		While: idx,
		Test:  this_.parseExpression(),
	}
	this_.ExpectAndNext("parseWhileStatement", token.RightParenthesis)
	res.Body = this_.parseIterationStatement()

	return res
}

func (this_ *Parser) parseIfStatement() node.Statement {
	idx := this_.ExpectAndNext("parseIfStatement", token.If)
	this_.ExpectAndNext("parseIfStatement", token.LeftParenthesis)
	res := &node.IfStatement{
		If:   idx,
		Test: this_.parseExpression(),
	}
	this_.ExpectAndNext("parseIfStatement", token.RightParenthesis)

	if this_.Token == token.LeftBrace {
		res.Consequent = this_.ParseBlockStatement()
	} else {
		this_.Scope.AllowLet = false
		res.Consequent = this_.parseStatement()
	}

	if this_.Token == token.Else {
		this_.Next()
		this_.Scope.AllowLet = false
		res.Alternate = this_.parseStatement()
	}

	return res
}

func (this_ *Parser) parseBreakStatement() node.Statement {
	idx := this_.ExpectAndNext("parseBreakStatement", token.Break)
	semicolon := this_.ImplicitSemicolon
	if this_.Token == token.Semicolon {
		semicolon = true
		this_.Next()
	}

	if semicolon || this_.Token == token.RightBrace {
		this_.ImplicitSemicolon = false
		if !this_.Scope.InIteration && !this_.Scope.InSwitch {
			goto illegal
		}
		return &node.BreakStatement{
			From: idx,
			To:   this_.Idx + 1,
		}
	}

	this_.TokenToBindingIdentifier()
	if this_.Token == token.Identifier {
		identifier := this_.ParseIdentifier()
		if !this_.Scope.HasLabel(identifier.Name) {
			_ = this_.Error("parseBreakStatement", idx, fmt.Sprintf("Undefined label '%s'", identifier.Name))
			return &node.BadStatement{From: idx, To: identifier.End()}
		}
		this_.Semicolon("parseBreakStatement")
		return &node.BreakStatement{
			From:  idx,
			To:    identifier.End(),
			Label: identifier,
		}
	}

	this_.ExpectAndNext("parseBreakStatement", token.Identifier)

illegal:
	this_.Error("parseBreakStatement", idx, "Illegal break statement")
	this_.nextStatement()
	return &node.BadStatement{From: idx, To: this_.Idx}
}

func (this_ *Parser) parseContinueStatement() node.Statement {
	idx := this_.ExpectAndNext("parseContinueStatement", token.Continue)
	semicolon := this_.ImplicitSemicolon
	if this_.Token == token.Semicolon {
		semicolon = true
		this_.Next()
	}

	if semicolon || this_.Token == token.RightBrace {
		this_.ImplicitSemicolon = false
		if !this_.Scope.InIteration {
			goto illegal
		}
		return &node.ContinueStatement{
			From: idx,
			To:   this_.Idx + 1,
		}
	}

	this_.TokenToBindingIdentifier()
	if this_.Token == token.Identifier {
		identifier := this_.ParseIdentifier()
		if !this_.Scope.HasLabel(identifier.Name) {
			_ = this_.Error("parseContinueStatement", idx, fmt.Sprintf("Undefined label '%s'", identifier.Name))
			return &node.BadStatement{From: idx, To: identifier.End()}
		}
		if !this_.Scope.InIteration {
			goto illegal
		}
		this_.Semicolon("parseContinueStatement")
		return &node.ContinueStatement{
			From:  idx,
			To:    identifier.End(),
			Label: identifier,
		}
	}

	this_.ExpectAndNext("parseContinueStatement", token.Identifier)

illegal:
	this_.Error("parseContinueStatement", idx, "Illegal continue statement")
	this_.nextStatement()
	return &node.BadStatement{From: idx, To: this_.Idx}
}

// Find the next statement after an error (recover)
func (this_ *Parser) nextStatement() {
	for {
		switch this_.Token {
		case token.Break, token.Continue,
			token.For, token.If, token.Return, token.Switch,
			token.Var, token.Do, token.Try, token.With,
			token.While, token.Throw, token.Catch, token.Finally:
			// Return only if parser made some progress since last
			// sync or if it has not reached 10 next calls without
			// progress. Otherwise, consume at least one token to
			// avoid an endless parser loop
			if this_.Idx == this_.Recover.Idx && this_.Recover.Count < 10 {
				this_.Recover.Count++
				return
			}
			if this_.Idx > this_.Recover.Idx {
				this_.Recover.Idx = this_.Idx
				this_.Recover.Count = 0
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
		this_.Next()
	}
}
