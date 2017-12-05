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
	filePath, interval, alertThreshold := parseArgs()
	quit, done := beginMonitoring(NewController(filePath, interval, alertThreshold))
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

func parseArgs() (filePath string, interval int, alertThreshold int) {
	var err error
	if arguments, err := docopt.Parse(usage, nil, true, version, false); err == nil {
		timeSeconds := defaultMonitoringInterval
		if arguments["--time"] != nil {
			proposed, _ := strconv.Atoi(arguments["--time"].(string))
			if proposed > 0 {
				timeSeconds = proposed
			}
		}
		alertThreshold := defaultHitCountToStartAlerting
		if arguments["--alert"] != nil {
			proposed, _ := strconv.Atoi(arguments["--alert"].(string))
			if proposed >= 0 {
				alertThreshold = proposed
			}
		}
		return arguments["<file>"].(string), timeSeconds, alertThreshold
	}
	log.Fatal(err)
	return
}

const version = "HTTP Log Monitoring Console Program 0.1"
const defaultMonitoringInterval = 10
const defaultHitCountToStartAlerting = 100
const usage = `HTTP Log Monitoring Console Program.

Usage:
  httpTrafficMonitoring [--time=10] [--alert=100.0] <file>

Options:
  -h --help     		Show this screen.
  --version     		Show version.
  -t --time=<seconds>	Time to refresh log stats in seconds [default: 10].
  --alert=<hits/second> Threshold to generate an alert that will persist until traffic falls below for two minutes.`
