package language

import "github.com/team-ide/go-interpreter/token"

/**
符号

算术符号：+，-，*，/，%

关系符号：==，!=，>，<，>=，<=

逻辑符号：&&，||，!

位运算符：&，|，^，~，<<，>>，>>>

赋值符号：=，+=，-=，*=，/=，%=，<<=，>>=，&=，^=，|=

其他符号：&（取地址），*（指针），->（Lambda 表达式），::（方法引用），? :（三元条件运算符），,（多重赋值），...（可变参数）

括号：()，{}，[]

分号：;

冒号：：

问号：？

点号：.

注释符：//（单行注释），/.../（多行注释）

除此之外，Java 还有一些特殊的符号，如 synchronized、try-catch-finally、throw、throws 等用于处理异常的关键字，以及 import、package 等用于导入和组织代码的关键字。这些符号都有特定的用途和语法规则，开发者需要熟练掌握才能正确地使用它们。
*/

/**
运算符

算术运算符：+，-，*，/，%，++，--

关系运算符：==，!=，>，<，>=，<=

逻辑运算符：&&，||，!

位运算符：&，|，^，~，<<，>>

赋值运算符：=，+=，-=，*=，/=，%=，<<=，>>=，&=，^=，|=

三元运算符：? :

instanceof 运算符：用于判断一个对象是否是某个类的实例。

空安全操作符：?.，用于避免空指针异常。

方法引用运算符：::，用于简化方法的调用。

除此之外，Java 还有一些特殊的运算符，如：

条件运算符：用于描述一组表达式与值的关系，如 switch 语句中的 case 关键字。

箭头运算符：->，用于描述 Lambda 表达式中的参数和方法体。

数组运算符：[]，用于访问数组元素。

这些运算符可以用于不同的数据类型，如整数、浮点数、布尔值、字符、字符串和对象等。
*/

/**
关键字

abstract    assert      boolean     break       byte
case        catch       char        class       const
continue    default     do          double      else
enum        extends     final       finally     float
for         goto        if          implements  import
instanceof  int         interface   long        native
new         package     private     protected   public
return      short       static      strictfp    super
switch      synchronized this        throw       throws
transient   try         void        volatile    while

*/

/**
基础类型

数值类型：

byte：8 位有符号整数类型，取值范围为 -128 到 127。

short：16 位有符号整数类型，取值范围为 -32768 到 32767。

int：32 位有符号整数类型，取值范围为 -2147483648 到 2147483647。

long：64 位有符号整数类型，取值范围为 -9223372036854775808 到 9223372036854775807。

float：32 位单精度浮点数类型。

double：64 位双精度浮点数类型。

字符类型：

char：16 位无符号 Unicode 字符类型，可以表示一个字符或者一个转义序列。
布尔类型：

boolean：表示真假值的布尔类型，只有 true 和 false 两个取值。
引用类型：

类类型：由类定义的引用类型。

接口类型：由接口定义的引用类型。

数组类型：由数组定义的引用类型。

这些基础类型都有特定的语法和使用方式，开发者需要熟练掌握才能正确地使用它们。

*/

var (
	javaKeywordToken = map[string]_keyword{
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

	javaIdentifierTokens = []token.Token{
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

	javaUnreservedWordTokens = []token.Token{
		token.Let,
		token.Static,
		token.Async,
		token.Await,
		token.Yield,
	}
)

type JavaSyntax struct {
}

func (this_ *JavaSyntax) IsDecimalDigit(chr rune) bool {
	return IsDecimalDigit(chr)
}
func (this_ *JavaSyntax) DigitValue(chr rune) int {
	return DigitValue(chr)
}
func (this_ *JavaSyntax) IsDigit(chr rune, base int) bool {
	return IsDigit(chr, base)
}
func (this_ *JavaSyntax) IsIdentifier(s string) bool {
	return IsIdentifier(s)
}
func (this_ *JavaSyntax) IsIdentifierStart(chr rune) bool {
	return IsIdentifierStart(chr)
}
func (this_ *JavaSyntax) IsIdentifierPart(chr rune) bool {
	return IsIdentifierPart(chr)
}
func (this_ *JavaSyntax) IsLineWhiteSpace(chr rune) bool {
	return IsLineWhiteSpace(chr)
}
func (this_ *JavaSyntax) IsLineTerminator(chr rune) bool {
	return IsLineTerminator(chr)
}
func (this_ *JavaSyntax) IsKeyword(literal string) (token.Token, bool) {
	if keyword, exists := javaKeywordToken[literal]; exists {
		if keyword.futureKeyword {
			return token.Keyword, keyword.strict
		}
		return keyword.token, false
	}
	return "", false
}
func (this_ *JavaSyntax) IsIdentifierToken(tkn token.Token) bool {
	return TokenIndexOf(javaIdentifierTokens, tkn) >= 0
}
func (this_ *JavaSyntax) IsUnreservedWordToken(tkn token.Token) bool {
	return TokenIndexOf(javaUnreservedWordTokens, tkn) >= 0
}
func (this_ *JavaSyntax) IsModifier(fromToken token.Token, modifierToken token.Token) bool {
	switch fromToken {
	case token.Class, token.Enum, token.Interface:
		return modifierToken == token.Public ||
			modifierToken == token.Protected ||
			modifierToken == token.Private ||
			modifierToken == token.Static ||
			modifierToken == token.Final ||
			modifierToken == token.Abstract
	case token.Field:
		return modifierToken == token.Public ||
			modifierToken == token.Protected ||
			modifierToken == token.Private ||
			modifierToken == token.Static ||
			modifierToken == token.Final ||
			modifierToken == token.Volatile ||
			modifierToken == token.Transient ||
			modifierToken == token.Strictfp
	case token.Method:
		return modifierToken == token.Public ||
			modifierToken == token.Protected ||
			modifierToken == token.Private ||
			modifierToken == token.Static ||
			modifierToken == token.Final ||
			modifierToken == token.Volatile ||
			modifierToken == token.Transient ||
			modifierToken == token.Strictfp
	}
	return false
}
