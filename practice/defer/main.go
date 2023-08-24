package main

import (
	"fmt"
)

// fn() = true
func fn() (closed bool) {
	defer func() {
		closed = true // 2. true
	}()
	return false // 1. false
}

func test() {
	defer func() {
		fmt.Println("1")
	}()

	if true {
		fmt.Println("2")
		return
	}

	// return 语句之后的 defer 不能被注册
	defer func() {
		fmt.Println("3")
	}()
}

func main() {
	test()
}
