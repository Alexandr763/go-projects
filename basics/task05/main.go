package main

import "fmt"

func main() {
	Warehouse := map[string]int{}
	//Warehouse["Яблоки"] = 10
	NewProduct(Warehouse, "Яблоки", 20)
	Restock(Warehouse, "Яблоки", 5)
	NewProduct(Warehouse, "Груши", 0)
	Restock(Warehouse, "Груши", 5)

	NewProduct(Warehouse, "", 100)
	fmt.Println(Warehouse)
	//Warehouse["Яблоки"] = 10
	//NewProduct(Warehouse)
	//Warehouse["Яблоки"] = 0
	//NewProduct(Warehouse)
	//Warehouse["Яблоки"] = -1
	//NewProduct(Warehouse)
	Sell(Warehouse, "Яблоки", 5)
	Sell(Warehouse, "Яблоки", 0)
	fmt.Println(Warehouse)
	Sell(Warehouse, "Яблоки", 20)

	fmt.Println(Warehouse)
	Sell(Warehouse, "Груши", 8)
	fmt.Println(Warehouse)

}

func NewProduct(a map[string]int, name string, quantity int) {
	if name == "" {
		fmt.Println("пустое название")
		return
	}

	_, ok := a[name]
	if ok {
		fmt.Println("товар уже на складе")
		return
	}

	if quantity > 0 {
		a[name] = quantity
	} else {
		fmt.Println("количество должно быть больше нуля")
	}
}

func Restock(a map[string]int, name string, quantity int) {
	if quantity <= 0 {
		fmt.Println("количество должно быть больше нуля")
		return
	}

	v, ok := a[name]
	if !ok {
		fmt.Println("такого товара нет")
		return
	}
	a[name] = v + quantity
}

func Sell(a map[string]int, name string, quantity int) {
	v, ok := a[name]

	if quantity <= 0 {
		fmt.Println("количество должно быть больше нуля")
		return
	}
	if !ok {
		fmt.Println("такого товара нет")
		return
	}

	if quantity > v {
		fmt.Println("нельзя продать больше остатка")
		return
	}

	if ok {
		a[name] = v - quantity
	}
}
