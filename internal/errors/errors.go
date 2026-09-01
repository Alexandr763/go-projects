package errors

import (
	"fmt"
)

var ErrDivideByZero = fmt.Errorf("деление на ноль")

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}
	return a / b, nil
}
