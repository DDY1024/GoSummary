package main

import (
	"fmt"
	"sync"
)

// Go 1.18 默认新增了 TryLock 方法
// 参考：https://mp.weixin.qq.com/s/nS-72MLogNmwUBcvC2Xq6g

func main() {
	var mux sync.Mutex
	if mux.TryLock() {
		fmt.Println("yes")
	}
	mux.Unlock()

	var rmux sync.RWMutex
	rmux.TryLock()
	rmux.TryRLock()
}
