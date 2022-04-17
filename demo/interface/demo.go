package main

// interface 完整性校验
var _ Gopher = (*WXY)(nil)

type Gopher interface {
	Print()
}

type WXY struct {
	Name string
}

func (wxy WXY) Print() {
}
