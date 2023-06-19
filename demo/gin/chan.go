package main

import (
	"fmt"
	"math/rand"
	"time"
)

func ch() {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 2)

	go func() {
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		ch1 <- 1
	}()

	go func() {
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		ch2 <- 1
	}()

BPOINT:
	for {
		select {
		case <-ch1:
			fmt.Println("ch1")
			break BPOINT
		case <-ch2:
			fmt.Println("ch2")
			break BPOINT
		}
	}
	fmt.Println("end")
}
