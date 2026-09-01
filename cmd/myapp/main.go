package main

import (
	"fmt"
	"go-projects/internal/hashtable"
	"go-projects/internal/linkedlist"
	"go-projects/internal/queue"
	"go-projects/internal/set"
	"go-projects/internal/stack"
)

func main() {
	list := &linkedlist.LinkedList{}

	list.PushFront(10)
	list.PushFront(20)
	list.PushFront(30)

	list.PushBack(40)
	list.PushBack(50)
	for n := list.Head; n != nil; n = n.Next {
		fmt.Print(n.Value, " ")
	}
	fmt.Println()

	// Проверка Insert
	err := list.Insert(2, 99) // вставка в середину
	if err != nil {
		fmt.Println("Ошибка Insert:", err)
	}

	// Проверка Insert в начало
	list.Insert(0, 1)

	// Проверка Insert в конец
	list.Insert(list.Len, 100)

	// Ещё раз напечатай список
	fmt.Println("После вставок:")
	for n := list.Head; n != nil; n = n.Next {
		fmt.Print(n.Value, " ")
	}
	fmt.Println()

	// Проверка Delete
	err = list.Delete(2) // удалить элемент на позиции 2
	if err != nil {
		fmt.Println("Ошибка Delete:", err)
	}

	fmt.Println("После удаления:")
	for n := list.Head; n != nil; n = n.Next {
		fmt.Print(n.Value, " ")
	}
	fmt.Println()
	// Проверка Find
	idx := list.Find(99)
	if idx != -1 {
		fmt.Println("Нашли 99 на позиции:", idx)
	} else {
		fmt.Println("99 не найдено")
	}

	idx = list.Find(999)
	if idx != -1 {
		fmt.Println("Нашли 999 на позиции:", idx)
	} else {
		fmt.Println("999 не найдено")
	}

	// Проверка Stack
	stack := &stack.Stack{}

	stack.Push(10)
	stack.Push(20)
	stack.Push(30)

	val, _ := stack.Pop()
	fmt.Println("Pop:", val) // 30

	val, _ = stack.Peek()
	fmt.Println("Peek:", val) // 20

	fmt.Println("IsEmpty:", stack.IsEmpty()) // false

	val, _ = stack.Pop()
	fmt.Println("Pop:", val) // 20

	val, _ = stack.Pop()
	fmt.Println("Pop:", val) // 10

	fmt.Println("IsEmpty:", stack.IsEmpty()) // true

	// Проверка Queue
	q := &queue.Queue{}
	q.Enqueue(10)
	q.Enqueue(20)
	q.Enqueue(30)

	valQ, _ := q.Dequeue()
	fmt.Println("Dequeue:", valQ) // 10

	valQ, _ = q.Peek()
	fmt.Println("Peek:", valQ) // 20

	fmt.Println("IsEmpty:", q.IsEmpty()) // false

	valQ, _ = q.Dequeue()
	fmt.Println("Dequeue:", valQ) // 20

	valQ, _ = q.Dequeue()
	fmt.Println("Dequeue:", valQ) // 30

	fmt.Println("IsEmpty:", q.IsEmpty()) // true

	// Проверка HashTable
	ht := hashtable.NewHashTable()

	ht.Set("Alice", 25)
	ht.Set("Bob", 30)

	age, _ := ht.Get("Alice")
	fmt.Println("Alice:", age)

	ht.Delete("Bob")
	fmt.Println("Has Bob:", ht.Has("Bob")) // false

	// Проверка Set
	s := set.NewSet()
	s.Add(10)
	s.Add(20)
	s.Add(10) // не добавится

	fmt.Println("Size:", s.Size())              // 2
	fmt.Println("Contains 10:", s.Contains(10)) // true
	fmt.Println("Contains 30:", s.Contains(30)) // false

	s.Remove(10)
	fmt.Println("Size after remove:", s.Size()) // 1
}
