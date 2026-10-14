package main

import "fmt"

func main() {
	ch := make(chan int)

	go func() {
		ch <- 1
		ch <- 2
		ch <- 3
		ch <- 4
		ch <- 5
	}()

	n := <-ch
	fmt.Println(n)

	n2 := <-ch
	fmt.Println(n2)

	n3 := <-ch
	fmt.Println(n3)

	n4 := <-ch
	fmt.Println(n4)

	n5 := <-ch
	fmt.Println(n5)
}
