package utils

import (
	"fmt"
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
	fmt.Println("done")
}