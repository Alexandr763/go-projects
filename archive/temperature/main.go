package main

import "fmt"

func main() {

}
package main

import (
	"errors"
	"fmt"
)


/*
package main

import (
	"fmt"
)

func main() {
	numbers := []int{7, 2, 9, 4, 2, 9, 1}
	currentMin := numbers[0]
	for _, value := range numbers {
		if value < currentMin {
			currentMin = value
		}
	}
	fmt.Println("Минимальное значение:", currentMin)

	currentMax := numbers[0]
	for _, value := range numbers {
		if value > currentMax {
			currentMax = value
		}
	}
	fmt.Println("Максимальное значение:", currentMax)

	sum := 0
	for _, value := range numbers {
		sum = sum + value
	}
	fmt.Println("Сумма значений:", sum)

	countNum := 0
	for _, value := range numbers {
		if value == 2 {
			countNum += 1
		}
	}
	fmt.Println("Количество двоек в слайсе:", countNum)

	for index, value := range numbers {
		fmt.Println("Индекс:", index, "Значение:", value)

	}

	fmt.Println("Длина слайса:", len(numbers))
	fmt.Println("Вместимость слайса:", cap(numbers))

	numbers = append(numbers, 5)
	fmt.Println("Обновленный слайс:", numbers)
	fmt.Println("Длина обновленного слайса:", len(numbers))
	fmt.Println("Вместимость обновленного слайса:", cap(numbers))

	backup := make([]int, len(numbers))
	copied := copy(backup, numbers)
	fmt.Println("Копия слайса:", backup)
	fmt.Println("Скопировано элементов:", copied)

	backup[0] = 100
	fmt.Println("1 слайс:", backup)
	fmt.Println("2 слайс:", numbers)

}


// объявили структуру
type Array struct {
	array []int
	len   int
}

// функция конструктор
func NewArray() *Array {
	return &Array{
		array: make([]int, 0),
		len:   0,
	}
}

// метод Push
func (a *Array) Push(val int) error {
	a.array = append(a.array, val)
	a.len = len(a.array)
	return nil
}

func main() {
	//создали новый элемент
	arr := NewArray()

	//добавляем элемент
	arr.Push(5)

	//выводим результаты
	fmt.Println("Слайс:", arr.array)
	fmt.Println("Длина:", arr.len)

}
*/
