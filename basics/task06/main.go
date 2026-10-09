package main

import (
	"errors"
	"fmt"
)

type Score struct {
	Name string
	Cash int
}

func main() {

	var a Score
	a = Score{Name: "Иван", Cash: 100}

	a.Add(10)
	fmt.Println(a.Name, a.Cash)

	a.Removal(30)
	fmt.Println(a.Name, a.Cash)

	m, err := a.Removal(130)
	if err != nil {
		fmt.Println(err, a.Cash)
		return
	}
	fmt.Println(m)

	n, err2 := a.Add(-20)
	if err2 != nil {
		fmt.Println(err2, a.Cash)
		return
	}
	fmt.Println(n)

	fmt.Println(a.Name, a.Cash)

}

func (b *Score) Add(n int) (int, error) {

	if n <= 0 {
		return 0, errors.New("Сумма пополдения ноль или меньше нуля, остаток на балансе:")
	}
	b.Cash += n
	return b.Cash, nil
}

func (c *Score) Removal(n int) (int, error) {
	if n <= 0 {
		return n, errors.New("Сумма снятия ноль или меньше нуля, остаток на балансе:")
	}
	if n > c.Cash {
		return n, errors.New("Сумма снятия выше доступного баланса, остаток на балансе:")
	} else {
		c.Cash -= n
		return c.Cash, nil
	}

}
