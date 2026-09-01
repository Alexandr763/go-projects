package main

import (
	"fmt"
)

func main() {

	s := "Go语言"

	fmt.Println("Длина в байтах:", len(s))
	fmt.Println("Длина в символах", len([]rune(s)))
	fmt.Println("Символы:")

	for _, r := range s {
		fmt.Println(string(r))
	}

}
