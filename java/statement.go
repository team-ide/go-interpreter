package java

import (
	"fmt"
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/token"
)

func (this_ *Parser) parseImportStatement() *node.ImportStatement {

	idx := this_.ExpectAndNext("parseImportStatement", token.Import)

	res := &node.ImportStatement{
		From: idx,
	}

	imp := ""
	for {
		if this_.Token == token.Identifier {
			identifier := this_.ParseIdentifier()
			imp += string(identifier.Name)
		} else if this_.Token == token.Period {
			this_.Next()
			imp += "."
		} else if this_.Token == token.Multiply {
			this_.Next()
			imp += "*"
		} else {
			break
		}
	}
	res.Import = imp
	res.To = this_.Idx
	if token.Semicolon == this_.Token {
		res.To++
	}
	this_.ExpectAndNext("parseImportStatement", token.Semicolon)
	//fmt.Println("parseImportStatement ", "imp:", imp, ",Next token:", this_.Token)

	return res
}

func (this_ *Parser) parseClassLiteral() *node.ClassLiteral {
	modifiers := this_.GetAndClearModifiers()
	idx := this_.ExpectAndNext("parseClassLiteral", token.Class)
	if len(modifiers) > 0 {
		idx = modifiers[0].Idx
	}
	//fmt.Println("modifiers:", this_.ToJSON(modifiers))
	res := &node.ClassLiteral{
		Class: idx,
	}
	res.Name = this_.ParseIdentifier()

	var implements []node.Expression

	for {
		if this_.Token == token.Extends {
			this_.Next()
			res.Extend = this_.ParseIdentifier()
		} else if this_.Token == token.Implements {
			this_.Next()
			for {
				if this_.Token == token.Identifier {
					implements = append(implements, this_.ParseIdentifier())
				} else if this_.Token == token.Comma {
					this_.Next()
					continue
				} else {
					break
				}
			}
		} else {
			break
		}
	}
	res.Implements = implements

	this_.ExpectAndNext("parseClassLiteral", token.LeftBrace)
	// 解析 class 内容

	var identifierCache []*node.Identifier
	for this_.Token != token.RightBrace && this_.Token != token.Eof {

		// 解析 属性 或 方法 定义
		if this_.Token == token.Identifier {
			identifierCache = append(identifierCache, this_.ParseIdentifier())
			continue
		}

		// 如果是左括号，则表示 方法开始
		if this_.Token == token.LeftParenthesis {
			start := this_.Idx
			modifiers = this_.GetAndClearModifiers()
			if len(modifiers) > 0 {
				start = modifiers[0].Idx
			} else if len(identifierCache) > 0 {
				start = identifierCache[0].Idx
			}

			// 解析参数 到 ) 结束
			for this_.Token != token.RightParenthesis && this_.Token != token.Eof {
				this_.Next()
			}
			this_.Next()
			// 解析参数 到 } 结束
			for this_.Token != token.RightBrace && this_.Token != token.Eof {
				this_.Next()
			}
			md := &node.MethodDefinition{
				Idx: start,
			}
			res.Body = append(res.Body, md)

			continue
		}

		if this_.Token == token.Assign {

		}

		fmt.Println("this token:", this_.Token)
		this_.Next()
	}

	res.RightBrace = this_.ExpectAndNext("parseClassLiteral", token.RightBrace)
	res.Source = this_.Slice(res.Class, res.RightBrace+1)

	return res
}
