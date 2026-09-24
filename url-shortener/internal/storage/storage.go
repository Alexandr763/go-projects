package storage

import (
	"context"
	"errors"
	"sync"
	"url-shortener/internal/generator"
)

type Storage struct {
	data map[string]string
	mu   sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func (s *Storage) Save(ctx context.Context, url string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	if url == "" {
		return "", errors.New("Пустой URL")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

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
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.data[code]
	if ok != true {
		return "", errors.New("Такого кода нет")
	}
	return value, nil

}

func (s *Storage) Delete(code string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, code)
}
