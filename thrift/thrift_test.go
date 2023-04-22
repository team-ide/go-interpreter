package thrift

import (
	"fmt"
	"github.com/team-ide/go-interpreter/parser"
	"os"
	"strings"
	"testing"
)

func TestThrift(t *testing.T) {
	testFileCode("code.txt")

	dir := `thrift`
	fs, _ := os.ReadDir(dir)
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name(), ".thrift") {
			continue
		}
		filename := dir + "/" + f.Name()
		testFileCode(filename)
	}
}

func testFileCode(filename string) {
	fmt.Println("testFileCode:", filename)
	bs, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	code := string(bs)

	err = testCode(filename, code)
	if err != nil {
		fmt.Println("error filename:", filename)
		panic(err)
	}
}
func testCode(filename string, code string) error {
	tree, err := Parse(filename, code)
	if tree != nil {
		parser.OutTree(code, tree)
	}
	return err
}
