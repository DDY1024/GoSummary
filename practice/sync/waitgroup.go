package main

import (
	"fmt"
)

// func main() {
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	go func() {
// 		wg.Wait() // 可以多次 wait
// 		fmt.Println("G1 end")
// 	}()
// 	go func() {
// 		wg.Wait()
// 		fmt.Println("G2 end")
// 	}()

// 	time.Sleep(2 * time.Second)
// 	wg.Done()
// 	time.Sleep(2 * time.Second)
// }

// func main() {
// 	wg := sync.WaitGroup{}
// 	wg.Add(100)
// 	for i := 0; i < 100; i++ {
// 		go func() {
// 			defer wg.Done()
// 			// time.Sleep(10 * time.Millisecond)
// 			fmt.Println(i)
// 		}()
// 	}
// 	wg.Wait()
// }

func main() {
	var a []int
	var b []int
	fmt.Println(a, b)

	c := new([]int)
	d := new([]int)
	fmt.Println(c, d)

	// reflect.SliceHeader
	copy()
}
