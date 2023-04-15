package thrift

import (
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

	res.Language = this_.ParsedLiteral
	this_.Next()

	res.Namespace = this_.ParseChainNameStatement()
	res.To = res.Namespace.To
	//fmt.Println("parseNamespaceStatement ", res, ",Next token:", this_.Token)

	return res
}

func (this_ *Parser) parseStructStatement() *StructStatement {

	idx := this_.ExpectAndNext("parseStructStatement", token.Struct)

	res := &StructStatement{
		From: idx,
	}

	res.Name = this_.ParsedLiteral
	this_.Next()

	if this_.Token == token.Extends {
		this_.Next()
		res.Extends = this_.ParseChainNameStatement()
	}

	for this_.Token != token.LeftBrace && this_.Token != token.Eof {
		this_.Next()
	}

	for this_.Token != token.RightBrace && this_.Token != token.Eof {
		if this_.Token == token.LeftBrace || this_.Token == token.Semicolon || this_.Token == token.Comma {
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

	res.Name = this_.ParsedLiteral
	this_.Next()

	if this_.Token == token.Extends {
		this_.Next()
		res.Extends = this_.ParseChainNameStatement()
	}

	for this_.Token != token.LeftBrace && this_.Token != token.Eof {
		this_.Next()
	}

	for this_.Token != token.RightBrace && this_.Token != token.Eof {
		if this_.Token == token.LeftBrace || this_.Token == token.Semicolon || this_.Token == token.Comma {
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

	res.Name = this_.ParsedLiteral
	this_.Next()

	if this_.Token == token.Extends {
		this_.Next()
		res.Extends = this_.ParseChainNameStatement()
	}

	for this_.Token != token.LeftBrace && this_.Token != token.Eof {
		this_.Next()
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

	res.Name = this_.ParsedLiteral
	this_.Next()

	if this_.Token == token.Extends {
		this_.Next()
		res.Extends = this_.ParseChainNameStatement()
	}

	for this_.Token != token.LeftBrace && this_.Token != token.Eof {
		this_.Next()
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

func (this_ *Parser) parseFieldDefinition() *FieldDefinition {
	idx := this_.Idx

	res := &FieldDefinition{
		From: idx,
	}
	num := ""
	//fmt.Println("parseFieldDefinition token:", this_.Token, ",Literal:", this_.Literal, ",ParsedLiteral:", this_.ParsedLiteral)
	// 表示有编号
	if this_.Literal != "" {
		var err error
		res.Num, err = strconv.Atoi(this_.Literal)
		if err == nil {
			this_.Next()
			this_.ExpectAndNext("parseFieldDefinition", token.Colon)
		}
	}
	if this_.Token == token.Optional {
		this_.Next()
	}
	res.Type = this_.parseFieldType()
	res.Num, _ = strconv.Atoi(num)

	res.Name = this_.ParsedLiteral
	res.To = this_.Idx + len(this_.ParsedLiteral)
	this_.Next()

	if this_.Token == token.Assign {
		this_.Next()
		if this_.ParsedLiteral == "" {
			res.Value = this_.Literal
			res.To = this_.Idx + len(this_.Literal)
		} else {
			res.Value = this_.ParsedLiteral
			res.To = this_.Idx + len(this_.ParsedLiteral)
		}
		this_.Next()
	}
	return res
}

func (this_ *Parser) parseIFaceMethodDefinition() *IFaceMethodDefinition {
	idx := this_.Idx

	res := &IFaceMethodDefinition{
		From: idx,
	}
	//fmt.Println("parseIFaceDefinition token:", this_.Token)

	res.Return = this_.parseFieldType()

	res.Name = this_.ParsedLiteral
	this_.Next()

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

func (this_ *Parser) parseEnumFieldDefinition() *FieldDefinition {
	idx := this_.Idx

	res := &FieldDefinition{
		From: idx,
	}
	res.Name = this_.ParsedLiteral
	res.To = this_.Idx + len(this_.ParsedLiteral)
	this_.Next()
	if this_.Token == token.Assign {
		this_.Next()
		if this_.ParsedLiteral == "" {
			res.Value = this_.Literal
			res.To = this_.Idx + len(this_.Literal)
		} else {
			res.Value = this_.ParsedLiteral
			res.To = this_.Idx + len(this_.ParsedLiteral)
		}
		this_.Next()
	}
	return res
}

func (this_ *Parser) parseFieldName() string {
	parsedLiteral := this_.ParsedLiteral
	this_.Next()
	return parsedLiteral
}

// parseFieldType 解析字段类型 如 xx、xx.x、xx<x>、xx<x,x>、xx.x<x,x>、xx.x<x,x.xx>、xx.x<x.xx,x.xx>
func (this_ *Parser) parseFieldType() FieldType {
	from, typeName := this_.Idx, this_.ParsedLiteral
	//fmt.Println("parseFieldType", ",token:", this_.Token, ",position:", this_.GetPosition(this_.Idx), ",typeName:", typeName)
	this_.OnlyReadGreater = true
	this_.OnlyReadLess = true
	this_.Next()
	this_.OnlyReadGreater = false
	this_.OnlyReadLess = false

	var res FieldType
	var genericTypes *[]FieldType
	var to *int
	if this_.Token == token.Period {
		//fmt.Println("parseFieldName Period")

		fieldTypeDot := &FieldTypeDot{
			From:  from,
			To:    from + len(typeName),
			Names: []string{typeName},
		}
		for this_.Token == token.Period {
			this_.Next()
			fieldTypeDot.To = this_.Idx + len(this_.ParsedLiteral)
			fieldTypeDot.Names = append(fieldTypeDot.Names, this_.ParsedLiteral)

			this_.OnlyReadGreater = true
			this_.OnlyReadLess = true
			this_.Next()
			this_.OnlyReadGreater = false
			this_.OnlyReadLess = false

		}
		to = &fieldTypeDot.To
		genericTypes = &fieldTypeDot.GenericTypes
		res = fieldTypeDot
	} else {
		fieldTypeName := &FieldTypeName{
			From: from,
			To:   from + len(typeName),
			Name: typeName,
		}
		to = &fieldTypeName.To
		genericTypes = &fieldTypeName.GenericTypes
		res = fieldTypeName
	}
	if this_.Token == token.Less {
		this_.Next()
		for this_.Token != token.Greater {
			if this_.Token == token.Comma {
				this_.Next()
				continue
			}
			gType := this_.parseFieldType()
			*genericTypes = append(*genericTypes, gType)
		}
		this_.OnlyReadGreater = true
		*to = this_.ExpectAndNext("parseFieldType", token.Greater)
		this_.OnlyReadGreater = false
	}
	return res
}
