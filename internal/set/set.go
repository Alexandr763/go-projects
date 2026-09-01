package set

type Set struct {
	data map[int]struct{}
}

func NewSet() *Set {
	return &Set{
		data: make(map[int]struct{}),
	}
}

func (s *Set) Add(value int) {
	s.data[value] = struct{}{}
}

func (s *Set) Remove(value int) {
	delete(s.data, value)
}

func (s *Set) Contains(value int) bool {
	_, ok := s.data[value]
	return ok
}

func (s *Set) Size() int {
	return len(s.data)
}
