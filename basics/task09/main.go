package main

import (
	"fmt"
	"sort"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	var Humans []Person

	ListPerson(Humans)
	Humans = append(Humans, Person{Name: "Маша", Age: 30})
	Humans = append(Humans, Person{Name: "Игорь", Age: 30})

	Humans = append(Humans, Person{Name: "Оля", Age: 20})

	ListPerson(Humans)
	b := SortAge(Humans)
	ListPerson(b)

	c := TrueAge(Humans, 20, 25)
	ListPerson(c)

}

func ListPerson(a []Person) {

	n := len(a)
	if n == 0 {
		fmt.Println("Людей нет")
	}
	for _, v := range a {
		fmt.Println(v.Name, v.Age)
	}
}

func SortAge(a []Person) []Person {
	cp := make([]Person, len(a))
	copy(cp, a)
	sort.Slice(cp, func(i, j int) bool {
		if cp[i].Age != cp[j].Age {
			return cp[i].Age < cp[j].Age
		}
		return cp[i].Name < cp[j].Name

	})
	return cp
}

func TrueAge(a []Person, b int, c int) []Person {
	n := []Person{}
	for _, v := range a {
		if b <= v.Age && v.Age <= c {
			n = append(n, v)
		}
	}
	return n
}
