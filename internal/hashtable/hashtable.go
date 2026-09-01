package hashtable

import (
	"errors"
)

type HashTable struct {
	data map[string]int
}

func NewHashTable() *HashTable {
	return &HashTable{
		data: make(map[string]int),
	}
}

func (h *HashTable) Set(key string, value int) {
	h.data[key] = value
}

func (h *HashTable) Get(key string) (int, error) {
	value, ok := h.data[key]
	if ok {
		return value, nil
	}
	return 0, errors.New("ключ не найден")
}

func (h *HashTable) Delete(key string) {
	delete(h.data, key)
}

func (h *HashTable) Has(key string) bool {
	_, ok := h.data[key]
	return ok
}
