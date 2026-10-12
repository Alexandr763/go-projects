package main

import (
	//"errors"
	"fmt"
)

type Task struct {
	Id   int
	Name string
	Done bool
}

func main() {

	a := []Task{}
	a = Add(a, Task{Id: 1, Name: "Помыть полы", Done: false})
	fmt.Println("Список задач:")
	(List(a))
	CompletionMark(a, 1)
	fmt.Println("Список задач:")
	CompletionMark(a, 2)

	a = Add(a, Task{Id: 2, Name: "Помыть посуду", Done: false})

	a = Add(a, Task{Id: 2, Name: "Помыть шкаф", Done: false})

	(List(a))
	a = DeleteTask(a, 1)
	(List(a))
	a = DeleteTask(a, 2)
	(List(a))

}

func List(b []Task) {

	n := len(b)
	if n == 0 {
		fmt.Println("Задач нет")
	}
	for _, v := range b {
		fmt.Println(v)
	}
}

func Add(b []Task, c Task) []Task {
	for _, v := range b {
		if v.Id == c.Id {
			fmt.Println("Номер уже занят")
			return b

		}
	}
	b = append(b, c)
	return b
}

func CompletionMark(a []Task, b int) {
	//v := a.Done
	found := false
	for i, v := range a {
		if v.Id == b {
			a[i].Done = true
			found = true
		}
	}
	if !found {
		fmt.Println("Такой задачи нет в списке")
	}

}

func DeleteTask(a []Task, b int) []Task {
	for i, v := range a {
		if v.Id == b {
			a = append(a[:i], a[i+1:]...)
			return a

		}
	}
	fmt.Println("Такой задачи нет в списке")
	return a
}
