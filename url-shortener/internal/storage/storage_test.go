package storage

import (
	"context"
	"testing"
)

// Тест 1 проверка сохранения и получения
func TestSaveAndGet(t *testing.T) {
	a := NewStorage()
	ctx := context.Background()
	code, err := a.Save(ctx, "https://example.com")
	if err != nil {
		t.Fatalf("Не удалось сохранить: %v", err)
	}
	url, err := a.Get(code)
	if err != nil {
		t.Fatalf("сообщение: %v", err)
	}
	if url != "https://example.com" {
		t.Errorf("ожидалось %q, получено %q", "https://example.com", url)
	}
}

// Тест 2 проверка удаления
func TestDelete(t *testing.T) {
	ctx := context.Background()

	a := NewStorage() // вызвал конструктор
	code, err := a.Save(ctx, "https://example.com")
	if err != nil {
		t.Fatalf("Не удалось сохранить: %v", err)
	}
	a.Delete(code)
	_, err = a.Get(code)
	if err == nil {
		t.Fatalf("Ожидалась ошибка удаления")
	}
}

// Тест 3 проверка получения после удаления
func TestGetAfterDelete(t *testing.T) {
	ctx := context.Background()

	a := NewStorage() // вызвал конструктор
	code, err := a.Save(ctx, "https://example.com")
	if err != nil {
		t.Fatalf("Не удалось сохранить: %v", err)
	}
	a.Delete(code)
	url, err := a.Get(code)
	if err == nil {
		t.Fatalf("Ожидалась ошибка удаления")
	}
	if url != "" {
		t.Errorf("ожидался пустой URL, получен: %q", url)
	}
}

// Тест 4 проверка что кода не существует
func TestGetNotFound(t *testing.T) {
	a := NewStorage() // вызвал конструктор
	_, err := a.Get("несуществующий код")
	if err == nil {
		t.Fatal("Ожидалась ошибка несуществующего кода")
	}
}

// Тест 5 проверка на несколько ссылок
func TestMultipleLinks(t *testing.T) {
	ctx := context.Background()

	a := NewStorage() // вызвал конструктор
	code1, err := a.Save(ctx, "https://example.com")
	if err != nil {
		t.Fatalf("Не удалось сохранить: %v", err)
	}
	code2, err := a.Save(ctx, "https://google.com")
	if err != nil {
		t.Fatalf("Не удалось сохранить: %v", err)
	}
	code3, err := a.Save(ctx, "https://github.com")
	if err != nil {
		t.Fatalf("Не удалось сохранить: %v", err)
	}

	url1, err := a.Get(code1)
	if err != nil {
		t.Fatalf("не удалось получить: %v", err)
	}
	if url1 != "https://example.com" {
		t.Errorf("ожидался  URL, получен: %q", url1)
	}

	url2, err := a.Get(code2)
	if err != nil {
		t.Fatalf("не удалось получить: %v", err)
	}
	if url2 != "https://google.com" {
		t.Errorf("ожидался  URL, получен: %q", url2)
	}

	url3, err := a.Get(code3)
	if err != nil {
		t.Fatalf("не удалось получить: %v", err)
	}
	if url3 != "https://github.com" {
		t.Errorf("ожидался  URL, получен: %q", url3)
	}

}

// Тест 6 проверка на пустой url
func TestSaveEmptyURL(t *testing.T) {
	ctx := context.Background()

	a := NewStorage() // вызвал конструктор
	_, err := a.Save(ctx, "")
	if err == nil {
		t.Errorf("код пустой: %v", err)
	}

}
