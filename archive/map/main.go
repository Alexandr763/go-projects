package main

import "fmt"

func main() {
	map1 := map[string]string{
		"Anna":   "82354235634",
		"Boris":  "43654674574",
		"Viktor": "23454385634",
	}

	fmt.Println(map1)
	fmt.Println(map1["Anna"])

	value, ok := map1["Boris"]
	if ok == true {
		fmt.Println("Boris есть в мапе, его номер:", value)
	} else {
		fmt.Println("Boris нет в мапе")
	}

	map1["Galina"] = "435667457754"
	delete(map1, "Viktor")
	for key, value2 := range map1 {
		fmt.Println(key, value2)

	}

}
