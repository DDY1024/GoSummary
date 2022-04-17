package main

// 装饰器模式
// 1. 解决继承关系过于复杂的问题，通过组合替代继承，在原始类基础上提供各类扩展能力
// 2. "基类" 只实现最基础的功能，"子类" 嵌套基类实现自定义的功能扩展

// 装饰器模式 vs 模板模式
// 1. 模板模式只会基于一套标准的 interface，自定义实现各类方法，并不会进行任何扩展
// 2. 装饰器模式，在基类基础上还可以继续扩展 **成员、方法** 等

type IDraw interface {
	Draw() string
}

type Square struct{}

func (s Square) Draw() string {
	return "this is a square"
}

// 扩展类型
type ColorSquare struct {
	square IDraw
	color  string
}

func NewColorSquare(square IDraw, color string) ColorSquare {
	return ColorSquare{color: color, square: square}
}

func (c ColorSquare) Draw() string {
	return c.square.Draw() + ", color is " + c.color
}
