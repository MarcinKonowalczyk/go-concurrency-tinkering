package main

import (
	"fmt"
	"sync"
	"time"
)

type none struct{}


// MOTE: The real mutex implementation is much more better than this.
// This is just a toy version to try to use chanel as a mutex.
type ChannelMutex struct {
	ch chan none
}

func NewMyMutex() *ChannelMutex {
	return &ChannelMutex{
		ch: make(chan none, 1),
	}
}

func (m *ChannelMutex) Lock() {
	m.ch <- none{}
}

func (m *ChannelMutex) Unlock() {
	<-m.ch
}

func NewMutex() *sync.Mutex {
	return &sync.Mutex{}
}

type NewFunc = func() *Mutex

// OPTION 1: sync.Mutex
// type Mutex = sync.Mutex
// var NewMutexFunc NewFunc = NewMutex

// OPTION 2: ChannelMutex
type Mutex = ChannelMutex
var NewMutexFunc NewFunc = NewMyMutex

func worker(id int, mu *Mutex, ch chan none, sleepTime time.Duration) {
	fmt.Printf("[%d] starting\n", id)
	mu.Lock()
	fmt.Printf("[%d] has the lock\n", id)
	defer func() {
		mu.Unlock()
		fmt.Printf("[%d] released the lock\n", id)
	}()
	time.Sleep(sleepTime)
	fmt.Printf("[%d] doing work\n", id)
	ch <- none{}
}

func main() {
	mu := NewMutexFunc()
	ch := make(chan none, 1)
	
	go worker(1, mu, ch, 1*time.Second)
	
	// NOTE: There is a race between worker 1 and worker 2 starting.
	// This is not the point of the exercise here, so we just add a sleep.
	time.Sleep(100 * time.Millisecond)

	go worker(2, mu, ch, 0)

	// synch with both workers finishing
	<-ch
	<-ch
	fmt.Println("done")
}