package main

import (
	"errors"
	"fmt"
)

type User struct {
	Name   string
	Number string
	Email  string
}

func main() {
	u := map[int]User{}
	for {
		var cmd string
		fmt.Print("команда: ")
		fmt.Scan(&cmd)
		if cmd == "стоп" {
			return
		}

		if cmd == "добавить" {
			var name, phone, email string
			fmt.Scan(&name, &phone, &email)
			n, err := Add(u, name, phone, email)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(n)
			continue
		}

		if cmd == "найти" {
			var id int
			fmt.Scan(&id)
			n, err := Search(u, id)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(n)
			continue
		}

		if cmd == "редактировать" {
			var id int
			var name, phone, email string
			fmt.Scan(&id, &name, &phone, &email)
			err := Edit(u, id, name, phone, email)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(u)
			continue
		}

		if cmd == "удалить" {
			var id int

			fmt.Scan(&id)
			err := deleteUser(u, id)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(u)
			continue
		}

		fmt.Println("неизвестная команда")

		//fmt.Println(u)

		//u[1] = User{Name: "Алексндр", Number: "89669801626", Email: "alex@mail.ru"}

		/*	Add(u, "Алексндр", "89669801626", "alex@mail.ru")
			Add(u, "Мария", "896777726", "Maria@mail.ru")
			fmt.Println(u)
			n, err := Add(u, "", "5678568", "Maria@mail.ru")
			if err != nil {
				fmt.Println(err)

			}
			fmt.Println(n)
			n3, err := Search(u, 1)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(n3)

			n2, err := Search(u, 5)
			if err != nil {
				fmt.Println(err)
				//return
			}
			fmt.Println(n2)

			err = Edit(u, 1, "Сергей", "89669801626", "alex@mail.ru")
			if err != nil {
				fmt.Println(err)
				//return
			}
			fmt.Println(u)

			err = Edit(u, 5, "Сергей", "89669801626", "alex@mail.ru")
			if err != nil {
				fmt.Println(err)
				//return
			}
			fmt.Println(u)

			deleteUser(u, 1)
			fmt.Println(u)
			err = deleteUser(u, 5)
			if err != nil {
				fmt.Println(err)
				//return
			}
			fmt.Println(u)

			Add(u, "Иван", "8966944801626", "alыпвex@mail.ru")
			fmt.Println(u)
		*/
	}
}

func Add(a map[int]User, Name, Number, Email string) (int, error) {
	nBig := 0

	if Name == "" || Number == "" {
		return 0, errors.New("Пустое имя или телефон")
	}
	for id := range a {
		if id > nBig {
			nBig = id
		}
	}
	n := nBig + 1
	a[n] = User{Name: Name, Number: Number, Email: Email}

	return n, nil

}

func Search(a map[int]User, b int) (User, error) {

	v, ok := a[b]
	if !ok {
		return v, errors.New("Такого пользователя нет")
	} else {
		return v, nil
	}
}

func Edit(a map[int]User, b int, Name, Number, Email string) error {
	_, ok := a[b]
	if !ok {
		return errors.New("Такого пользователя нет")
	}
	if Name == "" || Number == "" {
		return errors.New("Пустое имя или телефон")
	} else {
		a[b] = User{Name: Name, Number: Number, Email: Email}
		return nil
	}

}

func deleteUser(a map[int]User, b int) error {
	_, ok := a[b]
	if !ok {
		return errors.New("Такого пользователя нет")
	}
	delete(a, b)
	return nil
}
