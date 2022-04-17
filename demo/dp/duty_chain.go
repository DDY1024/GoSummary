package main

// 职责链模式 --> 函数链
// 1. 敏感词过滤
// 2. 参数校验

// 职责模型: 敏感词过滤标准接口
type SensitiveWordFilter interface {
	Filter(content string) bool
}

// 职责链 --> 函数链 --> 存放各种敏感词过滤策略
type SensitiveWordFilterChain struct {
	filters []SensitiveWordFilter
}

func (c *SensitiveWordFilterChain) AddFilter(filter SensitiveWordFilter) {
	c.filters = append(c.filters, filter)
}

func (c *SensitiveWordFilterChain) Filter(content string) bool {
	for _, filter := range c.filters { // 遍历函数链校验内容是否满足条件
		if filter.Filter(content) {
			return true
		}
	}
	return false
}

// 具体职责功能实现
type AdSensitiveWordFilter struct{}

func (f *AdSensitiveWordFilter) Filter(content string) bool {
	return false
}

type PoliticalWordFilter struct{}

func (f *PoliticalWordFilter) Filter(content string) bool {
	return true
}

// 一种经典的职责链模式实现: http/rpc 框架中 middleware 实现（参考 gin 框架）
// middleware 两种执行模式
// 1. 直接按序遍历执行一遍
// 2. 前置处理（正向遍历），后置处理（逆向遍历）
/*
type Context struct {
	// handlers 包含一组执行函数
	// type HandlersChain []HandlerFunc
	handlers HandlersChain

	// 当前执行函数的索引位置
	index    int8
}

// Next 会按照顺序将一个个中间件执行完毕，并且 Next 也可以在中间件中进行调用，达到请求前以及请求后的处理
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}
*/
