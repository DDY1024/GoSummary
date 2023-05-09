package main

import (
	"fmt"
)

// reflect.SliceHeader
type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

// 零值切片、nil 切片、空切片
func nilSlice() {
	var s1 []int
	fmt.Println(s1 == nil) // true

	s2 := new([]int)
	fmt.Println(s2 == nil) // false

	s3 := make([]int, 0)
	fmt.Println(s3 == nil) // false

	s4 := []int{}
	fmt.Println(s4 == nil) // false
}

func main() {
	nilSlice()
}
