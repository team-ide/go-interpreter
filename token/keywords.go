package token

// 关键字 可能出现的所有关键字定义
const (

	// Golang中出现的关键字

	Break       Token = "break"       // break
	Default     Token = "default"     // default
	Interface   Token = "interface"   // interface
	Func        Token = "func"        // func
	Select      Token = "select"      // select
	Case        Token = "case"        // case
	Defer       Token = "defer"       // defer
	Go          Token = "go"          // go
	Map         Token = "map"         // map
	Struct      Token = "struct"      // struct
	Chan        Token = "chan"        // chan
	Else        Token = "else"        // else
	Goto        Token = "goto"        // goto
	Package     Token = "package"     // package
	Switch      Token = "switch"      // switch
	Const       Token = "const"       // const
	Fallthrough Token = "fallthrough" // fallthrough
	If          Token = "if"          // if
	Range       Token = "range"       // range
	Type        Token = "type"        // type
	Continue    Token = "continue"    // continue
	For         Token = "for"         // for
	Import      Token = "import"      // import
	Return      Token = "return"      // return
	Var         Token = "var"         // var

	// Java中出现的补充

	Abstract Token = "abstract" // abstract
	Assert   Token = "assert"   // assert
	Boolean  Token = "boolean"  // boolean

	//Break         Token = "break"        // break

	Byte Token = "byte" // byte

	//Case          Token = "case"         // case

	Catch Token = "catch" // catch
	Char  Token = "char"  // char
	Class Token = "class" // class

	//Const         Token = "const"        // const
	//Continue      Token = "continue"     // continue
	//Default       Token = "default"      // default

	Do     Token = "do"     // do
	Double Token = "double" // double

	//Else          Token = "else"         // else

	Enum    Token = "enum"    // enum
	Extends Token = "extends" // extends
	Final   Token = "final"   // final
	Finally Token = "finally" // finally
	Float   Token = "float"   // float

	//For           Token = "for"          // for
	//Goto          Token = "goto"         // goto
	//If            Token = "if"           // if

	Implements Token = "implements" // implements

	//Import        Token = "import"       // import

	Instanceof Token = "instanceof" // instanceof
	Int        Token = "int"        // int

	//Interface     Token = "interface"    // interface

	Long   Token = "long"   // long
	Native Token = "native" // native
	New    Token = "new"    // new

	//Package       Token = "package"      // package

	Private   Token = "private"   // private
	Protected Token = "protected" // protected
	Public    Token = "public"    // public

	//Return        Token = "return"       // return

	Short    Token = "short"    // short
	Static   Token = "static"   // static
	Strictfp Token = "strictfp" // strictfp
	Super    Token = "super"    // super

	//Switch        Token = "switch"       // switch

	Synchronized Token = "synchronized" // synchronized
	This         Token = "this"         // this
	Throw        Token = "throw"        // throw
	Throws       Token = "throws"       // throws
	Transient    Token = "transient"    // transient
	Try          Token = "try"          // try
	Void         Token = "void"         // void
	Volatile     Token = "volatile"     // volatile
	While        Token = "while"        // while

	// JavaScript中出现的补充
	//Break         Token = "break" // break
	//Case  Token = "case"          // case
	//Catch  Token = "catch" // catch
	//Class  Token = "class" // class
	//Const  Token = "const"       // const
	//Continue  Token = "continue" // continue

	Debugger Token = "debugger" // debugger

	//Default  Token = "default" // default

	Delete Token = "delete" // delete

	//Do  Token = "do"         // do
	//Else  Token = "else" // else

	Export Token = "export" // export

	//Extends     Token = "extends"    // extends
	//Finally     Token = "finally"    // finally
	//For         Token = "for"        // for

	Function Token = "function" // function

	//If          Token = "if"         // if
	//Import      Token = "import"     // import

	In Token = "in" // in

	//Instanceof  Token = "instanceof" // instanceof
	//New         Token = "new"        // new
	//Return      Token = "return"     // return
	//Super       Token = "super"      // super
	//Switch      Token = "switch"     // switch
	//This        Token = "this"       // this
	//Throw       Token = "throw"      // throw
	//Try         Token = "try"        // try

	Typeof Token = "typeof" // typeof

	//Var         Token = "var"        // var
	//Void        Token = "void"       // void
	//While       Token = "while"      // while

	With Token = "with" // with

	// JavaScript严格模式补充

	//Implements  Token = "implements" // implements
	//Interface   Token = "interface"  // interface

	Let Token = "let" // let

	//Package     Token = "package"    // package
	//Private     Token = "private"    // private
	//Protected   Token = "protected"  // protected
	//Public      Token = "public"     // public
	//Static      Token = "static"     // static

	Yield Token = "yield" // yield

	// Thrift中出现的补充

	Include   Token = "include"   // include
	Namespace Token = "namespace" // namespace
	Exception Token = "exception" // exception
	Service   Token = "service"   // service

)
