package generator

import (
	"math/rand"
	"time"
)

func Generate() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := ""
	for i := 0; i < 6; i++ {
		index := rand.Intn(len(chars))
		result += string(chars[index])
	}
	return result
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
