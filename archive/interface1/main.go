package main

import "fmt"

type Speaker interface {
	Speak() string
}

type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return "Гав!"

}

type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return "Мяу!"

}

func main() {
	s := []Speaker{
		Dog{Name: "Шарик"},
		Cat{Name: "Мурка"},
	}

	for _, v := range s {
		fmt.Println(v.Speak())
	}
}
