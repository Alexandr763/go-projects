package main

import (
	"errors"
	"fmt"
)

type Wallet struct {
	Money int
}

func main() {
	var a Wallet
	//a = Wallet{Money: 0}

	fmt.Println("Баланс до:")
	fmt.Println(a.Money)
	a.Insert(100)
	fmt.Println("Баланс после:")
	fmt.Println(a.Money)

	n, err := a.Buy(20)
	if err != nil {
		fmt.Println(err, a.Money)
		return
	}
	fmt.Println("Баланс после:")
	fmt.Println(n)

	n, err2 := a.Buy(120)
	if err2 != nil {
		fmt.Println(err2, a.Money)
		//return
	}
	//fmt.Println("Баланс после:")
	//fmt.Println(n)

	n, err3 := a.Buy(-10)
	if err3 != nil {
		fmt.Println(err3, a.Money)
		//return
	}
	//fmt.Println("Баланс после:")
	//fmt.Println(n)

	n, err4 := a.Insert(-100)
	if err4 != nil {
		fmt.Println(err4, a.Money)
		//return
	}
	//fmt.Println("Баланс после:")
	//fmt.Println(n)

	//fmt.Println("Баланс после:")
	//a.Buy(50)
	//fmt.Println(a.Money)
}

func (a *Wallet) Insert(n int) (int, error) {
	if n <= 0 {
		return a.Money, errors.New("Сумма пополнения отрицательная или равна нулю, текущий баланс:")
	}
	a.Money += n
	return a.Money, nil
}

func (a *Wallet) Buy(n int) (int, error) {
	if n <= 0 {
		return a.Money, errors.New("Сумма покупки отрицательная, баланс:")
	}
	if n > a.Money {
		return a.Money, errors.New("Сумма покупки выше доступного баланса, остаток на балансе:")
	} else {

		a.Money -= n
		return a.Money, nil
	}
}
