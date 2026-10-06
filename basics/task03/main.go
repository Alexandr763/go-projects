package main

import (
	"fmt"
	"strings"
)

func main() {

	word1 := "go, Go sql!"
	word2 := "Go"
	word3 := ""

	//fmt.Println(say(word1))
	//fmt.Println(say(word2))
	//fmt.Println(say(word3))

	m1 := say(word1)
	for k, n := range m1 {
		fmt.Println("word1:", k, n)
	}

	m2 := say(word2)
	for k, n := range m2 {
		fmt.Println("word2:", k, n)
	}

	m3 := say(word3)
	for k, n := range m3 {
		fmt.Println("word3:", k, n)
	}
}

func say(s string) map[string]int {
	m := map[string]int{}
	s = strings.ToLower(s)

	parts := strings.Fields(s)

	for _, w := range parts {
		w = strings.Trim(w, ".,!")
		m[w] = m[w] + 1
		//fmt.Println(w)
	}
	//fmt.Println(m)
	return m
}
