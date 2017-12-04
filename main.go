package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	quit := make(chan int)
	done := make(chan int)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go monitor(quit, done)
	for range signalChan {
		fmt.Println("sending quit signal")
		quit <- 1
		fmt.Println("waiting for monitoring to be done")
		<- done
		fmt.Println("done")
		break
	}
}

func monitor(quit chan int, done chan int) {
	fmt.Println("beginning monitoring")
	loop:
	for {
		select {
		case <- quit:
			fmt.Println("received quit signal")
			break loop
		default:
		}
	}
	fmt.Println("sending done signal")
	done <- 1
}