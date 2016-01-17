package stack

import "fmt"

// StringStack is a LIFO stack of strings
type StringStack struct {
	data []string
}

// Contains returns true if an item is in the stack, false otherwise
func (s *StringStack) Contains(item string) bool {
	for _, elem := range s.data {
		if elem == item {
			return true
		}
	}
	return false
}

// Push adds an item to the top of the stack
func (s *StringStack) Push(item string) {
	s.data = append(s.data, item)
}

// Pop removes an item from the top of the stack
func (s *StringStack) Pop() (string, error) {
	var item string
	last := len(s.data) - 1
	if last < 0 {
		return "", fmt.Errorf("Can't pop empty stack.")
	}
	item, s.data = s.data[last], s.data[:last]
	return item, nil
}

// Reset removes all items from the stack
func (s *StringStack) Reset() {
	s.data = []string{}
}

// Items returns all the items in the stack in order
func (s *StringStack) Items() []string {
	return s.data
}

// NewStringStack returns a new empty stack
func NewStringStack() *StringStack {
	return &StringStack{data: []string{}}
}
