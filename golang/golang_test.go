package golang

import (
	"github.com/team-ide/go-interpreter/parser"
	"os"
	"testing"
)

func TestGolang(t *testing.T) {
	bs, err := os.ReadFile("code.txt")
	if err != nil {
		panic(err)
	}
	code := string(bs)

	tree, err := Parse(code)
	if tree != nil {
		parser.OutTree(code, tree)
	}
	if err != nil {
		panic(err)
	}

}
