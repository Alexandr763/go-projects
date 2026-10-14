package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)

	a := "worker 1"
	go Name(a, ch)

	b := "worker 2"
	go Name(b, ch)

	<-ch
	<-ch

}

func Name(v string, d chan int) {
	for i := 0; i < 5; i++ {
		fmt.Println(v)
		time.Sleep(time.Second)
	}
	d <- 1
}
