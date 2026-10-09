package main

import "fmt"

type Expense struct {
	Sum      int
	Category string
	Comment  string
}

func main() {
	allExpenses := []Expense{}
	fmt.Println(allExpenses)

	byCat0 := sumByCategories(allExpenses)
	fmt.Println("Общая сумма расходов по категориям:")
	for k, n := range byCat0 {
		fmt.Println(k, n)
	}

	allExpenses = Add(allExpenses, Expense{Sum: 500, Category: "Колбаса", Comment: "куплено"})
	allExpenses = Add(allExpenses, Expense{Sum: 1000, Category: "Мясо", Comment: "Не куплено"})
	fmt.Println(allExpenses)
	allExpenses = Add(allExpenses, Expense{Sum: 300, Category: "Сок", Comment: "куплено"})
	fmt.Println(allExpenses)
	allExpenses = Add(allExpenses, Expense{Sum: 0, Category: "Шоколад", Comment: "Не куплено"})
	fmt.Println(allExpenses)
	allExpenses = Add(allExpenses, Expense{Sum: -100, Category: "Вода", Comment: "Не куплено"})
	fmt.Println(allExpenses)
	allExpenses = Add(allExpenses, Expense{Sum: 300, Category: "Сок", Comment: "куплено"})

	fmt.Println("Общая сумма расходов:", AllSum(allExpenses))
	//fmt.Println("Общая сумма расходов по категориям:", sumByCategories(allExpenses))
	byCat := sumByCategories(allExpenses)
	fmt.Println("Общая сумма расходов по категориям:")
	for k, n := range byCat {
		fmt.Println(k, n)
	}

}

func Add(c []Expense, b Expense) []Expense {
	//for _, b := range c {
	if b.Sum <= 0 {
		fmt.Println("Сумма меньше нуля или ноль")
		return c
	}
	//	}
	c = append(c, b)

	return c
}

func AllSum(c []Expense) int {
	total := 0
	for _, it := range c {
		total += it.Sum
		//fmt.Println("Общий расход:", it.Sum)

	}

	return total
}

func sumByCategories(m []Expense) map[string]int {
	a := map[string]int{}

	for _, it := range m {
		a[it.Category] = a[it.Category] + it.Sum
		//total += it.Sum
		//fmt.Println("Общий расход:", it.Sum)

	}

	return a
}
