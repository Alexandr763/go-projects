package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		time.Sleep(500 * time.Millisecond)

		ch1 <- 1

	}()

	go func() {
		time.Sleep(500 * time.Millisecond)
		ch2 <- 2

	}()

	select {
	case v := <-ch1:
		fmt.Println("первый", v)
	case v := <-ch2:
		fmt.Println("второй", v)
	}
}
