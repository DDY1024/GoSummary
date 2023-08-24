package main

import "fmt"

// 1. 函数
func MapKeys[K comparable, V any](data map[K]V) []K {
	r := make([]K, 0)
	for k := range data {
		r = append(r, k)
	}
	return r
}

// 2. 结构体
type List[T any] struct {
	head, tail *element[T]
}

type element[T any] struct {
	next *element[T]
	val  T
}

// 3. 方法
func (lst *List[T]) Push(v T) {
	if lst.tail == nil {
		lst.head = &element[T]{val: v}
		lst.tail = lst.head
	} else {
		lst.tail.next = &element[T]{val: v}
		lst.tail = lst.tail.next
	}
}

func (lst *List[T]) GetAll() []T {
	var elems []T
	for e := lst.head; e != nil; e = e.next {
		elems = append(elems, e.val)
	}
	return elems
}

func test() {
	m := map[string]int{
		"1": 1,
		"2": 2,
	}
	fmt.Println(MapKeys[string, int](m))

	lst := List[int]{}
	lst.Push(1)
	lst.Push(2)
	fmt.Println(lst.GetAll())
}
