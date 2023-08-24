package main

import "fmt"

// 1. const 块内声明的变量模式具有传递性
// 2. iota 在一个 const 块内从 0 开始递增，不同 const 块内均从 0 开始递增

const (
	a = iota // iota = 0
	b        // iota = 1, b = iota
	c        // iota = 2, c = iota
)

const (
	d   = 1    // iota = 0, d = 1
	e   = iota // iota = 1, e = iota
	f          // iota = 2, f = iota
	g   = 1    // iota = 3, g = 1
	h          // iota = 4, h = 1
	xxx = iota // iota = 5,
)

const (
	i = 1 << iota // iota = 0, i = 1<<iota
	j             // iota = 1, j = 1<<iota
	k             // iota = 2, k = 1<<iota

	z = iota // iota = 3, z = iota
	l        // iota = 4, l = iota
	m        // iota = 5, m = iota
	n        // iota = 6, n = iota

	o = 0 // iota = 7, o = 0
	p     // iota = 8, p = 0
	q     // iota = 9, q = 0
)

func main() {
	fmt.Println(xxx) // 5
}
