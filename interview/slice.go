package main

import (
	"fmt"
)

// func main() {
// 	var a []int
// 	fmt.Println(a == nil) // true, a is nil

// 	b := a[:]
// 	fmt.Println(b == nil) // true, b is nil

// 	// SliceHeader 已经被创建
// 	a = make([]int, 0)
// 	// a = make([]int, 5)
// 	fmt.Println(a == nil) // false, a is not nil

// 	b = a[:]
// 	fmt.Println(b == nil) // false, b is not nil
// }

// go 中任何参数均为传值
func op(arr []int) {
	arr = append(arr, 1, 2, 3)
	fmt.Println("T2:", arr) // [1, 2, 3, 1, 2, 3]
}

func main() {
	arr := []int{1, 2, 3}
	op(arr)                 // append 操作导致 slice 底层数据数组 data 发生了改变，外部变量无法捕获到这一变化
	fmt.Println("T1:", arr) // [1,2,3]
}
