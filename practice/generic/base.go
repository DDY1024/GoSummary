package main

import (
	"fmt"
)

// 1. 类型形参
// 2. 类型实参
// 3. 类型形参列表
// 4. 类型约束
// 5. 实例化
// 6. 泛型类型
// 7. 泛型接收器
// 8. 泛型函数

type Slice[T int | float32 | float64] []T

// 泛型 receiver
// 不支持泛型方法
func (s Slice[T]) Sum() T {
	var sum T
	for _, value := range s {
		sum += value
	}
	return sum
}

// interface{} 包裹类型形参，避免语义上的歧义
type NewSlice[T interface{ *int | *float64 }] []T

type MyMap[Key int | string, Value float32 | float64] map[Key]Value

type Node[T int | string] struct {
	Name string
	Data T
}

type Printer[T int | float64 | string] interface {
	Print(data T)
}

type MyChan[T int | string] chan T

// 类型形参互相嵌套使用
// 类型形参 S 套用类型形参 T
type WordStruct[T int | float64 | string, S []T] struct {
	Data  S
	Value T
}

// 泛型类型嵌套
type WoWMap[T int | float32 | float64] map[T]Slice[T]

// 注意: 匿名结构体不支持泛型
// 注意：type switch 针对类型形参是不被允许的，但是可以通过 reflect 机制获取其对应的类型

// 泛型函数
func Add[T int | float32 | float64 | string](a, b T) T {
	return a + b
}

// 类型支持算术运算
func Sub[T int](a, b T) T {
	return a - b
}

// 注意：匿名函数不支持泛型

// Go 支持泛型函数，但是并不支持泛型方法；可以通过泛型 receiver 达到泛型方法的效果

// 泛型 & interface
// 接口中声明具体的类型约束
type Int interface {
	int | int8 | int16 | int32 | int64

	// ~ 表示所有以 int 为底层类型的所有类型
	// 注意：
	// ~ 后面的类型不能是 interface{}，必须为基本类型
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Uint interface {
	uint | uint8 | uint16 | uint32
}

type Float interface {
	float32 | float64
}

type SliceElement interface {
	Int | Uint | Float
}

type XSlice[T SliceElement] []T
type XXSlice[T Int | Uint | Float] []T

type MyInt int

type XXXSlice[T ~int] []T

// 接口类型
type AllInt interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32
}

type XUint interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// 接口 A 表示 AllInt 和 XUint 类型的【交集】
// 类型并集、类型交集
// 空集 vs 空接口
// 空集：类型的交集为空集
// 空接口：interface{} 代表所有类型
type A interface {
	AllInt
	XUint
}

// type any = interface{}
type XXXXSlice[T interface{}] []T
type XXXXXSlice[T any] []T

// 可比较: comparable
// 可排序: golang.org/x/exp/constraints；并没有内置关键字实现 ordered
type XXMap[T comparable, S any] map[T]S

// 基本接口、一般接口
// 基本接口：接口定义中只有方法集
// 一般接口：类型约束 + 方法集；【一般接口类型不能用来定义变量，只能用于泛型的类型约束中】
type DataProcessor[T any] interface {
	Process(oriData T) (newData T)
	Save(data T) error
}

type DataProcessor2[T any] interface {
	int | ~struct{ Data interface{} }

	Process(data T) (newData T)
	Save(data T) error
}

var xxx = func() int {
	fmt.Println("Test:", 10086) // 并没有输出
	return 0
}

func getKeys[K comparable, V any](m map[K]V) []K {
	res := []K{}
	for k := range m {
		res = append(res, k)
	}
	return res
}

func baseTest() {
	var a Slice[int] = []int{1, 2, 3}
	fmt.Println(a)
	fmt.Printf("%T\n", a)
	fmt.Println("Slice Sum:", a.Sum())

	var b Slice[float64] = []float64{1.0, 2.0, 3.0}
	fmt.Println(b)

	var c MyMap[string, float64] = map[string]float64{
		"a": 1.0,
		"b": 2.0,
	}
	fmt.Println(c)

	// queue
	var q1 Queue[int]
	// q1 := Queue[int]{}
	// q1 := new(Queue[int])
	q1.Put(1)
	q1.Put(2)
	fmt.Println(q1.Size())

	fmt.Println(Add[int](1, 2))
	fmt.Println(Add[float32](1.0, 2.0))
	fmt.Println(Add[string]("a", "b"))

	var s1 XXXSlice[MyInt] = []MyInt{1, 2, 3}
	fmt.Println("XXX:", s1)
}

func main() {
	// ss := []string{"a", "b", "c"}
	// for _, v := range ss { // 值复制
	// 	fmt.Println(v)
	// }
	// for i := range ss {
	// 	fmt.Println(i)
	// }
	// // fmt.Println()

	// arr := []int{1, 2, 3}
	// for range arr { // 表达式只会在初始被求值一次，因此最终结果为 [1, 2, 3, 10, 10, 10]
	// 	arr = append(arr, 10)
	// }
	// fmt.Println(arr)

	ch1 := make(chan int, 3) // ❶
	go func() {
		ch1 <- 0
		ch1 <- 1
		ch1 <- 2
		close(ch1)
	}()

	ch2 := make(chan int, 3) // ❷
	go func() {
		ch2 <- 10
		ch2 <- 11
		ch2 <- 12
		close(ch2)
	}()

	ch := ch1           // ❸
	for v := range ch { // ❹ 切记表达式只会被求值一次
		fmt.Println(v)
		ch = ch2 // ❺
	}
	// 0, 1, 2

	a := [3]int{0, 1, 2}  // ❶
	for i, v := range a { // ❷  值复制
		a[2] = 10   // ❸
		if i == 2 { // ❹
			fmt.Println(v) // 2
		}
	}

	// 只会起一个变量来进行复制

loop: // ❶
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", i)

		switch i {
		default:
		case 2:
			break loop // ❷
			// 标签 break
		}
	}
}
