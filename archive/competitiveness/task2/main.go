package main

import "fmt"

func main() {
	ch := make(chan int)

	go func() {
		ch <- 10
		ch <- 20
		ch <- 30
		close(ch)
	}()

	for value := range ch {
		fmt.Println(value)
	}

}
