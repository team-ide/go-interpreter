package node

// Node 节点
type Node interface {
	Start() int // 节点所在 开始位置
	End() int   // 节点所在 结束位置
}
