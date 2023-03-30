package token

const (
	Keyword Token = "keyword" // 标记是关键字
	Nil     Token = "nil"     // 标记 nil
	Null    Token = "null"    // 标记 null

	Eof            Token = "eof"            // 标记 结束
	LineTerminator Token = "LineTerminator" // 标记 结束

	Illegal Token = "illegal" // 标记 非法的
	Empty   Token = "empty"   // 标记 是空的

	Identifier        Token = "identifier"        // 标记 是标识符
	PrivateIdentifier Token = "PrivateIdentifier" // 标记 是 私有 标识符

	EscapedReservedWord Token = "escapedReservedWord" // 标记 是 转义的保留字

	Async Token = "async" // 标记 同步
	Await Token = "await" // 标记 等待

	String Token = "string" // 标记 是 字符串
	Number Token = "number" // 标记 是 数字

	Of         Token = "of"         //
	Field      Token = "field"      // 表示 字段
	Method     Token = "method"     // 表示 方法
	BlankSpace Token = "blankSpace" // 表示 空白

)
