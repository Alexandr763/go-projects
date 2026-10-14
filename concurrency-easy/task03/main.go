package main

import "fmt"

func main() {
	ch := make(chan string)

	go sendMessages(ch)
	n := <-ch
	n2 := <-ch
	n3 := <-ch

	fmt.Println(n)
	fmt.Println(n2)
	fmt.Println(n3)

}

func sendMessages(ch chan string) {
	a := "hello"
	b := "how are you?"
	c := "bye"

	ch <- a
	ch <- b
	ch <- c
	close(ch)
}
