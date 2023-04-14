package thrift

import (
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/token"
	"strconv"
)

func (this_ *Parser) parseIncludeStatement() *IncludeStatement {

	idx := this_.ExpectAndNext("parseIncludeStatement", token.Include)

	res := &IncludeStatement{
		From: idx,
	}
	toIdx := this_.Idx

	inc := ""
	if this_.Token == token.String {
		inc = string(this_.ParsedLiteral)
		toIdx += len(this_.Literal)
		this_.ExpectAndNext("parseIncludeStatement", token.String)
	}
	res.Include = inc
	res.To = toIdx
	//fmt.Println("parseIncludeStatement ", res, ",Next token:", this_.Token)

	return res
}

func (this_ *Parser) parseNamespaceStatement() *NamespaceStatement {

	idx := this_.ExpectAndNext("parseNamespaceStatement", token.Namespace)

	res := &NamespaceStatement{
		From: idx,
	}

	if this_.Token == token.Identifier {
		identifier := this_.ParseIdentifier()
		res.Language = string(identifier.Name)
	}
	toIdx := this_.Idx
	namespace := ""
	for {
		if this_.Token == token.Period {
			namespace += "."
			this_.Next()
			continue
		} else if this_.ParsedLiteral != "" {
			namespace += string(this_.ParsedLiteral)
			this_.Next()
			if this_.Token != token.Period {
				break
			}
		} else {
			break
		}
	}
	res.Namespace = namespace
	res.To = toIdx + len(namespace)
	//fmt.Println("parseNamespaceStatement ", res, ",Next token:", this_.Token)

	return res
}

func (this_ *Parser) parseStructStatement() *StructStatement {

	idx := this_.ExpectAndNext("parseStructStatement", token.Struct)

	res := &StructStatement{
		From: idx,
	}

	if this_.Token == token.Identifier {
		identifier := this_.ParseIdentifier()
		res.Name = string(identifier.Name)
	}
	for this_.Token != token.RightBrace && this_.Token != token.Eof {
		if this_.Token == token.LeftBrace || this_.Token == token.Semicolon {
			this_.Next()
			continue
		}
		field := this_.parseFieldDefinition()
		res.Fields = append(res.Fields, field)
	}
	res.To = this_.Idx
	if this_.Token == token.RightBrace {
		res.To++
	}
	this_.ExpectAndNext("parseStructStatement", token.RightBrace)
	//fmt.Println("parseStructStatement ", res, ",Next token:", this_.Token)

	return res
}

func (this_ *Parser) parseExceptionStatement() *ExceptionStatement {

	idx := this_.ExpectAndNext("parseExceptionStatement", token.Exception)

	res := &ExceptionStatement{
		From: idx,
	}

	if this_.Token == token.Identifier {
		identifier := this_.ParseIdentifier()
		res.Name = string(identifier.Name)
	}
	for this_.Token != token.RightBrace && this_.Token != token.Eof {
		if this_.Token == token.LeftBrace || this_.Token == token.Semicolon {
			this_.Next()
			continue
		}
		field := this_.parseFieldDefinition()
		res.Fields = append(res.Fields, field)
	}
	res.To = this_.Idx
	if this_.Token == token.RightBrace {
		res.To++
	}
	this_.ExpectAndNext("parseExceptionStatement", token.RightBrace)
	//fmt.Println("parseExceptionStatement ", res, ",Next token:", this_.Token)

	return res
}

func (this_ *Parser) parseEnumStatement() *EnumStatement {

	idx := this_.ExpectAndNext("parseEnumStatement", token.Enum)

	res := &EnumStatement{
		From: idx,
	}

	if this_.Token == token.Identifier {
		identifier := this_.ParseIdentifier()
		res.Name = string(identifier.Name)
	}
	for this_.Token != token.RightBrace && this_.Token != token.Eof {
		if this_.Token == token.LeftBrace || this_.Token == token.Semicolon || this_.Token == token.Comma {
			this_.Next()
			continue
		}
		field := this_.parseEnumFieldDefinition()
		res.Fields = append(res.Fields, field)
	}
	res.To = this_.Idx
	if this_.Token == token.RightBrace {
		res.To++
	}
	this_.ExpectAndNext("parseEnumStatement", token.RightBrace)
	//fmt.Println("parseEnumStatement ", res, ",Next token:", this_.Token)
	return res
}

func (this_ *Parser) parseServiceStatement() *ServiceStatement {

	idx := this_.ExpectAndNext("parseServiceStatement", token.Service)

	res := &ServiceStatement{
		From: idx,
	}

	if this_.Token == token.Identifier {
		identifier := this_.ParseIdentifier()
		res.Name = string(identifier.Name)
	}
	for this_.Token != token.RightBrace && this_.Token != token.Eof {
		if this_.Token == token.LeftBrace || this_.Token == token.Semicolon || this_.Token == token.Comma {
			this_.Next()
			continue
		}
		method := this_.parseIFaceMethodDefinition()
		res.Methods = append(res.Methods, method)
	}
	res.To = this_.Idx
	if this_.Token == token.RightBrace {
		res.To++
	}
	this_.ExpectAndNext("parseServiceStatement", token.RightBrace)
	//fmt.Println("parseServiceStatement ", res, ",Next token:", this_.Token)
	return res
}

func (this_ *Parser) parseIFaceMethodDefinition() *IFaceMethodDefinition {
	idx := this_.Idx

	res := &IFaceMethodDefinition{
		From: idx,
	}
	//fmt.Println("parseIFaceDefinition token:", this_.Token)

	str, keyName, value, tkn := this_.parseFieldName()
	if str == "" && keyName == "" && tkn == "" {

	}
	res.Return = value

	if this_.Token == token.Less {
		this_.Next()
		str, keyName, value, tkn = this_.parseFieldName()
		//fmt.Println("parseFieldDefinition type ", ",num:", num, ",str:", str, ",keyName:", keyName, ",value:", value, ",tkn:", tkn)
		if this_.Token == token.Comma {
			this_.Next()
			str, keyName, value, tkn = this_.parseFieldName()
			//fmt.Println("parseFieldDefinition type ", ",num:", num, ",str:", str, ",keyName:", keyName, ",value:", value, ",tkn:", tkn)
		}
		if this_.Token == token.Greater {
			this_.Next()
		}
	}

	str, keyName, value, tkn = this_.parseFieldName()
	res.Name = value

	for this_.Token != token.RightParenthesis && this_.Token != token.Eof {

		if this_.Token == token.LeftParenthesis || this_.Token == token.Semicolon || this_.Token == token.Comma {
			this_.Next()
			continue
		}

		field := this_.parseFieldDefinition()
		res.Params = append(res.Params, field)
	}

	res.To = this_.Idx
	if this_.Token == token.RightParenthesis {
		res.To = this_.ExpectAndNext("parseIFaceDefinition", token.RightParenthesis) + 1
	}
	return res
}

func (this_ *Parser) parseFieldDefinition() *FieldDefinition {
	idx := this_.Idx

	res := &FieldDefinition{
		Idx: idx,
	}
	num := ""
	//fmt.Println("parseFieldDefinition token:", this_.Token)
	for {
		if this_.Token == token.Colon {
			this_.Next()
			break
		} else if this_.ParsedLiteral != "" {
			num += string(this_.ParsedLiteral)
			this_.Next()
		} else if this_.Literal != "" {
			num += this_.Literal
			this_.Next()
		} else {
			break
		}
	}
	if this_.Token == token.Optional {
		this_.Next()
	}
	str, keyName, value, tkn := this_.parseFieldName()
	if str == "" && keyName == "" && tkn == "" {

	}
	res.Type = value
	res.FieldNum, _ = strconv.Atoi(num)
	//fmt.Println("parseFieldDefinition type ", ",num:", num, ",str:", str, ",keyName:", keyName, ",value:", value, ",tkn:", tkn)
	if this_.Token == token.Less {
		this_.Next()
		str, keyName, value, tkn = this_.parseFieldName()
		//fmt.Println("parseFieldDefinition type ", ",num:", num, ",str:", str, ",keyName:", keyName, ",value:", value, ",tkn:", tkn)
		if this_.Token == token.Comma {
			this_.Next()
			str, keyName, value, tkn = this_.parseFieldName()
			//fmt.Println("parseFieldDefinition type ", ",num:", num, ",str:", str, ",keyName:", keyName, ",value:", value, ",tkn:", tkn)
		}
		if this_.Token == token.Greater {
			this_.Next()
		}
	}
	str, keyName, value, tkn = this_.parseFieldName()
	res.Key = value
	//fmt.Println("parseFieldDefinition name ", ",num:", num, ",str:", str, ",keyName:", keyName, ",value:", value, ",tkn:", tkn)
	if this_.Token == token.Assign {
		this_.Next()
		str, keyName, value, tkn = this_.parseFieldName()
		res.Value = value
		//fmt.Println("parseFieldDefinition value ", ",num:", num, ",str:", str, ",keyName:", keyName, ",value:", value, ",tkn:", tkn)
	}
	return res
}

func (this_ *Parser) parseEnumFieldDefinition() *FieldDefinition {
	idx := this_.Idx

	res := &FieldDefinition{
		Idx: idx,
	}
	str, keyName, value, tkn := this_.parseFieldName()
	if str == "" && keyName == "" && tkn == "" {

	}
	res.Key = value
	//fmt.Println("parseFieldDefinition name ", ",num:", num, ",str:", str, ",keyName:", keyName, ",value:", value, ",tkn:", tkn)
	if this_.Token == token.Assign {
		this_.Next()
		res.Value = &node.StringLiteral{
			Idx:     this_.Idx,
			Literal: this_.Literal,
			Value:   this_.ParsedLiteral,
		}
		//fmt.Println("parseFieldDefinition value ", ",num:", num, ",str:", str, ",keyName:", keyName, ",value:", value, ",tkn:", tkn)
	}
	return res
}

func (this_ *Parser) parseFieldName() (string, node.String, *node.StringLiteral, token.Token) {
	idx, tkn, literal, parsedLiteral := this_.Idx, this_.Token, this_.Literal, this_.ParsedLiteral
	this_.Next()
	if this_.Token == token.Period {
		//fmt.Println("parseFieldName Period")
		this_.Next()
		literal_, parsedLiteral_, _, _ := this_.parseFieldName()
		literal += "." + literal_
		parsedLiteral += "." + parsedLiteral_
	}
	value := &node.StringLiteral{
		Idx:     idx,
		Literal: literal,
		Value:   parsedLiteral,
	}
	return literal, parsedLiteral, value, tkn
}
