package parser

import (
	"fmt"
	"github.com/team-ide/go-interpreter/token"
	"sort"
)

const (
	errUnexpectedToken      = "Unexpected token %v"
	errUnexpectedEndOfInput = "Unexpected end of input"
)

func (this_ *Parser) Error(from string, place int, msg string) *Error {
	idx := place
	err := &Error{msg: "from:" + from + ",msg:" + msg, idx: idx}
	po := this_.GetPosition(place)
	if po != nil {
		err.line = po.Line
		err.column = po.Column
	}

	this_.Errors.Add(err)
	return (this_.Errors)[len(this_.Errors)-1]
}
func (this_ *Parser) ErrorUnexpected(from string, offset int, chr rune) error {
	if chr == -1 {
		return this_.Error(from, offset, errUnexpectedEndOfInput)
	}
	return this_.Error(from, offset, fmt.Sprintf(errUnexpectedToken, token.Illegal))
}

func (this_ *Parser) ErrorUnexpectedToken(from string, tkn token.Token) error {
	switch tkn {
	case token.Eof:
		return this_.Error(from, this_.Idx, errUnexpectedEndOfInput)
	}
	value := tkn.String()
	switch tkn {
	case token.Boolean, token.Null:
		value = this_.Literal
	case token.Identifier:
		return this_.Error(from, this_.Idx, "Unexpected identifier")
	case token.Keyword:
		// TODO Might be a future reserved word
		return this_.Error(from, this_.Idx, "Unexpected reserved word")
	case token.EscapedReservedWord:
		return this_.Error(from, this_.Idx, "Keyword must not contain escaped characters")
	case token.Number:
		return this_.Error(from, this_.Idx, "Unexpected number")
	case token.String:
		return this_.Error(from, this_.Idx, "Unexpected string")
	}
	return this_.Error(from, this_.Idx, fmt.Sprintf(errUnexpectedToken, value))
}

type Error struct {
	filename string
	idx      int
	line     int
	column   int
	msg      string
}

// FIXME Should this be "SyntaxError"?

func (this_ *Error) Error() string {
	filename := this_.filename
	if filename == "" {
		filename = "(anonymous)"
	}
	return fmt.Sprintf("%s: Idx:%d Line %d:%d %s",
		filename,
		this_.idx,
		this_.line,
		this_.column,
		this_.msg,
	)
}

type ErrorList []*Error

func (this_ *ErrorList) Add(err *Error) {
	*this_ = append(*this_, err)
}

func (this_ *ErrorList) Reset() { *this_ = (*this_)[0:0] }

func (this_ *ErrorList) Len() int { return len(*this_) }
func (this_ *ErrorList) Swap(i, j int) {
	x := (*this_)[i]
	y := (*this_)[j]
	(*this_)[i] = y
	(*this_)[j] = x
}
func (this_ *ErrorList) Less(i, j int) bool {
	x := (*this_)[i]
	y := (*this_)[j]
	if x.filename < y.filename {
		return true
	}
	if x.filename == y.filename {
		if x.line < y.line {
			return true
		}
		if x.line == y.line {
			return x.column < y.column
		}
	}
	return false
}

func (this_ *ErrorList) Sort() {
	sort.Sort(this_)
}

func (this_ *ErrorList) Error() string {
	var out = fmt.Sprintf("has %d errors", len(*this_))
	for inx, one := range *this_ {
		out += fmt.Sprintf("\nerror %d:%s", inx, one.Error())
	}
	return out
}

func (this_ *ErrorList) Err() error {
	if len(*this_) == 0 {
		return nil
	}
	return this_
}
