package stack

import "errors"

type Stack struct {
	data []int
}

func (s *Stack) Push(value int) {
	s.data = append(s.data, value)
}

func (s *Stack) Pop() (int, error) {
	if len(s.data) == 0 {
		return 0, errors.New("стек пуст")
	}
	last := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return last, nil

}

func (s *Stack) Peek() (int, error) {
	if len(s.data) == 0 {
		return 0, errors.New("стек пуст")
	} else {
		last := s.data[len(s.data)-1]
		return last, nil
	}
}

func (s *Stack) IsEmpty() bool {
	if len(s.data) == 0 {
		return true
	}
	return false
}
