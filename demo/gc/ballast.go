package main

import "runtime"

// 原文参考：https://www.cnblogs.com/457220157-FTD/p/15567442.html
func main() {
	bs := make([]byte, 100<<20)
	// do something ...
	runtime.KeepAlive(bs)
}
