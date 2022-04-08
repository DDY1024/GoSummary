package main

import (
	"fmt"
	"unsafe"
)

type demo3 struct {
	a struct{}
	b int32
}
type demo4 struct {
	b int32
	a struct{}
}

func main() {
	fmt.Println(unsafe.Sizeof(demo3{})) // 4
	fmt.Println(unsafe.Sizeof(demo4{})) // 8
}
