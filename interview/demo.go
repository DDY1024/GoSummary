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

	for i := 0; i < 4; i++ {
		ch <- i
	}

	close(ch)
	<-quit
	fmt.Println("End")

	ballast := make([]byte, 10*1024*1024)
	runtime.KeepAlive(ballast)
}
