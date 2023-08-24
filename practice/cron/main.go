package main

import (
	"fmt"

	"github.com/robfig/cron"
)

// https://github.com/robfig/cron

func main() {
	c := cron.New()
	c.AddFunc("@every 1s", func() {
		fmt.Println("tick every 1 second run once")
	})
	c.Start()

	//
	select {}
}
