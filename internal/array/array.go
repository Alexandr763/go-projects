package array

import (
	"errors"
	"fmt"
)

type Array struct {
	Data []int
	Len  int
}

func (a *Array) Push(value int) {
	a.Data = append(a.Data, value)
	a.Len = a.Len + 1
}

func (a *Array) Get(index int) (int, error) {
	if index < 0 || index >= a.Len {
		return 0, errors.New("Индекс все диапазона")
	}
	return a.Data[index], nil
}

func (a *Array) Delete(index int) error {
	if index < 0 || index >= a.Len {
		return errors.New("индекс вне диапазона")
	}
	a.Data = append(a.Data[:index], a.Data[index+1:]...)
	a.Len = a.Len - 1
	return nil
}

func (a *Array) Search(value int) int {
	for i, v := range a.Data {
		if v == value {
			return i
		}
	}
	return -1
}

func (a *Array) Reverse() {
	left := 0
	right := a.Len - 1
	for left < right {
		a.Data[left], a.Data[right] = a.Data[right], a.Data[left]
		left = left + 1
		right = right - 1
	}
}

func main() {
	arr := Array{Data: []int{}, Len: 0}

	arr.Push(10)
	arr.Push(20)
	arr.Push(30)

	fmt.Println(arr.Data)
	fmt.Println(arr.Len)

	val1, err := arr.Get(1) // сохраняем результат
	if err != nil {         // проверяем, есть ли ошибка
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Значение:", val1)
	}

	val2, err := arr.Get(99) // сохраняем результат
	if err != nil {          // проверяем, есть ли ошибка
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Значение:", val2)
	}

	fmt.Println("До удаления:", arr.Data)

	err2 := arr.Delete(0)
	if err2 != nil {
		fmt.Println("Ошибка удаления:", arr.Data)
	} else {
		fmt.Println("После удаления:", arr.Data)
	}

	err2 = arr.Delete(99)
	if err2 != nil {
		fmt.Println("Ошибка:", err2)
	}

	idx := arr.Search(20)
	if idx == -1 {
		fmt.Println("Элемент не найден")
	} else {
		fmt.Println("Индекс 20:", idx)
	}

	idx = arr.Search(99)
	if idx == -1 {
		fmt.Println("Элемент не найден")
	} else {
		fmt.Println("Индекс 99:", idx)
	}

	fmt.Println("До разворота:", arr.Data)
	arr.Reverse()
	fmt.Println("После разворота:", arr.Data)

	var s []int
	for i := 0; i < 5; i++ {
		s = append(s, i)
		fmt.Println("len:", len(s), "cap:", cap(s))
	}

	src := []int{5, 10, 15}
	src[0] = 999
	dst := make([]int, len(src))
	copy(dst, src)
	fmt.Println(dst)

	fmt.Println(dst)

}

/*
type Student struct {
	Name  string
	Grade int
}

func Promote(s *Student) {
	s.Grade = s.Grade +1
}

func main() {

	array := [3]int{1, 2, 3}
	slice := []int{10, 20, 30}
	slice = append(slice, 40)
	fmt.Println(array)
	fmt.Println(slice)
	student := Student{Name: "Alex", Grade: 4}
	fmt.Println("Имя студента:", student.Name)
	fmt.Println("Грейд студента:", student.Grade)
	fmt.Println("До повышения грейда студента:", student.Grade)

	Promote(&student)
	fmt.Println("После повышения грейда студента:", student.Grade)

}

/*
type User struct {
	Name string
	Age  int
}

func makeAdult(u *User) {
	u.Age = 18
	fmt.Println("Возраст юзера изменился:", u.Age)

}

func main() {

	user1 := User{Name: "Алексей", Age: 10}
	fmt.Println("Возраст юзера:", user1.Age)
	makeAdult(&user1)
}



type Book struct {
	Title  string
	Author string
	Pages  int
}

	Book1 := Book{Title: "Война и мир", Author: "Лев Толстой", Pages: 1125}
	fmt.Println(Book1.Title)
	fmt.Println(Book1.Author)
	fmt.Println(Book1.Pages)



numbers := [5]int{10, 20, 30, 40, 50}
fmt.Println(numbers[0])
fmt.Println(numbers[2])
fmt.Println(numbers[4])

slice := []int{5, 10, 15}
slice = append(slice, 20)
slice = append(slice, 25)

fmt.Println(slice)
*/
