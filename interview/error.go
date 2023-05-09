package main

import (
	"fmt"
)

type MyError struct {
}

// 实现 string 方法
func (m *MyError) String() {

}

// type error interface {
// 	Error() string
// }

func main() {
	var err error
	x, ok := err.(error)
	fmt.Println(x, ok) // <nil> false
}
