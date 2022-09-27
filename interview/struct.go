package main

import (
	"fmt"
	"time"
	"unsafe"
)

// 字节对齐相关
// 1. 对于任何类型的变量x：unsafe.Alignof(x)至少为 1
// 2. 对于struct 类型的变量x：unsafe.Alignof(x)是 （所有字段字节对齐的最大值）unsafe.Alignof(x.f)，但至少为 1
// 3. 对于数组类型的变量x ：unsafe.Alignof(x)与数组元素类型的变量的对齐方式相同

type A struct {
	Time time.Time
}

func main() {
	var a A
	var b struct{}
	fmt.Println(unsafe.Sizeof(a))          // 24
	fmt.Println(unsafe.Sizeof(b))          // 0
	fmt.Println(unsafe.Sizeof(struct{}{})) // 0
}
