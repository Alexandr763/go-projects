package main

import (
	"fmt"
	"time"
)

func main() {

	ch := make(chan int)

	select {
	case v := <-ch:
		fmt.Println("значение", v)

	default:
		fmt.Println("сейчас пусто")
	}

	go func() {
		ch <- 1
	}()

	time.Sleep(2 * time.Second)

	select {
	case v := <-ch:
		fmt.Println("значение", v)

	default:
		fmt.Println("сейчас пусто")
	}
}
