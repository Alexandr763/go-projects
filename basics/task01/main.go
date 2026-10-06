package main

import (
	"errors"
	"fmt"
	"strings"
)

type Book struct {
	Name   string
	Author string
	Year   int
	Read   bool
}

func main() {
	catal := []Book{}

	catal = Add(catal, Book{Name: "Война и мир",
		Author: "Л. Толстой", Year: 1800, Read: true})

	//List(catal)
	catal = Add(catal, Book{Name: "Идиот",
		Author: "Достоевский", Year: 1900, Read: false})
	List(catal)

	book, err := Search(catal, "Призрак")
	if err != nil {
		fmt.Println("Призрак:", err)
	} else {
		fmt.Println(book.Name, book.Author, book.Year, book.Read)
	}
	book2, err := Search(catal, "война и мир")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Книга есть:", book2.Name, book2.Author, book2.Year, book2.Read)
	}

}

func List(c []Book) {
	if len(c) == 0 {
		fmt.Println("Каталог пуст")
	}
	for _, v := range c {
		fmt.Println("Книга в каталоге:", v.Name, v.Author, v.Year, v.Read)
	}
}

func Add(c []Book, b Book) []Book {
	c = append(c, b)

	return c
}

func Search(c []Book, Name string) (Book, error) {

	for _, it := range c {
		if strings.EqualFold(it.Name, Name) {
			return it, nil
		}
	}
	return Book{}, errors.New("Такой книги нет")
}
