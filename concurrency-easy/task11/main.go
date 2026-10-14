package main

import (
	"context"
	"fmt"
	//"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch := make(chan int)

	go func() {
		time.Sleep(1 * time.Second)
		select {
		//time.Sleep(1 * time.Second)
		case ch <- 1:
			//time.Sleep(1 * time.Second)
		case <-ctx.Done():
			return
		}
		time.Sleep(1 * time.Second)
		select {
		case ch <- 2:
			//time.Sleep(1 * time.Second)
		case <-ctx.Done():
			return
		}
		time.Sleep(1 * time.Second)
		select {
		case ch <- 3:
			//time.Sleep(1 * time.Second)
		case <-ctx.Done():
			return
		}

		close(ch)
	}()
	for {
		select {
		case v := <-ch:
			fmt.Println(v)
		case <-ctx.Done():
			fmt.Println("Время вышло")
			return
		}

	}

}
