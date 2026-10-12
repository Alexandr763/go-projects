package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	var a string
	a = "42"

	//Transformation(a)

	n, err := Transformation(a)
	if err != nil {
		fmt.Println(err)
		//return
	}
	fmt.Println(n)

	b := "Привет"
	_, err2 := Transformation(b)
	if err2 != nil {
		fmt.Println(err2)
		//return
	}

	c := ""
	_, err3 := Transformation(c)
	if err3 != nil {
		fmt.Println(err3)
		//return
	}

	d := "0"
	_, err4 := Transformation(d)
	if err4 != nil {
		fmt.Println(err4)
		//return
	}

	e := "-5"
	_, err5 := Transformation(e)
	if err5 != nil {
		fmt.Println(err5)
		//return
	}

}

func Transformation(a string) (int, error) {
	var ErrEmpty = errors.New("Пустая строка")
	var ErrBad = errors.New("Не число")

	if strings.TrimSpace(a) == "" {
		//fmt.Println(ErrBad)
		return 0, ErrEmpty
	}

	n, err := strconv.Atoi(a)
	if err != nil {
		//fmt.Println(ErrBad)
		return 0, ErrBad

	}

	if n == 0 {
		return 0, errors.New("Число равно нулю")
	}
	//fmt.Println(n)

	if n < 0 {
		return 0, errors.New("Число отрицательное")
	}
	return n, nil
}
