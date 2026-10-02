package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	ch := make(chan int)
	var wg sync.WaitGroup

	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()

	wg.Add(3)
	go func() {
		defer wg.Done()
		for {
			select {
			case a, ok := <-ch:
				duration := time.Duration(500) * time.Millisecond
				time.Sleep(duration)
				if ok != true {
					return
				}
				if ok != false {
					fmt.Printf("worker %d processed task %d\n", 1, a)

				}
			case <-ctx.Done():
				return

			}
		}

	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case a, ok := <-ch:
				duration := time.Duration(500) * time.Millisecond
				time.Sleep(duration)
				if ok != true {
					return
				}
				if ok != false {
					fmt.Printf("worker %d processed task %d\n", 2, a)

				}
			case <-ctx.Done():
				return

			}
		}

	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case a, ok := <-ch:
				duration := time.Duration(500) * time.Millisecond
				time.Sleep(duration)
				if ok != true {
					return
				}
				if ok != false {
					fmt.Printf("worker %d processed task %d\n", 3, a)

				}
			case <-ctx.Done():
				return

			}
		}

	}()
	/*
		go func() {
			defer wg.Done()
			for value := range ch {
				//fmt.Println(value)
				duration := time.Duration(500) * time.Millisecond
				time.Sleep(duration)
				fmt.Printf("worker %d processed task %d\n", 2, value)

			}
		}()
		go func() {
			defer wg.Done()
			for value := range ch {
				//fmt.Printf("worker", %d value)
				duration := time.Duration(500) * time.Millisecond
				time.Sleep(duration)
				fmt.Printf("worker %d processed task %d\n", 3, value)
			}
		}() */
	//tasks := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	tasks := make([]int, 100)
	for i := range tasks {
		tasks[i] = i + 1
	}

	for _, v := range tasks {
		select {
		case ch <- v:

		case <-ctx.Done():

		}
		//cancel()
		//wg.Wait()
		//ch <- 1
		//ch <- 2
		//ch <- 3

		//wg.Wait()

	}
	wg.Wait()
}
