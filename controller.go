package main

import (
	"log"
	"os"
	"time"
	"github.com/olekukonko/tablewriter"
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
	log.Println("Beginning monitoring")
	table := tablewriter.NewWriter(os.Stdout)
loop:
	for {
		select {
		case <-quit:
			log.Println("received quit signal")
			break loop
		case <-time.After(c.interval):
			log.Println("Getting updates...")
			logEntries := c.monitor.GetUpdates(c.file)
			c.stats.Update(logEntries)
			c.stats.Display(table)
		}
	}
	log.Println("sending done signal")
	done <- 1
}
