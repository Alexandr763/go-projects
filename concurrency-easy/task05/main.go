package main

import "fmt"

func main() {

	ch := make(chan int)
	a := 5
	go square(a, ch)

	b := 7
	go square(b, ch)

	c := 9
	go square(c, ch)

	n := <-ch
	fmt.Println(n)

	n2 := <-ch
	fmt.Println(n2)

	n3 := <-ch
	fmt.Println(n3)

}

func square(n int, ch chan int) {
	a := n * n

	ch <- a

}
