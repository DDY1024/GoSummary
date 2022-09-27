package main

import (
	"fmt"
	"sync"
)

func main() {
	var once sync.Once
	once.Do(func() {
		fmt.Println("very good!")
	})
}
