package language

import "github.com/team-ide/go-interpreter/token"

/**
符号

算术运算符：+、-、*、/、% 分别表示加、减、乘、除、取模运算。

赋值运算符：=、+=、-=、*=、/=、%= 分别表示赋值、加等于、减等于、乘等于、除等于、取模等于运算。

比较运算符：==、===、!=、!==、>、>=、<、<= 分别表示等于、恒等于、不等于、不恒等于、大于、大于等于、小于、小于等于运算。

逻辑运算符：&&、||、! 分别表示与、或、非运算。

位运算符：&、|、^、~、<<、>>、>>> 分别表示按位与、按位或、按位异或、按位取反、左移、右移、无符号右移运算。

三目运算符：? : 表示条件运算符，用于简单的条件判断。

其他运算符：typeof、delete、in、instanceof、new、void、yield 等，用于类型检查、删除对象属性、判断属性是否存在、判断对象是否为某个类型、创建对象实例、计算表达式并返回 undefined、生成迭代器等。

JavaScript 中的符号和运算符都有其特定的用法和优先级，开发者需要熟练掌握才能正确地使用它们。
*/

/**
运算符

算术运算符：+，-，*，/，%，++，--

关系运算符：==，===，!=，!==，>，<，>=，<=

逻辑运算符：&&，||，!

位运算符：&，|，^，~，<<，>>，>>>

赋值运算符：=，+=，-=，*=，/=，%=，<<=，>>=，&=，^=，|=

三元运算符：? :

instanceof 运算符：用于判断一个对象是否是某个类的实例。

in 运算符：用于判断一个对象是否包含指定的属性。

delete 运算符：用于删除对象的属性或数组中的元素。

typeof 运算符：用于返回一个值的类型。

void 运算符：用于指定表达式没有返回值。

除此之外，JavaScript 还有一些特殊的运算符，如：

箭头函数运算符：=>，用于描述箭头函数的参数和方法体。

条件运算符：用于描述一组表达式与值的关系，如 switch 语句中的 case 关键字。

点运算符和中括号运算符：用于访问对象的属性或方法。

这些运算符可以用于不同的数据类型，如整数、浮点数、布尔值、字符、字符串和对象等。
*/

/**
关键字

保留字：break、case、catch、class、const、continue、debugger、default、delete、do、else、export、extends、finally、for、function、if、import、in、instanceof、new、return、super、switch、this、throw、try、typeof、var、void、while、with。

严格模式下的保留字：implements、interface、let、package、private、protected、public、static、yield。

未来的保留字：await、enum。

这些关键字都有特定的用途和语法规则，开发者在使用时需要注意避免与关键字冲突。
*/

/**
基础类型

Undefined：表示未定义或未初始化的值。

Null：表示空对象指针。

Boolean：表示布尔值，即 true 或 false。

Number：表示数值，包括整数和浮点数，还包括特殊的值，如 Infinity、-Infinity 和 NaN。

String：表示字符串，由一组 16 位 Unicode 字符序列组成。

Symbol：表示唯一的标识符。

BigInt：表示大整数，用于处理超出 JavaScript Number 类型范围的整数。

需要注意的是，JavaScript 中的基础类型（除了对象类型）都是不可变的，也就是说，一旦创建就不能再修改其值，而是会创建一个新的值。
*/

var (
	javaScriptKeywordToken = map[string]_keyword{
		"if": {
			token: token.If,
		},
		"in": {
			token: token.In,
		},
		"do": {
			token: token.Do,
		},
		"var": {
			token: token.Var,
		},
		"for": {
			token: token.For,
		},
		"new": {
			token: token.New,
		},
		"try": {
			token: token.Try,
		},
		"this": {
			token: token.This,
		},
		"else": {
			token: token.Else,
		},
		"case": {
			token: token.Case,
		},
		"void": {
			token: token.Void,
		},
		"with": {
			token: token.With,
		},
		"async": {
			token: token.Async,
		},
		"while": {
			token: token.While,
		},
		"break": {
			token: token.Break,
		},
		"catch": {
			token: token.Catch,
		},
		"throw": {
			token: token.Throw,
		},
		"return": {
			token: token.Return,
		},
		"typeof": {
			token: token.Typeof,
		},
		"delete": {
			token: token.Delete,
		},
		"switch": {
			token: token.Switch,
		},
		"default": {
			token: token.Default,
		},
		"finally": {
			token: token.Finally,
		},
		"function": {
			token: token.Function,
		},
		"continue": {
			token: token.Continue,
		},
		"debugger": {
			token: token.Debugger,
		},
		"instanceof": {
			token: token.Instanceof,
		},
		"const": {
			token: token.Const,
		},
		"class": {
			token: token.Class,
		},
		"enum": {
			token:         token.Keyword,
			futureKeyword: true,
		},
		"export": {
			token:         token.Keyword,
			futureKeyword: true,
		},
		"extends": {
			token: token.Extends,
		},
		"import": {
			token:         token.Keyword,
			futureKeyword: true,
		},
		"super": {
			token: token.Super,
		},
		/*
			"implements": {
				token:         KEYWORD,
				futureKeyword: true,
				strict:        true,
			},
			"interface": {
				token:         KEYWORD,
				futureKeyword: true,
				strict:        true,
			},*/
		"let": {
			token:  token.Let,
			strict: true,
		},
		/*"package": {
			token:         KEYWORD,
			futureKeyword: true,
			strict:        true,
		},
		"private": {
			token:         KEYWORD,
			futureKeyword: true,
			strict:        true,
		},
		"protected": {
			token:         KEYWORD,
			futureKeyword: true,
			strict:        true,
		},
		"public": {
			token:         KEYWORD,
			futureKeyword: true,
			strict:        true,
		},*/
		"static": {
			token:  token.Static,
			strict: true,
		},
		"await": {
			token: token.Await,
		},
		"yield": {
			token: token.Yield,
		},
		"false": {
			token: token.Boolean,
		},
		"true": {
			token: token.Boolean,
		},
		"null": {
			token: token.Null,
		},
	}

	javaScriptIdentifierTokens = []token.Token{
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

	javaScriptUnreservedWordTokens = []token.Token{
		token.Let,
		token.Static,
		token.Async,
		token.Await,
		token.Yield,
	}
)

type JavaScriptSyntax struct {
}

func (this_ *JavaScriptSyntax) IsDecimalDigit(chr rune) bool {
	return IsDecimalDigit(chr)
}
func (this_ *JavaScriptSyntax) DigitValue(chr rune) int {
	return DigitValue(chr)
}
func (this_ *JavaScriptSyntax) IsDigit(chr rune, base int) bool {
	return IsDigit(chr, base)
}
func (this_ *JavaScriptSyntax) IsIdentifier(s string) bool {
	return IsIdentifier(s)
}
func (this_ *JavaScriptSyntax) IsIdentifierStart(chr rune) bool {
	return IsIdentifierStart(chr)
}
func (this_ *JavaScriptSyntax) IsIdentifierPart(chr rune) bool {
	return IsIdentifierPart(chr)
}
func (this_ *JavaScriptSyntax) IsLineWhiteSpace(chr rune) bool {
	return IsLineWhiteSpace(chr)
}
func (this_ *JavaScriptSyntax) IsLineTerminator(chr rune) bool {
	return IsLineTerminator(chr)
}
func (this_ *JavaScriptSyntax) IsKeyword(literal string) (token.Token, bool) {
	if keyword, exists := javaScriptKeywordToken[literal]; exists {
		if keyword.futureKeyword {
			return token.Keyword, keyword.strict
		}
		return keyword.token, false
	}
	return "", false
}
func (this_ *JavaScriptSyntax) IsIdentifierToken(tkn token.Token) bool {
	return TokenIndexOf(javaScriptIdentifierTokens, tkn) >= 0
}
func (this_ *JavaScriptSyntax) IsUnreservedWordToken(tkn token.Token) bool {
	return TokenIndexOf(javaScriptUnreservedWordTokens, tkn) >= 0
}
func (this_ *JavaScriptSyntax) IsModifier(fromToken token.Token, modifierToken token.Token) bool {
	return false
}
