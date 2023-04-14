package node

// Position 坐标
type Position struct {
	Filename string `json:"filename,omitempty"` // 坐标所在文件名
	Offset   int    `json:"offset"`             // 索引
	Line     int    `json:"line"`               // 所在函数
	Column   int    `json:"column"`             // 所在列
}
