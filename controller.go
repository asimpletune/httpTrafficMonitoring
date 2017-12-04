package main

import (
	"log"
	"os"
	"time"
)

const defaultAlertDuration = time.Duration(10) * time.Second

type Controller struct {
	file           *os.File
	interval       time.Duration
	monitor        Monitor
	stats          *SiteStats
	alertThreshold int
}

func NewController(filePath string, timeInSeconds int, alertThreshold int) *Controller {
	if file, err := os.Open(filePath); err == nil {
		return &Controller{
			file:           file,
			interval:       time.Duration(timeInSeconds) * time.Second,
			monitor:        Monitor{},
			stats:          NewSiteStats(),
			alertThreshold: alertThreshold}
	} else {
		log.Fatal(err)
		return nil
	}
}

func (c *Controller) Monitor(quit chan int, done chan int) {
	log.Println("Beginning monitoring")
	view := NewView()
	alert := NewAlert(float64(c.alertThreshold), defaultAlertDuration, c.stats)
	warn := make(chan bool)
	onAlert := false
	go alert.Monitor(warn)
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
			select {
			case onAlert = <-warn:
				view.Display(c.stats, onAlert)
			default:
				view.Display(c.stats, onAlert)
			}
		}
	}
	log.Println("sending done signal")
	done <- 1
}
