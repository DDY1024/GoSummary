package unsafe

import (
	"fmt"
	"unsafe"
)

// 相关文章
// unsafe&&uintptr 坑总结: https://blog.csdn.net/u010853261/article/details/103826830
//
// uintptr 就是一个 16 进制的整数，这个数字表示对象的地址，但是 **uintptr 没有指针的语义**
// 一、如果一个对象只有一个 uintptr 表示的地址表示"引用"关系，那么这个对象会在 GC 时被无情的回收掉，那么 uintptr 将表示一个 "野地址"
// 二、如果 uintptr 表示的地址指向的对象发生了 copy 移动(比如协程栈增长，slice 扩容等)，那么 uintptr 也表示一个"野地址"
// 三、unsafe.Pointer 有指针语义，可以保护它所指向的对象在“有用”的时候不会被垃圾回收，并且在发生移动时候更新地址值

type Object struct {
	// addr uintptr
	addr unsafe.Pointer
}

var obj Object

// func doTestOne() {
// 	x := 2
// 	obj.addr = uintptr(unsafe.Pointer(&x))
// 	fmt.Println(obj.addr) // 824634167000
// 	// runtime.GC()
// }

func doTestTwo() {
	x := 2
	obj.addr = unsafe.Pointer(&x)
	fmt.Printf("%d\n", obj.addr)
}

// func main() {
// 	doTestOne()

// 	runtime.GC()
// 事实证明 uinptr 指向的地址的对象仍然被 GC 回收掉了，此时变成了一个野地址
// 	y := (*int)(unsafe.Pointer(obj.addr)) // uintptr --> unsafe.Pointer 是有问题的，容易读脏/写脏内存
// 	fmt.Printf("%d\n", y)                 // 824633802928
// 	fmt.Println(*y)                       // 824634834944
// }

// func main() {
// 	doTestTwo()

// 	runtime.GC()
// 	y := (*int)(obj.addr)
// 	fmt.Printf("%d\n", y) //
// 	fmt.Println(*y)       //
// 	/*
// 		824633802904
// 		824633802904
// 		2  // 没有被 gc，所以指向的对象内容不变
// 	*/
// }

// 案例分析一
/*
// 假设此函数不会被内联（inline）。
func createInt() *int {
	return new(int)
}

func foo() {
	p0, y, z := createInt(), createInt(), createInt()
	var p1 = unsafe.Pointer(y) // 和y一样引用着同一个值
	var p2 = uintptr(unsafe.Pointer(z))
	// 由于 p2 是 uintptr 类型，因此对 z 并不存在实际上的引用关系，z 随时可能被 GC 掉

	// 此时，即使z指针值所引用的int值的地址仍旧存储
	// 在p2值中，但是此int值已经不再被使用了，所以垃圾
	// 回收器认为可以回收它所占据的内存块了。另一方面，
	// p0和p1各自所引用的int值仍旧将在下面被使用。

	// uintptr值可以参与算术运算。
	p2 += 2; p2--; p2--

	*p0 = 1                         // okay
	*(*int)(p1) = 2                 // okay

	// 由于 z 随时可能被 GC 调，因此 p2 指向地址在 GC 过后可能会变成一个野地址，因此 p2 赋值操作是一个危险操作
	*(*int)(unsafe.Pointer(p2)) = 3 // 危险操作！操作野地址内容
}
*/

// 案例分析二
// 一个协程栈的大小改变时，开辟在此栈上的内存块需要移动，从而相应的值地址将改变
// uintptr 指向旧栈地址，针对旧栈地址的操作会引发问题

// 案例分析三
// runtime.KeepAlive：对象保活，防止被 GC
//
/*
func foo() {
	p0, y, z := createInt(), createInt(), createInt()
	var p1 = unsafe.Pointer(y)
	var p2 = uintptr(unsafe.Pointer(z)) // uintptr 赋值

	p2 += 2; p2--; p2--

	*p0 = 1
	*(*int)(p1) = 2
	*(*int)(unsafe.Pointer(p2))) = 3 // 转危为安！

	runtime.KeepAlive(z) // 确保z所引用的值仍在使用中，防止被 gc 掉 ? 那这看起来是一个编译期行为？
}
*/

// 案例分析四
// 一个值的可被使用范围可能并没有代码中看上去的大
/*
 p = uintptr(unsafe.Pointer(&t.y[0]))  // 虽然引用了 t.y[0]，但是 t.y[0] 地址随时可能被替换，重新针对该地址的操作会很危险 --> 野地址

	... // 使用t.x和t.y

	// 一个聪明的编译器能够觉察到值t.y将不会再被用到，
	// 所以认为t.y值所占的内存块可以被回收了。

	*(*byte)(unsafe.Pointer(p)) = 1 // 危险操作！

	println(t.x) // ok。继续使用值t，但只使用t.x字段。
}
*/

// unsafe.Pointer 是一个类型安全指针类型，可以将其当做是 void*

// 六种使用 unsafe.Pointer 模式
// https://golang.google.cn/pkg/unsafe/#Pointer
// https://golang.google.cn/pkg/unsafe/#Pointer
//
// type Pointer *ArbitraryType
//
// **注意**：任何时刻将一个 uintptr 类型转化成 unsafe.Pointer 指针类型都是一个危险操作，因为我们没有办法
// 确定 uintptr 指向地址的有效性
//
// 1. *T1 --> *T2
/*
func Float64bits(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}
*/

// 2. unsafe.Pointer --> uintptr

// 3. unsafe.Pointer --> uintptr --> 地址偏移 --> unsafe.Pointer
// 因为本身 unsafe.Pointer 是不能做地址运算的，而 uintptr 可以，因此需要在两者之间来回转换
/*
	p = unsafe.Pointer(uintptr(p) + offset)

	// 注意此处所有 uintptr 地址偏移计算后转 unsafe.Pointer 必须用一行表示，防止中间过程中
	// uintptr 指向的地址发生变化
	// equivalent to f := unsafe.Pointer(&s.f)
	// unsafe.Offset 获取 field 地址偏移
	f := unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + unsafe.Offsetof(s.f))
	// equivalent to e := unsafe.Pointer(&x[i])
	// unsafe.Sizeof 获取大小
	e := unsafe.Pointer(uintptr(unsafe.Pointer(&x[0])) + i*unsafe.Sizeof(x[0]))
*/

// 4. syscall
// syscall.Syscall(SYS_READ, uintptr(fd), uintptr(unsafe.Pointer(p)), uintptr(n))

// 5. reflect
// 6. reflect.StringHeader、reflect.SliceHeader
/*
var s string
hdr := (*reflect.StringHeader)(unsafe.Pointer(&s)) // case 1
hdr.Data = uintptr(unsafe.Pointer(p))              // case 6 (this case)
hdr.Len = n
*/

// unsafe.Pointer 最佳实践
// 具体参考：https://louyuting.blog.csdn.net/article/details/100178972
//
//
//
/*
type StringHeader struct {
    Data uintptr
    Len  int
}

type SliceHeader struct {
    Data uintptr
    Len  int
    Cap  int
}

// 理论上这种方式才是 "真正安全" 的 string 和 []byte 之间 zero-copy 实现
// 本质上是 reflect.SliceHeader 和 reflect.StringHeader 之间的相互转化
func stringtobyte(s string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: stringHeader.Data,
		Len:  stringHeader.Len,
		Cap:  stringHeader.Len,
	}
	runtime.KeepAlive(&s)  // 之所以使用 keepalive 是担心 uintptr 指向的内存地址在这期间发生变化
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func bytes2string(b []byte) string{
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	sh := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}
	runtime.KeepAlive(&b) // 之所以使用 keepalive 是担心 uintptr 指向的内存地址在这期间发生变化
	return *(*string)(unsafe.Pointer(&sh))
}
*/
