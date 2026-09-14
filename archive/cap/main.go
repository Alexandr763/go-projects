package main

import "fmt"

func main() {
	a := make([]int, 2, 4)
	fmt.Printf("Старт: len=%d cap=%d %v\n", len(a), cap(a), a)

	a = append(a, 1)
	fmt.Printf("После append(1): len=%d cap=%d %v\n", len(a), cap(a), a)

	a = append(a, 2)
	fmt.Printf("После append(2): len=%d cap=%d %v\n", len(a), cap(a), a)

	a = append(a, 3)
	fmt.Printf("После append(3): len=%d cap=%d %v\n", len(a), cap(a), a)

	a = append(a, 4)
	fmt.Printf("После append(4): len=%d cap=%d %v\n", len(a), cap(a), a)

	a = append(a, 5)
	fmt.Printf("После append(5): len=%d cap=%d %v\n", len(a), cap(a), a)
}
