package main

import (
	"fmt"
)

// Go 内置类型声明
// type any = interface{}
// func panic(v any)
// func recover() any
// var nil Type
// type Type int

func main() {
	//fmt.Println(test())
	ret := test()
	fmt.Println("TT:", ret)
}

func test() int {
	defer func() interface{} {
		if err := recover(); err != nil {
			x, ok := err.(int)
			fmt.Println(x, ok) // 1, true
			// return err
			return x
		}
		return nil
	}()

	panic(1)
}

// func testTwo() (res int) {
// 	defer func() interface{} {
// 		if err := recover(); err != nil { //
// 			x, ok := err.(int)
// 			fmt.Println(x, ok) // 1, true
// 			res = x

// 			y, ok := err.(string)
// 			fmt.Println(y, ok)
// 			return err
// 		}
// 		return nil
// 	}()

// 	// panic(1)
// 	// panic(10)
// 	panic("woca") // panic 传入什么类型，recover 便返回什么类型
// }
