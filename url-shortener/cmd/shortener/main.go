package main

import (
	"context"
	//"errors"
	"fmt"
	"time"
	"url-shortener/internal/storage"
)

func main() {
	s := storage.NewStorage()
	ctx := context.Background()

	for {
		var choice int

		fmt.Println("Выберите действие:")
		fmt.Println(1, "Создать короткую ссылку")
		fmt.Println(2, "Получить оригинальный URL")
		fmt.Println(3, "Удалить ссылку")
		fmt.Println(4, "Создать группу ссылок")
		fmt.Println(5, "Выйти")

		fmt.Scan(&choice)

		if choice == 5 {
			fmt.Println("Программа завершена.")
			break
		}

		if choice == 1 {
			fmt.Println("Введите URL:")
			var url string
			fmt.Scan(&url)

			code, err := s.Save(ctx, url)
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

		if choice == 4 {
			var urls []string

			fmt.Println("Введите URL по одному (или exit для завершения):")
			for {
				var url string
				fmt.Scan(&url)

				if url == "exit" {
					break
				}
				urls = append(urls, url)

			}

			if len(urls) == 0 {
				fmt.Println("Нет URL для обработки")
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			results := s.CreateBatch(ctx, urls)

			fmt.Println("Результаты:")
			for _, r := range results {
				if r.Err != nil {
					fmt.Printf("URL: %s — Ошибка: %v\n", r.URL, r.Err)
				} else {
					fmt.Printf("URL: %s — Код: %s\n", r.URL, r.Code)
				}
			}
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
