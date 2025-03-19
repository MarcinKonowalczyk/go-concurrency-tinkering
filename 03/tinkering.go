package main

import (
	"fmt"
	"go-concurrency-tinkering/utils"
)

func main() {
	ch := make(chan int)
	abort := make(chan struct{})
	go func (N int, ch chan int, abort chan struct{}) {
		n_send := 0
		loop: for {
			select {
			case <-abort:
				break loop
			case ch <- n_send:
				n_send++
				if n_send >= N {
					break loop
				}
			// The default case here should not be used
			// default:
			// 	fmt.Println("default")
			}
			fmt.Println("loop")
		}
		close(ch)
	}(10, ch, abort)

	i := 0
	for value := range ch {
		fmt.Println(value)
		i++
		if i == 5 {
			fmt.Println("Sending abort signal...")
			close(abort)
		}
	}

	// Wait for the goroutine to finish
	fmt.Println("Waiting for goroutine to finish...", <-ch)

	utils.PrintMemStats()
	fmt.Println("done")
}
