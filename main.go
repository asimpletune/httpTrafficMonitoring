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

	controller := NewController("/Users/Spencer/Desktop/server/access.log", 3)
	go controller.Monitor(quit, done)

	for range signalChan {
		fmt.Println("sending quit signal")
		quit <- 1
		fmt.Println("waiting for monitoring to be done")
		<- done
		fmt.Println("done")
		break
	}
}
