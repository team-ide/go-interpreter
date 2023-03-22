package node

type Position struct {
	Filename string // 坐标所在文件名
	Idx      int    // 索引
	Line     int    // 所在函数
	Column   int    // 所在列
}

func (this_ *Position) NewByColumnOffset(offset int) *Position {
	newPosition := &Position{
		Filename: this_.Filename,
		Idx:      this_.Idx + offset,
		Line:     this_.Line,
		Column:   this_.Column + offset,
	}
	return newPosition
}
