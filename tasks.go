package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func setupTasks() {
	c := cron.New()
	c.AddFunc("30 * * * *", func() { fmt.Println("Every hour on the half hour") })
	c.AddFunc("30 3-6,20-23 * * *", func() { fmt.Println(".. in the range 3-6am, 8-11pm") })
	c.AddFunc("CRON_TZ=Asia/Tokyo 30 04 * * *", func() { fmt.Println("Runs at 04:30 Tokyo time every day") })
	c.AddFunc("@hourly", func() { fmt.Println("Every hour, starting an hour from now") })
	c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty, starting an hour thirty from now") })
	// c.AddFunc("@every 0h0m1s", func() { fmt.Printf("Hi %v\n", time.Now().Format("2006-01-02 15:04:05")) })
	c.AddFunc("@every 0h1m", func() { fmt.Printf("Hi %v\n", time.Now().Format("2006-01-02 15:04:05")) })
	c.Start()
}
