package main

// 建造者模式
// 1. 必填参数显示传入
// 2. 可选参数通过方法进行传入

type People struct {
	Age    int
	Name   string
	Gender int
}

type Option func(p *People)

func WithAge(age int) Option {
	return func(p *People) { p.Age = age }
}

func WithName(name string) Option {
	return func(p *People) { p.Name = name }
}

func WithGender(gender int) Option {
	return func(p *People) { p.Gender = gender }
}

func NewPeople(opts ...Option) *People {
	people := &People{}
	for _, opt := range opts {
		opt(people)
	}
	return people
}
