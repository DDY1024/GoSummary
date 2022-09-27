package main

// 应用场景: 公共库初始化构建对象，解决入参比较多的问题
// option 模式：必须参数显示传入，可选参数 option 方式传入

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
