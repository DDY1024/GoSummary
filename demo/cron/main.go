package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
)

// 单机cron: github.com/robfig/cron

func main() {
	c := cron.New()
	c.AddFunc("@every 1s", func() {
		fmt.Println("tick every 1 second run once")
	})
	c.Start()
	time.Sleep(time.Second * 10)
}
