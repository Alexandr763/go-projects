package storage

import (
	"errors"
	//"fmt"
	"url-shortener/internal/generator"
)

type Storage struct {
	data map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func (s *Storage) Save(url string) (string, error) {
	if url == "" {
		return "", errors.New("Пустой URL")
	}
	for {
		code := generator.Generate()
		_, ok := s.data[code]
		if !ok {
			s.data[code] = url
			return code, nil
		}
	}

}

func (s *Storage) Get(code string) (string, error) {
	value, ok := s.data[code]
	if ok != true {
		return "", errors.New("Такого кода нет")
	}
	return value, nil
}

func (s *Storage) Delete(code string) {
	delete(s.data, code)
}
