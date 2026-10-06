package main

import "fmt"

func main() {
	num := []int{0, 7, 11, 5, -1, 10}
	num2 := []int{}
	num3 := []int{0, 11, 15}
	num4 := []int{1, 4, 6, 7, 8, 11}
	//fmt.Println(num)
	//numGrade(num)

	m := numGrade(num)
	numGrade(num2)
	numGrade(num3)
	m4 := numGrade(num4)

	minMax(m)
	//minMax(m2)
	//minMax(m3)
	minMax(m4)

	avg, _ := mean(m)
	avg2, _ := mean(m4)

	aboveAverage(m, avg)
	aboveAverage(m4, avg2)

}

func numGrade(n []int) []int {
	a := []int{}

	for _, v := range n {
		if 1 <= v && v <= 10 {
			a = append(a, v)
		}
	}

	if len(a) == 0 {
		fmt.Println("Нет данных для анализа")
		return a
	}

	fmt.Println(a)
	return a
}

func minMax(n []int) (int, int) {
	min := n[0]
	max := n[0]
	if len(n) > 0 {

		for i := 1; i < len(n); i++ {
			if n[i] < min {
				min = n[i]
			}
			if n[i] > max {
				max = n[i]
			}
		}

	}
	fmt.Println("min:", min, "max:", max)
	return min, max
}

func mean(n []int) (float64, bool) {
	if len(n) == 0 {
		return 0, false
	}
	sum := 0
	for _, n := range n {
		sum += n
	}
	fmt.Println("Средняя оценка:", float64(sum)/float64(len(n)))

	return float64(sum) / float64(len(n)), true

}

func aboveAverage(n []int, m float64) int {

	count := 0
	for _, v := range n {
		if float64(v) > float64(m) {
			count = count + 1
		}

	}

	fmt.Println("Колличество оценок выше среднего:", count)
	return count
}
