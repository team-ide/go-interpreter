package node

type Node interface {
	Start() *Position // 节点所在 开始位置
	End() *Position   // 节点所在 结束位置
}
