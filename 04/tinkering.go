package main

import (
	"fmt"
	"go-concurrency-tinkering/utils"
	"runtime"
)


func main() {
	// NOTE: we're making a purposeful mistake here
	// channels and aborts are length 5, but we only start 2 goroutines
	N := 5
	M := 2

	channels := make([]chan int, N)
	aborts := make([]chan struct{}, N)
	
	for i := 0; i < M; i++ {
		channels[i] = make(chan int)
		aborts[i] = make(chan struct{})
		go func (N int, offset int, ch chan int, abort chan struct{}) {
			n := 0 // number of values sent
			loop: for {
				select {
				case <-abort:
					break loop
				case ch <- (n + offset):
					n++
					if n >= N {
						break loop
					}
				}
			}
			close(ch)
		}(9, i*10, channels[i], aborts[i])
	}

	// Create a channel to receive values from all channels
	ch := make(chan int)
	go utils.FanIn(channels, ch)

	abort := make(chan struct{})
	go utils.FanOut(abort, aborts)


	received := make([]int, 0)
	i := 0
	for value := range ch {
		fmt.Println(value)
		received = append(received, value)
		if i == 5 {
			fmt.Println("Sending abort signal...")
			close(abort)
		}
		i++
	}

	fmt.Println("Received values:", received)

	// Wait for the goroutine to finish
	// We can use the value from the channel since the FanIn will close the channel
	// when all the channels are closed
	fmt.Println("Waiting for goroutine to finish...", <-ch)

	// print the current running goroutines
	fmt.Println("Current Goroutines:", runtime.NumGoroutine())

	utils.PrintMemStats()
	fmt.Println("done")
}
