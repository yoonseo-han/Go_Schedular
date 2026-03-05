package runtime

type Stack struct {
	lo uintptr
	hi uintptr

	data []byte
}

func newStack() *Stack {
	return &Stack{
		lo:   0,
		hi:   0,
		data: make([]byte, 0),
	}
}

func (s *Stack) push(x uintptr) {
	s.data = append(s.data, byte(x))
}

func (s *Stack) pop() byte {
	x := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return x
}

func (s *Stack) peek() byte {
	return s.data[len(s.data)-1]
}

func (s *Stack) size() int {
	return len(s.data)
}

func (s *Stack) isEmpty() bool {
	return len(s.data) == 0
}
