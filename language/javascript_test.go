package language

import (
	"encoding/json"
	"fmt"
	"github.com/team-ide/go-interpreter/node"
	"github.com/team-ide/go-interpreter/parser"
	"testing"
)

func TestJavaScript(t *testing.T) {
	code := `

var a = 1;
var b = 1;
`
	tree, err := parser.Parse(code, &JavaScriptSyntax{})
	outTree(tree)
	if err != nil {
		panic("parser.Parse error:" + err.Error())
	}

}

func outTree(tree *node.Tree) {
	fmt.Println("tree:", tree)
	for _, one := range tree.Children {
		bs, _ := json.Marshal(one)
		fmt.Println("tree one:", string(bs))
	}
}
