package main

import (
	"fmt"
	"log"
)

var (
	todo bool
)

// 声明变量的作用域问题

// func init() {
// 	todo, err := fn() // 1. todo 为局部变量，其作用域仅仅为 init 函数内

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("init:", todo) // 局部变量: true
// }

func init() {
	var err error
	todo, err = fn()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("init:", todo) // 全局变量: true
}

func fn() (bool, error) {
	return true, nil
}

func main() {
	fmt.Println("main:", todo) // 全局变量: true
}
