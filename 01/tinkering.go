package main

import (
	"fmt"
	"time"

	"go-concurrency-tinkering/utils"
)

func main() {

	// NOTE: THIS IS BAD!! we are sending a value from many goroutines to a single channel!!!
	ch := make(chan int)

	fmt.Println("starting a pile of functions...")
	for range 10 {
		go func(N int, ch chan int) {

			defer func() {
				if r := recover(); r != nil {
					// NOTE: comparing error strings is not great, but that's what we have here
					// ( we get runtime.plainError as the recovered type )
					// fmt.Printf("%T\n", r)
					if utils.IsClosedChannelErr(r) {
						fmt.Println("channel already closed")
						// panic(r)
					} else {
						// some other panic. continue panicking
						panic(r)
					}
				}
			}()

			for i := 0; i < N; i++ {
				ch <- i
			}
			close(ch)
		}(10, ch)
	}

	for value := range ch {
		fmt.Println(value)
	}

	// wait for all the goroutines to finish
	// NOTE: again, not a fantastic way of doing this. If we want to guarantee that all goroutines finish, we should use a WaitGroup
	// but this is just a quick and dirty
	fmt.Println("waiting for a bit...")
	time.Sleep(100 * time.Millisecond)

	utils.PrintMemStats()
	fmt.Println("done")
}
