package parser

import (
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/token"
)

// Parse 解析
func (this_ *Parser) Parse() (tree *node.Tree, err error) {
	tree = this_.parseTree()
	//this_.errors.Sort()
	err = this_.Errors.Err()
	return
}

func (this_ *Parser) parseTree() (tree *node.Tree) {

	var statements []node.Statement
	this_.Read()
	this_.Next()
	for this_.Token != token.Eof {
		statements = append(statements, this_.ParseStatement())
	}

	tree = &node.Tree{
		Children:       statements,
		OffsetPosition: this_.OffsetPosition,
	}
	//this_.file.SetSourceMap(this_.parseSourceMap())
	return
}

func (this_ *Parser) Next() {
	this_.Token, this_.Literal, this_.ParsedLiteral, this_.Idx = this_.Scan()
}

func (this_ *Parser) IdxOf(offset int) int {
	return offset
}
func (this_ *Parser) Slice(start, end int) string {
	from := start
	to := end
	//if from >= 0 && to <= len(this_.str) {
	return this_.Str[from:to]
	//}

	//return ""
}

func (this_ *Parser) ParseChainNameStatement() *node.ChainNameStatement {

	res := &node.ChainNameStatement{
		From: this_.Idx,
	}

	for {
		if this_.Token == token.Period {
			if this_.ImplicitRead() == '.' {
				break
			}
			res.To = this_.Idx + 1
			this_.Next()
			continue
		} else if this_.ParsedLiteral != "" {
			res.Names = append(res.Names, this_.ParsedLiteral)
			res.To = this_.Idx + len(this_.ParsedLiteral)
			this_.Next()
			if this_.Token != token.Period {
				break
			}
		} else {
			break
		}
	}

	return res
}
