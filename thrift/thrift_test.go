package thrift

import (
	"fmt"
	"github.com/team-ide/go-interpreter/node"
	"testing"
)

const thriftCode = `
include "a.thrift"

namespace java com.thrift.service
namespace go com.thrift.service
namespace cpp com.thrift.service

/**
 * 异常信息
 */
exception XException{
	/**
     * 错误码
	 */
    1: optional i32 code;
	/**
     * 错误信息
	 */
    2: optional string msg;
}


/**
 * 对象
 */
struct Obj {
  1: i32 num1 = 0;		//默认值
  2: i32 num2;
  3: Operation op;		//可以嵌套其他类型
  4: optional string comment;	//可选字段
  5: list<i32> l		//list
  6: map<i32,string> m	//map
  7: set<string> s	//set
  8: bool vBool	//bool
  9: byte vByte	//byte
  10: i16 vI16	//i16
  11: i64 vI64	//i64
  11: double vDouble	//double
  12: string vString	//string
  13: Obj2 vObj2	//Obj2
}

struct Obj2 {
  1: i32 num1 = 0;		//默认值
}

/**
 * 枚举
 */
enum Operation {
  ADD = 1,
  SUBTRACT = 2,
  MULTIPLY = 3,
  DIVIDE = 4
}

/**
 * 响应
 */
struct XResponse{
	/**
     * 状态码
	 */
    1: optional i32 code;
	/**
     * 状态信息
	 */
    2: optional string msg;
}


/**
 * 服务
 */
service Service {
	/**
	 * 方法
	 * 参数：
	 *	    Obj request
	 * 返回：
	 *		XResponse res
	 */
	XResponse method1(1: Obj request);

}
`

func TestThrift(t *testing.T) {
	tree, err := Parse(thriftCode)
	if err != nil {
		fmt.Println(err.Error())
	}
	node.OutTree(thriftCode, tree)

}
