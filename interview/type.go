package main

type m struct{}

func (m *m) Lock()   {}
func (m *m) Unlock() {}

// 1. 编译错误
// type ptr *m
// type newM m

// 2. 正常运行
type ptr = *m
type newM = m

// var a ptr  // ptr 为一种新的类型
// var b newM // newM 为一种新的类型
// 新类型并没有声明其对应的 Lock 和 Unlock 方法

// which code section can compile success?
// func main() {
// 	// Why none of them can run ?
// 	a.Lock()
// 	b.Lock()
// }

// type hash []byte
// type xmap map[int]int
// type number float64
// type f func(int, int)

// var a hash
// var b []byte
// var i number
// var j float64

// func main() {

// 	// 内置类型 和 非内置类型 赋值转换问题
// 	// 1. 内置类型无法直接进行隐式类型转换
// 	// 2. 非内置类型均可以直接进行隐式类型转换

// 	// right here?
// 	a = b
// 	b = a

// 	// right here?
// 	a = make([]byte, 0)

// 	// right here?
// 	a = make(hash, 0)

// 	// right here?
// 	i = j
// 	j = i

// 	// righer here?
// 	i = number(j)
// 	j = float64(i)

// 	var c xmap
// 	var d map[int]int
// 	c = d
// 	d = c

// 	var e f
// 	var f func(int, int)
// 	e = f
// 	f = e
// }
