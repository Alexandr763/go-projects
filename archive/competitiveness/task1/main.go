package main


import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	num := []int{1, 2, 3}

	fmt.Println(Process(num))

	num1 := []int{}

	fmt.Println(Process(num1))

}

func Process(num []int) []int {
	var wg sync.WaitGroup

	ch := make(chan int, len(num))

	rand.Seed(time.Now().UnixNano())

	wg.Add(len(num))
	c := 0
	d := []int{}
	for _, v := range num {

		go func() {
			a := 100
			b := 500
			r := rand.Intn(b - a)
			duration := time.Duration(r) * time.Millisecond
			time.Sleep(duration)
			ch <- v * 2

			wg.Done()
		}()

	}
	n := len(num)
	for i := 0; i < n; i++ {
		c = <-ch
		d = append(d, c)
	}

	wg.Wait()

	return d

}

