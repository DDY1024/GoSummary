package main

// 建造者模式: 解决入参比较多的场景，必传参数显示传入，其余参数通过别的方式传入
// 应用场景: 公共库初始化构建对象

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
