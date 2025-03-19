package main

import (
	"fmt"
	"go-concurrency-tinkering/utils"
	"sync"
)

type not_atomic_int64 struct {
	value int64
}

func (a *not_atomic_int64) Add(delta int64) {
	a.value += delta
}

func (a *not_atomic_int64) Load() int64 {
	return a.value
}

func main() {
	
	var counter not_atomic_int64
	// var counter atomic.Int64

	var wg sync.WaitGroup

	N := 50
	M := 1000

	for range N {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < M; j++ {
				counter.Add(1)
			}
		}()
	}

	wg.Wait()

	fmt.Println("counter:", counter.Load())

	if counter.Load() != int64(N*M) {
		fmt.Printf("ocounter... something went wrong. Expected %d, got %d\n", N*M, counter.Load())
	}

	utils.PrintMemStats()
	fmt.Println("done")
}
