package token

const (
	Keyword = "keyword" // 标记是关键字
	Nil     = "nil"     // 标记 nil
	Null    = "null"    // 标记 null

	Eof = "eof" // 标记 结束

	Illegal = "illegal" // 标记 非法的
	Empty   = "empty"   // 标记 是空的

	Identifier        = "identifier"        // 标记 是标识符
	PrivateIdentifier = "PrivateIdentifier" // 标记 是 私有 标识符

	EscapedReservedWord = "escapedReservedWord" // 标记 是 转义的保留字

	Async = "async" // 标记 同步
	Await = "await" // 标记 等待

	String = "string" // 标记 是 字符串
	Number = "number" // 标记 是 数字
)
