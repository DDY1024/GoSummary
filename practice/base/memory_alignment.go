package main

import (
	"fmt"
	"unsafe"
)

// 内存对齐
// 参考资料：https://www.liwenzhou.com/posts/Go/struct_memory_layout/

type Foo struct {
	A int8
	B int8
	C int8
}

// Go 编译时会按照【一定的规则】自动进行【内存对齐】
// 减少 CPU 访问内存的次数，加大 CPU 访问内存的吞吐量
type Bar struct {
	X int32  // 4
	Y *int32 // 8
	Z bool   // 1
}

// 通过优化结构体字段的顺序，减少结构体占用的空间大小
type Bar2 struct {
	X int32  // 4
	Z bool   // 1
	Y *int32 // 8
}

// 由于空结构体struct{}的大小为 0，所以当一个【结构体中】包含空结构体类型的字段时，通常不需要进行内存对齐
type Demo1 struct {
	M struct{} // 0
	N int8     // 1
}

// 【当空结构体类型作为结构体的最后一个字段时】，如果有指向该字段的指针，那么【就会返回该结构体之外的地址】；为了避免内存泄露会额外进行一次内存对齐
type Demo2 struct {
	N int8     // 1
	M struct{} // 1
}

// struct{} 作为结构体最后一个字段，存在特殊情况，需要额外进行一次内存对齐
type Demo3 struct {
	N int      // 8
	M struct{} // 8
}

// 1. 结构体热点字段放在第一个位置，避免计算偏移地址，提升效率
// 2. 结构体尾部添加 padding 字段，cache line 独占，避免伪共享
// 3. 强制结构体某个字段进行内存对齐，满足原子操作对于内存对齐的要求

func main() {
	var f Foo
	fmt.Println(unsafe.Sizeof(f))  // 3
	fmt.Println(unsafe.Alignof(f)) // 1

	var b Bar
	fmt.Println(unsafe.Sizeof(b))  // 24
	fmt.Println(unsafe.Alignof(b)) // 8

	// 注意 go 中 slice 和 array 的区别
	var arr1 = [3]int{1, 2, 3} // 8
	fmt.Println(unsafe.Alignof(arr1))

	var arr2 = [1]bool{true} // 1
	fmt.Println(unsafe.Alignof(arr2))

	var b2 Bar2
	fmt.Println(unsafe.Sizeof(b2)) // 16

	var d1 Demo1
	fmt.Println(unsafe.Sizeof(d1)) // 1

	var d2 Demo2
	fmt.Println(unsafe.Sizeof(d2)) // 2

	var d3 Demo3
	fmt.Println(unsafe.Sizeof(d3)) // 16
}
