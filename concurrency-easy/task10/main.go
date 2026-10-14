package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go func() {
		ticker := time.NewTicker(time.Second)

		for {

			select {
			case <-ctx.Done():
				fmt.Println("Время вышло")
				ch <- 1

				ticker.Stop()
				return
			case <-ticker.C:
				fmt.Println("worker is working")

			}

		}
	}()
	time.Sleep(5 * time.Second)

	<-ch

}
