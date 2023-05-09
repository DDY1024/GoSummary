package main

import (
	"fmt"
	"runtime"
)

func main() {
	ch := make(chan int, 4)
	quit := make(chan bool)

	go func() {
		defer close(quit)
		for {
			select {
			case v, ok := <-ch:
				fmt.Println(v, ok)
				if !ok {
					return
				}
			default:
			}
		}
	}()

	// 缓冲 channel，直接写入，不阻塞
	for i := 0; i < 4; i++ {
		ch <- i
	}

	close(ch)
	<-quit
	fmt.Println("End")

	// GC 调优
	ballast := make([]byte, 10*1024*1024)
	runtime.KeepAlive(ballast)
}
