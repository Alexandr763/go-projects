package main

import (
	"fmt"
	"sort"
)

func findMinMax(nums []int) (min, max int, ok bool) {
	if len(nums) == 0 {
		return 0, 0, false
	}
	min = nums[0]
	max = nums[0]

	for _, v := range nums[1:] {
		if v < min {
			min = v
		} else if v > max {
			max = v
		}
	}
	return min, max, true
}

func CountDoubles(nums []int) map[int]int {
	counts := make(map[int]int)
	for _, v := range nums {
		counts[v]++
	}
	return counts
}

func main() {

	numbers := []int{5, 2, 8, 2, 19, 1, 5, 5} //слайс чисел
	min, max, ok := findMinMax(numbers)
	if !ok {
		fmt.Println("Срез пуст")
		return
	}
	fmt.Printf("Min: %d, Max: %d\n", min, max)

	dupCounts := CountDoubles(numbers)
	fmt.Println("Сколько раз встретилось каждое числоо:")
	for num, count := range dupCounts {
		if count > 1 {
			fmt.Printf("%d встречается %d раз\n", num, count)
		}
	}

	keys := []int{}
	for key := range dupCounts {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	fmt.Println("Уникальные значения:", keys)

}
