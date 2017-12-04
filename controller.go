package main

import (
	"log"
	"os"
	"time"
)

type Controller struct {
	file     *os.File
	interval time.Duration
	monitor  Monitor
	stats    *SiteStats
}

func NewController(filePath string, timeInSeconds int) *Controller {
	if file, err := os.Open(filePath); err == nil {
		return &Controller{
			file:     file,
			interval: time.Duration(timeInSeconds) * time.Second,
			monitor:  Monitor{},
			stats:    NewSiteStats()}
	} else {
		log.Fatal(err)
		return nil
	}
}

func (c *Controller) Monitor(quit chan int, done chan int) {
	log.Println("beginning monitoring")
loop:
	for {
		select {
		case <-quit:
			log.Println("received quit signal")
			break loop
		case <-time.After(c.interval):
			log.Println("Interval")
			logEntries := c.monitor.GetUpdates(c.file)
			c.stats.update(logEntries)
		}
	}
	log.Println("sending done signal")
	done <- 1
}
