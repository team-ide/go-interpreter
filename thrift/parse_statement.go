package thrift

import (
	"github.com/apache/thrift/lib/go/thrift"
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
		inc = this_.ParsedLiteral
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
		res.ExtendsInclude, res.ExtendsName, res.To = this_.ParseIncludeName()
	}

	for this_.Token != token.LeftBrace && this_.Token != token.Eof {
		this_.Next()
	}

	for this_.Token != token.RightBrace && this_.Token != token.Eof {
		if this_.Token == token.LeftBrace || this_.Token == token.Semicolon || this_.Token == token.Comma {
			this_.Next()
			continue
		}
		field := this_.parseFieldNode()
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
		StructStatement: &StructStatement{
			From: idx,
		},
	}

	res.Name = this_.ParsedLiteral
	this_.Next()

	if this_.Token == token.Extends {
		this_.Next()
		res.ExtendsInclude, res.ExtendsName, res.To = this_.ParseIncludeName()
	}

	for this_.Token != token.LeftBrace && this_.Token != token.Eof {
		this_.Next()
	}

	for this_.Token != token.RightBrace && this_.Token != token.Eof {
		if this_.Token == token.LeftBrace || this_.Token == token.Semicolon || this_.Token == token.Comma {
			this_.Next()
			continue
		}
		field := this_.parseFieldNode()
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
		res.ExtendsInclude, res.ExtendsName, res.To = this_.ParseIncludeName()
	}

	for this_.Token != token.LeftBrace && this_.Token != token.Eof {
		this_.Next()
	}

	for this_.Token != token.RightBrace && this_.Token != token.Eof {
		if this_.Token == token.LeftBrace || this_.Token == token.Semicolon || this_.Token == token.Comma {
			this_.Next()
			continue
		}
		field := this_.parseEnumFieldNode()
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
		res.ExtendsInclude, res.ExtendsName, res.To = this_.ParseIncludeName()
	}

	for this_.Token != token.LeftBrace && this_.Token != token.Eof {
		this_.Next()
	}
	for this_.Token != token.RightBrace && this_.Token != token.Eof {
		if this_.Token == token.LeftBrace || this_.Token == token.Semicolon || this_.Token == token.Comma {
			this_.Next()
			continue
		}
		method := this_.parseServiceMethodNode()
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

func (this_ *Parser) parseFieldNode() *FieldNode {
	idx := this_.Idx

	res := &FieldNode{
		From: idx,
	}
	//fmt.Println("parseFieldDefinition token:", this_.Token, ",Literal:", this_.Literal, ",ParsedLiteral:", this_.ParsedLiteral)
	var num int
	// 表示有编号
	if this_.Literal != "" {
		var err error
		num, err = strconv.Atoi(this_.Literal)
		if err == nil {
			this_.Next()
			this_.ExpectAndNext("parseFieldDefinition", token.Colon)
		}
	}
	if this_.Token == token.Optional {
		res.Optional = true
		this_.Next()
	}
	res.Type = this_.parseFieldType()
	res.Num = int16(num)

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

func (this_ *Parser) parseServiceMethodNode() *ServiceMethodNode {
	idx := this_.Idx

	res := &ServiceMethodNode{
		From: idx,
	}
	//fmt.Println("parseIFaceDefinition token:", this_.Token)

	if this_.ParsedLiteral == "oneway" {
		res.Oneway = true
		this_.Next()
	}

	res.Return = this_.parseFieldType()

	res.Name = this_.ParsedLiteral

	this_.Next()

	for this_.Token != token.RightParenthesis && this_.Token != token.Eof {

		if this_.Token == token.LeftParenthesis || this_.Token == token.Semicolon || this_.Token == token.Comma {
			this_.Next()
			continue
		}

		field := this_.parseFieldNode()
		res.Params = append(res.Params, field)
	}

	res.To = this_.Idx
	if this_.Token == token.RightParenthesis {
		res.To = this_.ExpectAndNext("parseIFaceDefinition", token.RightParenthesis) + 1
	}
	if this_.ParsedLiteral == "throws" {
		this_.Next()
		for this_.Token != token.RightParenthesis && this_.Token != token.Eof {

			if this_.Token == token.LeftParenthesis || this_.Token == token.Semicolon || this_.Token == token.Comma {
				this_.Next()
				continue
			}

			field := this_.parseFieldNode()
			res.Exceptions = append(res.Exceptions, field)
		}
		if this_.Token == token.RightParenthesis {
			res.To = this_.ExpectAndNext("parseIFaceDefinition", token.RightParenthesis) + 1
		}

	}
	return res
}

func (this_ *Parser) parseEnumFieldNode() *FieldNode {
	idx := this_.Idx

	res := &FieldNode{
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

// parseFieldType 解析字段类型 如 xx、xx.x、xx<x>、xx<x,x>、xx.x<x,x>、xx.x<x,x.xx>、xx.x<x.xx,x.xx>
func (this_ *Parser) parseFieldType() *FieldType {

	from, typeName := this_.Idx, this_.ParsedLiteral

	res := &FieldType{
		From: from,
		To:   this_.Idx + len(typeName),
	}

	res.TypeName = typeName
	switch typeName {
	case "void":
		res.TypeId = thrift.VOID
	case "bool":
		res.TypeId = thrift.BOOL
	case "byte":
		res.TypeId = thrift.BYTE
	case "i8":
		res.TypeId = thrift.I08
	case "double":
		res.TypeId = thrift.DOUBLE
	case "i16":
		res.TypeId = thrift.I16
	case "i32":
		res.TypeId = thrift.I32
	case "i64":
		res.TypeId = thrift.I64
	case "string":
		res.TypeId = thrift.STRING
	case "utf7":
		res.TypeId = thrift.UTF7
	case "map":
		res.TypeId = thrift.MAP
	case "set":
		res.TypeId = thrift.SET
	case "list":
		res.TypeId = thrift.LIST
	case "utf8":
		res.TypeId = thrift.STRING
	case "utf16":
		res.TypeId = thrift.STRING
	case "uuid":
		//res.TypeId = thrift.UUID
		res.TypeId = thrift.STRING
	case "binary":
		res.TypeId = thrift.STRING
	default:
		res.TypeId = thrift.STRUCT
		res.StructInclude, res.StructName, res.To = this_.ParseIncludeName()

		res.TypeName = res.StructInclude + "." + res.StructName
		//if res.StructName == "" {
		//	panic(this_.Filename + " Eof err " + this_.ToJSON(this_.GetPosition(from)))
		//}

	}

	if res.TypeId != thrift.STRUCT {
		this_.OnlyReadGreater = true
		this_.OnlyReadLess = true
		this_.Next()
		this_.OnlyReadGreater = false
		this_.OnlyReadLess = false
	}

	var genericTypes []*FieldType
	if this_.Token == token.Less {
		this_.OnlyReadGreater = true
		this_.OnlyReadLess = true
		this_.Next()
		this_.OnlyReadGreater = false
		this_.OnlyReadLess = false
		for this_.Token != token.Greater && this_.Token != token.Eof {
			if this_.Token == token.Comma {
				this_.Next()
				continue
			}
			gType := this_.parseFieldType()
			if gType != nil {
				genericTypes = append(genericTypes, gType)
			}
		}
		//if this_.Token != token.Greater {
		//	panic(this_.Filename + " Eof err " + this_.ToJSON(this_.GetPosition(from)))
		//}
		this_.OnlyReadGreater = true
		res.To = this_.ExpectAndNext("parseFieldType", token.Greater)
		this_.OnlyReadGreater = false
	}

	if res.TypeId == thrift.LIST {
		res.ListType = genericTypes[0]
	} else if res.TypeId == thrift.SET {
		res.SetType = genericTypes[0]
	} else if res.TypeId == thrift.MAP {
		res.MapKeyType = genericTypes[0]
		res.MapValueType = genericTypes[1]
	}
	//fmt.Println("parseFieldType type:", this_.ToJSON(res))

	return res
}

func (this_ *Parser) ParseIncludeName() (include string, name string, to int) {

	name = this_.ParsedLiteral
	to = this_.Idx + len(name)
	this_.OnlyReadGreater = true
	this_.OnlyReadLess = true
	this_.Next()
	if this_.Token == token.Period {
		this_.Next()
		include = name
		name = this_.ParsedLiteral
		to = this_.Idx + len(name)
		this_.Next()
	}
	this_.OnlyReadGreater = false
	this_.OnlyReadLess = false

	return
}
