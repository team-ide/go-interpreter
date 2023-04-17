package node

import (
	"reflect"
	"unsafe"
)

// Node 节点
type Node interface {
	Start() int // 节点所在 开始位置
	End() int   // 节点所在 结束位置
}

func FromUtf16(b []uint16) string {
	var str string
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&str))
	hdr.Data = uintptr(unsafe.Pointer(&b[0]))
	hdr.Len = len(b) * 2

	return str
}
