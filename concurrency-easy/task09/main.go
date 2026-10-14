package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(time.Second)

		for {

			select {
			case <-ctx.Done():
				fmt.Println("worker stopped")
				ch <- 1

				ticker.Stop()
				return
			case <-ticker.C:
				fmt.Println("worker is working")

			}

		}
	}()
	time.Sleep(5 * time.Second)
	cancel()
	<-ch

}
