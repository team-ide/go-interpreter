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
		"if": {
			Token: token.If,
		},
		"in": {
			Token: token.In,
		},
		"do": {
			Token: token.Do,
		},
		"for": {
			Token: token.For,
		},
		"new": {
			Token: token.New,
		},
		"try": {
			Token: token.Try,
		},
		"this": {
			Token: token.This,
		},
		"else": {
			Token: token.Else,
		},
		"case": {
			Token: token.Case,
		},
		"void": {
			Token: token.Void,
		},
		"while": {
			Token: token.While,
		},
		"break": {
			Token: token.Break,
		},
		"catch": {
			Token: token.Catch,
		},
		"throw": {
			Token: token.Throw,
		},
		"return": {
			Token: token.Return,
		},
		"switch": {
			Token: token.Switch,
		},
		"default": {
			Token: token.Default,
		},
		"finally": {
			Token: token.Finally,
		},
		"continue": {
			Token: token.Continue,
		},
		"instanceof": {
			Token: token.Instanceof,
		},
		"class": {
			Token: token.Class,
		},
		"enum": {
			Token: token.Enum,
		},
		"super": {
			Token: token.Super,
		},
		/*
			"implements": {
				Token:         KEYWORD,
				FutureKeyword: true,
				Strict:        true,
			},
			"interface": {
				Token:         KEYWORD,
				FutureKeyword: true,
				Strict:        true,
			},*/
		"let": {
			Token:  token.Let,
			Strict: true,
		},
		/*"package": {
			Token:         KEYWORD,
			FutureKeyword: true,
			Strict:        true,
		},*/
		"private": {
			Token: token.Private,
		},
		"protected": {
			Token: token.Protected,
		},
		"public": {
			Token: token.Public,
		},
		"final": {
			Token: token.Final,
		},
		"static": {
			Token:  token.Static,
			Strict: true,
		},
		"await": {
			Token: token.Await,
		},
		"yield": {
			Token: token.Yield,
		},
		"false": {
			Token: token.Boolean,
		},
		"true": {
			Token: token.Boolean,
		},
		"null": {
			Token: token.Null,
		},
	}

	IdentifierTokens = []token.Token{
		token.Identifier,
		token.Keyword,
		token.Boolean,
		token.Null,

		token.If,
		token.In,
		token.Of,
		token.Do,

		token.Var,
		token.For,
		token.New,
		token.Try,

		token.This,
		token.Else,
		token.Case,
		token.Void,
		token.With,

		token.Const,
		token.While,
		token.Break,
		token.Catch,
		token.Throw,
		token.Class,
		token.Super,

		token.Return,
		token.Typeof,
		token.Delete,
		token.Switch,

		token.Default,
		token.Finally,
		token.Extends,

		token.Function,
		token.Continue,
		token.Debugger,

		token.Instanceof,

		token.EscapedReservedWord,
		// Non-reserved keywords below

		token.Let,
		token.Static,
		token.Async,
		token.Await,
		token.Yield,
	}

	UnreservedWordTokens = []token.Token{
		token.Let,
		token.Static,
		token.Async,
		token.Await,
		token.Yield,
	}

	ModifierTokens = []token.Token{
		token.Public,
		token.Private,
		token.Protected,
		token.Final,
	}
)
