package main

import (
	"fmt"
	"go-concurrency-tinkering/utils"
)


func main() {
	ch := make(chan int)
	abort := make(chan struct{})
	go func (N int, ch chan int, abort chan struct{}) {
		loop: for i := 0; i < N; i++ {
			select {
			case ch <- i:
			case <-abort:
				break loop
			}
		}
		close(ch)
	}(10, ch, abort)

	i := 0
	for value := range ch {
		fmt.Println(value)
		i++
		if i == 5 {
			fmt.Println("Aborting...")
			close(abort)
			break
		}
	}

	// Wait for the goroutine to finish
	<-ch

	utils.PrintMemStats()
	fmt.Println("done")
}
