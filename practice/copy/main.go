package main

import (
	"fmt"

	"github.com/jinzhu/copier"
)

// reflect --> deep copy

type Node struct {
	Name string
	Age  int
	Next *Node
}

func main() {
	nd1 := &Node{"w", 18, nil}
	nd2 := &Node{"x", 18, nd1}
	// copier.Copy()
	var nd3 Node
	copier.Copy(&nd3, nd2)
	fmt.Println(nd2)
	fmt.Println(nd3)
}
