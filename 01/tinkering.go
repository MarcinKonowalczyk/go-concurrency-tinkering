package main

import (
	"fmt"
	"runtime"
	"time"
)



func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func main() {

	// NOTE: THIS IS BAD!! we are sending a value from many goroutines to a single channel!!!
	ch := make(chan int)

	fmt.Println("starting a pile of functions...")
	for range 10 {
		go func(N int, ch chan int) {

			defer func() {
				if r := recover(); r != nil {
					err_string := fmt.Sprintf("%v", r)
					// NOTE: comparing error strings is not great, but that's what we have here
					// ( we get runtime.plainError as the recovered type )
					// fmt.Printf("%T\n", r)
					if err_string == "send on closed channel" {
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
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
	
	cpuUsage := runtime.NumCPU()
	fmt.Printf("NumCPU = %v\n", cpuUsage)

	fmt.Println("N Goroutines:", runtime.NumGoroutine())
	fmt.Println("done")
}
