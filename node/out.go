package node

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func OutTree(code string, tree *Tree) {
	for _, one := range tree.Children {
		bs, _ := json.Marshal(one)
		fmt.Println("tree one type:", reflect.TypeOf(one).String(), "start:", one.Start()-1, ",end:", one.End()-1, ",json:", string(bs))
		fmt.Println(code[one.Start()-1 : one.End()-1])
		outSub(code, 1, one)
		fmt.Println("--------------------------------------------")
	}
}

func outSub(code string, leven int, one interface{}) {

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
				outOne(code, fmt.Sprintf(fieldT.Name+"-%d", n), leven, iV.Interface())
			}
		default:
			if fieldV.Kind() == reflect.Ptr {
				if fieldV.IsNil() {
					continue
				}
			}
			outOne(code, fieldT.Name, leven, v)
		}
	}
}

func outOne(code string, name string, leven int, one interface{}) {
	if one == nil {
		return
	}
	var n Node = nil

	s, ok := one.(Statement)
	if ok && s != nil {
		n = s
	}
	e, ok := one.(Expression)
	if ok && e != nil {
		n = e
	}
	c, ok := one.(ClassElement)
	if ok && c != nil {
		n = c
	}
	if n != nil {
		var bef = ""
		for i := 0; i < leven; i++ {
			bef += "\t"
		}
		bs, _ := json.Marshal(one)
		fmt.Print(bef+"field:", name, ",type:", reflect.TypeOf(n).String())
		fmt.Println(",start:", n.Start()-1, ",end:", n.End()-1, ",json:", string(bs))
		str := code[n.Start()-1 : n.End()-1]
		str = strings.ReplaceAll(str, "\n", "\n"+bef)
		fmt.Println(bef + str)
		outSub(code, leven+1, one)
	}
}
