package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"log"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	filePath, interval := parseArgs()
	quit, done := beginMonitoring(NewController(filePath, interval))
	handleSignals(quit, done)
}

func handleSignals(quit chan int, done chan int) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	for range signalChan {
		fmt.Println("sending quit signal")
		quit <- 1
		fmt.Println("waiting for monitoring to be done")
		<-done
		fmt.Println("done")
		break
	}
}

func beginMonitoring(controller *Controller) (quit chan int, done chan int) {
	q := make(chan int)
	d := make(chan int)

	go controller.Monitor(q, d)

	return q, d
}

func parseArgs() (filePath string, interval int) {
	var err error
	if arguments, err := docopt.Parse(usage, nil, true, version, false); err == nil {
		if timeString, _ := arguments["--time"]; timeString != nil {
			if timeSeconds, err := strconv.Atoi(timeString.(string)); err == nil {
				if timeSeconds <= 0 {
					timeSeconds = defaultInterval
				}
				return arguments["<file>"].(string), timeSeconds
			}
		} else {
			return arguments["<file>"].(string), defaultInterval
		}
	}
	log.Fatal(err)
	return
}

const version = "HTTP Log Monitoring Console Program 0.1"
const defaultInterval = 10
const usage = `HTTP Log Monitoring Console Program.

Usage:
  monitor [--time=10] <file>

Options:
  -h --help     		Show this screen.
  --version     		Show version.
  -t --time=<seconds>	Time to refresh log stats in seconds [default: 10].`
