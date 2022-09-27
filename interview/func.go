package main

// 1. Variable parameters（可变参数）

import (
	"fmt"
	"strings"
)

func greet(names ...string) {
	names[0] = "hello"
	names[1] = "world"
}

func main() {
	params := []string{"init", "value"}

	// 可变参数列表 <--> slice
	greet(params...)
	res := strings.Join(params, ",")

	fmt.Println(res) // hello,world
}
