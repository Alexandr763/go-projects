package main

import (
	//"errors"
	"fmt"
	"url-shortener/internal/storage"
)

func main() {
	s := storage.NewStorage()

	for {
		var choice int

		fmt.Println("Выберите действие:")
		fmt.Println(1, "Создать короткую ссылку")
		fmt.Println(2, "Получить оригинальный URL")
		fmt.Println(3, "Удалить ссылку")
		fmt.Println(4, "Выйти")

		fmt.Scan(&choice)

		if choice == 4 {
			fmt.Println("Программа завершена.")
			break
		}

		if choice == 1 {
			fmt.Println("Введите URL:")
			var url string
			fmt.Scan(&url)

			code, err := s.Save(url)
			if err != nil {
				fmt.Println("Ошибка")
			}
			fmt.Println("Короткая ссылка создана:", code)
		}

		if choice == 2 {
			var code string

			fmt.Println("Введите код который ищете:")
			fmt.Scan(&code)
			url2, err2 := s.Get(code)
			if err2 != nil {
				fmt.Println("Ошибка", err2)
			} else {
				fmt.Println("URL:", url2)
			}

		}

		if choice == 3 {
			fmt.Println("Введите код который хотите удалить:")
			var code string

			fmt.Scan(&code)
			s.Delete(code)
			fmt.Println("Ссылка удалена")

		}

	}
}

/*code, err := s.Save("https://example.com")

	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Сохранено с кодом:", code)
	}

	url, err := s.Get(code)
	if err != nil {
		fmt.Println("Ошибка", err)
	} else {
		fmt.Println("URL:", url)
	}

	s.Delete(code)
	url2, err2 := s.Get(code)
	if err2 != nil {
		fmt.Println("Ошибка", err2)
	} else {
		fmt.Println("URL:", url2)
	}

}
*/
