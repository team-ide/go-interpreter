package parser

import (
	"encoding/json"
	"fmt"
	"github.com/team-ide/go-interpreter/node"
	"reflect"
	"strings"
)

func toJSON(obj interface{}) string {
	bs, _ := json.Marshal(obj)
	return string(bs)
}
func OutTree(code string, tree *node.Tree) {
	fmt.Println("-------------------out code tree start-------------------------")
	for _, one := range tree.Children {
		fmt.Println("-----------------------out code tree one start---------------------")
		fmt.Println("tree one type:", reflect.TypeOf(one).String(), "start:", toJSON(tree.GetPosition(one.Start())), ",end:", toJSON(tree.GetPosition(one.End())), ",json:", toJSON(one))
		fmt.Println(code[one.Start():one.End()])
		outSub(tree, code, 1, one)
		fmt.Println("-----------------------out code tree one end---------------------")
	}
	fmt.Println("-------------------out code tree end-------------------------")
}

func outSub(tree *node.Tree, code string, leven int, one interface{}) {

	// 获取结构体实例的反射类型对象
	oneVOf := reflect.ValueOf(one).Elem()
	oneTOf := reflect.TypeOf(one).Elem()
	// 遍历结构体所有成员
	for i := 0; i < oneVOf.NumField(); i++ {
		// 获取每个成员的结构体字段类型
		fieldV := oneVOf.Field(i)
		fieldT := oneTOf.Field(i)
		v := fieldV.Interface()
		if v == nil {
			continue
		}
		switch fieldV.Kind() {
		case reflect.Array, reflect.Slice:
			size := fieldV.Len()
			for n := 0; n < size; n++ {
				iV := fieldV.Index(n)
				outOne(tree, code, fmt.Sprintf(fieldT.Name+"-%d", n), leven, iV.Interface())
			}
		default:
			if fieldV.Kind() == reflect.Ptr {
				if fieldV.IsNil() {
					continue
				}
			}
			outOne(tree, code, fieldT.Name, leven, v)
		}
	}
}

func outOne(tree *node.Tree, code string, name string, leven int, one interface{}) {
	if one == nil {
		return
	}
	var n node.Node = nil

	c, ok := one.(node.Node)
	if ok && c != nil {
		n = c
	}
	if n != nil {
		var bef = ""
		for i := 0; i < leven; i++ {
			bef += "\t"
		}
		fmt.Print(bef+"field:", name, ",type:", reflect.TypeOf(n).String())
		fmt.Println(",start:", toJSON(tree.GetPosition(n.Start())), ",end:", toJSON(tree.GetPosition(n.End())), ",json:", toJSON(one))
		str := code[n.Start():n.End()]
		str = strings.ReplaceAll(str, "\r\n", "\n")
		str = strings.ReplaceAll(str, "\n", "\n"+bef)
		fmt.Println(bef + str)
		outSub(tree, code, leven+1, one)
	}
}
