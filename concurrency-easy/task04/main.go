package main

import "fmt"

func main() {
	ch := make(chan string)

	go sendMessages(ch)

	for v := range ch {
		fmt.Println(v)
	}

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
