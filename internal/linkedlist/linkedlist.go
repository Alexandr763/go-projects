package linkedlist

import "errors"

type Node struct {
	Value int
	Next  *Node
}

type LinkedList struct {
	Head *Node
	Len  int
}

func (l *LinkedList) PushFront(value int) {
	newNode := &Node{
		Value: value,
		Next:  l.Head,
	}
	l.Head = newNode
	l.Len++
}

func (l *LinkedList) PushBack(value int) {
	newNode := &Node{
		Value: value,
		Next:  nil,
	}

	if l.Head == nil {
		l.Head = newNode
		l.Len++
		return
	}

	current := l.Head
	for current.Next != nil {
		current = current.Next
	}
	current.Next = newNode
	l.Len++
}

func (l *LinkedList) Insert(index int, value int) error {
	if index < 0 || index > l.Len {
		return errors.New("индекс вне диапазона")
	}

	if index == 0 {
		l.PushFront(value)
		return nil
	}

	if index == l.Len {
		l.PushBack(value)
		return nil
	}

	// вставка в середину
	current := l.Head
	for i := 0; i < index-1; i++ {
		current = current.Next
	}

	newNode := &Node{
		Value: value,
		Next:  current.Next,
	}
	current.Next = newNode
	l.Len++

	return nil
}

func (l *LinkedList) Delete(index int) error {
    if index < 0 || index >= l.Len {
        return errors.New("индекс вне диапазона")
    }

    if index == 0 {
        l.Head = l.Head.Next
        l.Len--
        return nil
    }

    current := l.Head
    for i := 0; i < index-1; i++ {
        current = current.Next
    }
    current.Next = current.Next.Next
    l.Len--
    return nil
}

func (l *LinkedList) Find(value int) int {
    current := l.Head
    index := 0
    for current != nil {
        if current.Value == value {
            return index
        }
        current = current.Next
        index++
    }
    return -1
}
