package queue

import (
	"errors"
)

type Queue struct {
	data []int
}

func (q *Queue) Enqueue(value int) {
	q.data = append(q.data, value)
}

func (q *Queue) Dequeue() (int, error) {
	if len(q.data) == 0 {
		return 0, errors.New("Пустая очередь")
	}
	first := q.data[0]
	q.data = q.data[1:]
	return first, nil

}

func (q *Queue) Peek() (int, error) {
	if len(q.data) == 0 {
		return 0, errors.New("Пустая очередь")
	} else {
		return q.data[0], nil
	}

}

func (q *Queue) IsEmpty() bool {
	return len(q.data) == 0
}
