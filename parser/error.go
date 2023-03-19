package parser

import (
	"fmt"
	"sort"
)

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
	return fmt.Sprintf("%s: Line %d:%d %s",
		filename,
		this_.line,
		this_.column,
		this_.msg,
	)
}

func (this_ *parser) error(place int, msg string) *Error {
	idx := place

	this_.errors.Add(&Error{msg: msg, idx: idx})
	return (this_.errors)[len(this_.errors)-1]
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
