package main

import (
	"fmt"
	"sync"
)

// 参考文章: https://jiajunhuang.com/articles/2018_11_12-golang_nocopy.md.html

type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) UnLock() {}

type Demo struct {
	// _ noCopy
	_ sync.Mutex
}

func TestFunc(d Demo) {

}

func main() {
	d := Demo{}
	fmt.Printf("%+v", d)

	TestFunc(d)

	fmt.Printf("%+v", d)
}

/*
➜  lock git:(main) ✗ go vet nocopy.go
# command-line-arguments
./nocopy.go:21:17: TestFunc passes lock by value: command-line-arguments.Demo contains sync.Mutex
./nocopy.go:27:20: call of fmt.Printf copies lock value: command-line-arguments.Demo contains sync.Mutex
./nocopy.go:29:11: call of TestFunc copies lock value: command-line-arguments.Demo contains sync.Mutex
./nocopy.go:31:20: call of fmt.Printf copies lock value: command-line-arguments.Demo contains sync.Mutex
*/
