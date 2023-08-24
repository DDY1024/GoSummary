package memory

// 内存对齐
// 1. https://ms2008.github.io/2019/08/01/golang-memory-alignment/
// 2. https://eddycjy.gitbook.io/golang/di-1-ke-za-tan/go-memory-align

type A struct {
	a int32
	_ int32 // 补齐 padding
	b int64 //
}
