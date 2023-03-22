package token

// 符号 可能出现的所有符号定义
const (
	Plus      = "+"  // +
	Minus     = "-"  // -
	Multiply  = "*"  // *
	Exponent  = "**" // **
	Slash     = "/"  // /
	Remainder = "%"  // /

	And                = "&"   // &
	Or                 = "|"   // |
	ExclusiveOr        = "^"   // ^
	ShiftLeft          = "<<"  // <<
	ShiftRight         = ">>"  // >>
	UnsignedShiftRight = ">>>" // >>>

	AddAssign       = "+="  // +=
	SubtractAssign  = "-="  // -=
	MultiplyAssign  = "*="  // *=
	ExponentAssign  = "**=" // **=
	QuotientAssign  = "/="  // /=
	RemainderAssign = "%="  // %=

	AndAssign                = "&="   // &=
	OrAssign                 = "|="   // |=
	ExclusiveOrAssign        = "^="   // ^=
	ShiftLeftAssign          = "<<="  // <<=
	ShiftRightAssign         = ">>="  // >>=
	UnsignedShiftRightAssign = ">>>=" // >>>=

	LogicalAnd = "&&" // &&
	LogicalOr  = "||" // ||
	Coalesce   = "??" // ??
	Increment  = "++" // ++
	Decrement  = "--" // --

	Equal       = "=="  // ==
	StrictEqual = "===" // ===
	Less        = "<"   // <
	Greater     = ">"   // >
	Assign      = "="   // =
	Not         = "!"   // !

	BitwiseNot = "~" // ~

	NotEqual       = "!="  // !=
	StrictNotEqual = "!==" // !==
	LessOrEqual    = "<="  // <=
	GreaterOrEqual = ">="  // >=

	LeftParenthesis = "(" // (
	LeftBracket     = "[" // [
	LeftBrace       = "{" // {
	Comma           = "," // ,
	Period          = "." // .

	RightParenthesis = ")"   // )
	RightBracket     = "]"   // ]
	RightBrace       = "}"   // }
	Semicolon        = ";"   // ;
	Colon            = ":"   // :
	QuestionMark     = "?"   // ?
	QuestionDot      = "?."  // ?.
	Arrow            = "=>"  // =>
	Ellipsis         = "..." // ...
	Backtick         = "`"   // `
	ArrowRight       = "->"  // ->
	ArrowLeft        = "<-"  // <-
)
