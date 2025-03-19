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

	ch := make(chan int)

	fmt.Println("starting a pile of ...")
	for range 1000 {
		go func(N int, ch chan int) {
			// fmt.Println("hi")
			for i := 0; i < N; i++ {
				ch <- i
			}
			close(ch)
		}(100, ch)
	}

	// for value := range ch {
	// 	fmt.Println(value)
	// }

	fmt.Println("waiting for a bit...")
	time.Sleep(time.Second * 1)
	
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
