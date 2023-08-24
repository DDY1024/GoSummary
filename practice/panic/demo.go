package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

func main() {
	go func() {
		// 1. 方式一
		defer func() {
			if err := recover(); err != nil {
				const size = 64 << 10 // 64kb 堆栈输出
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
			}
		}()

		// 2. 方式二
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(string(debug.Stack()))
			}
		}()
	}()
}
