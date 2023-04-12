package thrift

import (
	"fmt"
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/token"
)

func (this_ *Parser) parseIncludeStatement() *node.IncludeStatement {

	idx := this_.ExpectAndNext("parseIncludeStatement", token.Include)

	res := &node.IncludeStatement{
		From: idx,
	}
	toIdx := this_.Idx

	inc := ""
	for {
		if this_.Token == token.String {
			inc = string(this_.ParsedLiteral)
			toIdx += len(this_.Literal)
			this_.ExpectAndNext("parseIncludeStatement", token.String)
		} else {
			break
		}
	}
	res.Include = inc
	res.To = toIdx
	fmt.Println("parseIncludeStatement ", res, ",Next token:", this_.Token)

	return res
}

func (this_ *Parser) parseNamespaceStatement() *node.NamespaceStatement {

	idx := this_.ExpectAndNext("parseNamespaceStatement", token.Namespace)

	res := &node.NamespaceStatement{
		From: idx,
	}
	toIdx := this_.Idx

	if this_.Token == token.Identifier {
		identifier := this_.ParseIdentifier()
		res.Language = string(identifier.Name)
	}
	namespace := ""
	for {
		if this_.Token == token.Identifier {
			identifier := this_.ParseIdentifier()
			namespace += string(identifier.Name)
			continue
		} else if this_.Token == token.Period {
			namespace += "."
			this_.Next()
			continue
		} else {
			namespace += this_.Literal
			s := this_.ImplicitRead()
			if s == '\n' || s == ';' {
				toIdx = this_.Idx
				this_.Next()
				this_.Next()
				break
			}
			this_.Next()
		}
	}
	res.Namespace = namespace
	res.To = toIdx
	fmt.Println("parseNamespaceStatement ", res, ",Next token:", this_.Token)

	return res
}

func (this_ *Parser) parseStructStatement() *node.StructStatement {

	idx := this_.ExpectAndNext("parseStructStatement", token.Struct)

	res := &node.StructStatement{
		From: idx,
	}

	if this_.Token == token.Identifier {
		identifier := this_.ParseIdentifier()
		res.Name = string(identifier.Name)
	}
	for this_.Token != token.RightBrace && this_.Token != token.Eof {
		this_.Next()
	}
	res.To = this_.Idx
	if this_.Token == token.RightBrace {
		res.To++
	}
	this_.ExpectAndNext("parseStructStatement", token.RightBrace)
	fmt.Println("parseStructStatement ", res, ",Next token:", this_.Token)

	return res
}

func (this_ *Parser) parseExceptionStatement() *node.ExceptionStatement {

	idx := this_.ExpectAndNext("parseExceptionStatement", token.Exception)

	res := &node.ExceptionStatement{
		From: idx,
	}

	if this_.Token == token.Identifier {
		identifier := this_.ParseIdentifier()
		res.Name = string(identifier.Name)
	}
	for this_.Token != token.RightBrace && this_.Token != token.Eof {
		this_.Next()
	}
	res.To = this_.Idx
	if this_.Token == token.RightBrace {
		res.To++
	}
	this_.ExpectAndNext("parseExceptionStatement", token.RightBrace)
	fmt.Println("parseExceptionStatement ", res, ",Next token:", this_.Token)

	return res
}

func (this_ *Parser) parseEnumStatement() *node.EnumStatement {

	idx := this_.ExpectAndNext("parseEnumStatement", token.Enum)

	res := &node.EnumStatement{
		From: idx,
	}

	if this_.Token == token.Identifier {
		identifier := this_.ParseIdentifier()
		res.Name = string(identifier.Name)
	}
	for this_.Token != token.RightBrace && this_.Token != token.Eof {
		this_.Next()
	}
	res.To = this_.Idx
	if this_.Token == token.RightBrace {
		res.To++
	}
	this_.ExpectAndNext("parseEnumStatement", token.RightBrace)
	fmt.Println("parseEnumStatement ", res, ",Next token:", this_.Token)
	return res
}

func (this_ *Parser) parseServiceStatement() *node.ServiceStatement {

	idx := this_.ExpectAndNext("parseServiceStatement", token.Service)

	res := &node.ServiceStatement{
		From: idx,
	}

	if this_.Token == token.Identifier {
		identifier := this_.ParseIdentifier()
		res.Name = string(identifier.Name)
	}
	for this_.Token != token.RightBrace && this_.Token != token.Eof {
		this_.Next()
	}
	res.To = this_.Idx
	if this_.Token == token.RightBrace {
		res.To++
	}
	this_.ExpectAndNext("parseServiceStatement", token.RightBrace)
	fmt.Println("parseServiceStatement ", res, ",Next token:", this_.Token)
	return res
}
