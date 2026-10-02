package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		Worker(ctx)
	}()
	duration := time.Duration(5000) * time.Millisecond
	time.Sleep(duration)
	cancel()
	wg.Wait()

}

func Worker(ctx context.Context) {

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("worker stopped")
			ticker.Stop()
			return

		case <-ticker.C:
			fmt.Println("worker is working")

		}
	}

}
