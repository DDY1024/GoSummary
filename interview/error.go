package main

import (
	"fmt"
)

type MyError struct {
}

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
