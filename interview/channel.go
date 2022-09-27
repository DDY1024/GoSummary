package main

import (
	"fmt"
)

// panic or print nil or print 10?
func main() {
	c := make(chan int, 1)
	c <- 10
	close(c)

	x, ok := <-c
	fmt.Println(x, ok) // 10, true

	x, ok = <-c
	fmt.Println(x, ok) // 0, false

	//runtime.SetFinalizer()
	//runtime.NumCPU()
	//runtime.GOMAXPROCS()
}
