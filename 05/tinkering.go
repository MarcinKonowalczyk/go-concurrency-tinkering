package main

import (
	"fmt"
	"go-concurrency-tinkering/utils"
	"time"
)

type limiter struct {
	limiter chan time.Time
	stop    chan struct{}
}

func newLimiter(N int, interval time.Duration) (lim limiter) {
	lim = limiter{
		limiter: make(chan time.Time, N),
		stop:    make(chan struct{}),
	};

	for range cap(lim.limiter) {
		lim.limiter <- time.Now()
	}
	
	go func(lim limiter) {
		loop: for {
			select {
			case <-lim.stop:
				fmt.Println("stopping ticker")
				break loop
			case <-time.After(interval):
				fmt.Println("tick")
				lim.limiter <- time.Now()
			}
		}
		close(lim.limiter)
	}(lim)
	
	return
}

func (l *limiter) close() {
	l.stop <- struct{}{}
	close(l.stop)
}

func main() {
	N := 7
	requests := make(chan int, N)

	go func(N int) {
		for i := 0; i < N; i++ {
			requests <- i
		}
		close(requests)
	}(N)

	func() {
		limiter := newLimiter(3, 500*time.Millisecond)
		defer func() {
			fmt.Println("closing limiter")
			limiter.close()
		}()
		

		for req := range requests {
			<-limiter.limiter
			fmt.Println("request", req, time.Now())
		}
	}()

	utils.PrintMemStats()
	fmt.Println("done")
}
