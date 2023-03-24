package token

// 符号 可能出现的所有符号定义
const (
	Plus      Token = "+"  // +
	Minus     Token = "-"  // -
	Multiply  Token = "*"  // *
	Exponent  Token = "**" // **
	Slash     Token = "/"  // /
	Remainder Token = "%"  // /

	And                Token = "&"   // &
	Or                 Token = "|"   // |
	ExclusiveOr        Token = "^"   // ^
	ShiftLeft          Token = "<<"  // <<
	ShiftRight         Token = ">>"  // >>
	UnsignedShiftRight Token = ">>>" // >>>

	AddAssign       Token = "+ Token ="  // + Token =
	SubtractAssign  Token = "- Token ="  // - Token =
	MultiplyAssign  Token = "* Token ="  // * Token =
	ExponentAssign  Token = "** Token =" // ** Token =
	QuotientAssign  Token = "/ Token ="  // / Token =
	RemainderAssign Token = "% Token ="  // % Token =

	AndAssign                Token = "& Token ="   // & Token =
	OrAssign                 Token = "| Token ="   // | Token =
	ExclusiveOrAssign        Token = "^ Token ="   // ^ Token =
	ShiftLeftAssign          Token = "<< Token ="  // << Token =
	ShiftRightAssign         Token = ">> Token ="  // >> Token =
	UnsignedShiftRightAssign Token = ">>> Token =" // >>> Token =

	LogicalAnd Token = "&&" // &&
	LogicalOr  Token = "||" // ||
	Coalesce   Token = "??" // ??
	Increment  Token = "++" // ++
	Decrement  Token = "--" // --

	Equal       Token = " Token = Token ="         //  Token = Token =
	StrictEqual Token = " Token = Token = Token =" //  Token = Token = Token =
	Less        Token = "<"                        // <
	Greater     Token = ">"                        // >
	Assign      Token = " Token ="                 //  Token =
	Not         Token = "!"                        // !

	BitwiseNot Token = "~" // ~

	NotEqual       Token = "! Token ="         // ! Token =
	StrictNotEqual Token = "! Token = Token =" // ! Token = Token =
	LessOrEqual    Token = "< Token ="         // < Token =
	GreaterOrEqual Token = "> Token ="         // > Token =

	LeftParenthesis Token = "(" // (
	LeftBracket     Token = "[" // [
	LeftBrace       Token = "{" // {
	Comma           Token = "," // ,
	Period          Token = "." // .

	RightParenthesis Token = ")"         // )
	RightBracket     Token = "]"         // ]
	RightBrace       Token = "}"         // }
	Semicolon        Token = ";"         // ;
	Colon            Token = ":"         // :
	QuestionMark     Token = "?"         // ?
	QuestionDot      Token = "?."        // ?.
	Arrow            Token = " Token =>" //  Token =>
	Ellipsis         Token = "..."       // ...
	Backtick         Token = "`"         // `
	ArrowRight       Token = "->"        // ->
	ArrowLeft        Token = "<-"        // <-
)
