package utils

import (
	"fmt"
	"reflect"
	"runtime"
)

func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func IsClosedChannelErr(err interface{}) bool {
	if err == nil {
		return false
	}
	var str string
	
	switch err.(type) {
	case error:
		str = err.(error).Error()
	case string:
		str = err.(string)
	default:
		return false
	}

	if str == "send on closed channel" {
		return true
	}
	if str == "close of closed channel" {
		return true
	}
	return false
}

func PrintMemStats() {
	fmt.Println("Memory stats:")
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", BToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", BToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", BToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
	
	cpuUsage := runtime.NumCPU()
	fmt.Printf("NumCPU = %v\n", cpuUsage)

	fmt.Println("N Goroutines:", runtime.NumGoroutine())
}


// Dynamic fan-in: multiple channels to one channel
func FanIn[T any](in []chan T, out chan T) {
	// Create a dynamic select statement
	cases := make([]reflect.SelectCase, 0)
	for _, c := range in {
		if c == nil {
			continue
		}
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(c),
		})
	}

	loop: for {
		chosen, value, ok := reflect.Select(cases)
		if !ok {
			cases = append(cases[:chosen], cases[chosen+1:]...)
			if len(cases) == 0 {
				break loop
			}
			continue loop
		}
		// send the value to the output channel
		out <- value.Interface().(T)
	}
	close(out)
}

// Dynamic fan-out: one channel to multiple channels
func FanOut[T any](in chan T, out []chan T) {
	value, ok := <-in;
	if !ok {
		for _, c := range out {
			// close the channel if it is not nil
			if c == nil {
				continue
			}
			close(c)
		}
		return
	}
	for _, c := range out {
		c <- value
	}	
}
