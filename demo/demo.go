package main

import (
	"fmt"
	"sync"
	"time"
)

// 1 done，N wait
func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		time.Sleep(500 * time.Millisecond)
		wg.Done()
	}()

	go func() {
		wg.Wait()
		fmt.Println("done other")
	}()

	wg.Wait()
	fmt.Println("done main")
	time.Sleep(2 * time.Second)
}
