package main

import (
	"context"
	"fmt"
	"time"
)

func sleepAndTalk(ctx context.Context, duration time.Duration, message string) {
	select {
	case <-time.After(duration):
		fmt.Println(message)
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}
}

func main() {
	fmt.Println("Hello, World!")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	_ = cancel
	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	sleepAndTalk(ctx, 1*time.Second, "Hello, World!")
}
