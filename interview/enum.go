package main

// 定义一种新的类型来进行枚举约束

type OrderStatus int

const (
	CREATE OrderStatus = iota + 1
	PAID
	DELIVERING
	COMPLETED
	CANCELLED
)
