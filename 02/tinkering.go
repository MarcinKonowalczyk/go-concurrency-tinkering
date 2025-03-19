package main

import (
	"fmt"
	"go-concurrency-tinkering/utils"
	"math/rand/v2"
	"sync"
)

// func isChannelClosed(ch chan int) bool {
// 	select {
// 	case _, ok := <-ch:
// 		if !ok {
// 			return true
// 		}
// 	default:
// 		return false
// 	}
// }

func main() {

	// NOTE: THIS IS BAD!! we are sending a value from many goroutines to a single channel!!!
	ch := make(chan int)

	fmt.Println("starting a pile of functions...")

	wg := sync.WaitGroup{}
	for range 10 {
		wg.Add(1)
		go func(N int, ch chan int) {
			id := rand.Int() % 10000
			defer func() {
				if r := recover(); r != nil {
					if utils.IsClosedChannelErr(r) {
						fmt.Printf("%d: channel already closed\n", id)
					} else {
						// some other panic. continue panicking
						panic(r)
					}
				}
			}()
			defer func() {
				fmt.Printf("%d: waitgroup done\n", id)
				wg.Done()
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
	fmt.Println("waiting for all goroutines to finish...")
	wg.Wait()
	
	utils.PrintMemStats()
}
