package main

import (
	"fmt"
)

func fn() (closed bool) {
	closed = true
	defer func() {
		fmt.Println("TT:", closed)
		closed = true // 2
	}()
	return false // 1
}

func main() {
	res := fn()
	fmt.Println("FF:", res)
}

/*
func main() {
	defer func() {
		fmt.Println("1")
	}()

	if true {
		fmt.Println("2")
		return
	}

	// return 之后的 defer 不能被注册
	defer func() {
		fmt.Println("3")
	}()

	// 输出结果: 2, 1
}
*/
