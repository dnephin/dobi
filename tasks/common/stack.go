package common

import (
	"fmt"
)

// Stack is a LIFO stack of strings
type Stack struct {
	data []TaskName
}

// Contains returns true if an item is in the stack, false otherwise
func (s *Stack) Contains(item TaskName) bool {
	for _, elem := range s.data {
		if elem.Equal(item) {
			return true
		}
	}
	return false
}

// Push adds an item to the top of the stack
func (s *Stack) Push(item TaskName) {
	s.data = append(s.data, item)
}

// Pop removes an item from the top of the stack
func (s *Stack) Pop() (TaskName, error) {
	var item TaskName
	last := len(s.data) - 1
	if last < 0 {
		return TaskName{}, fmt.Errorf("Can't pop empty stack.")
	}
	item, s.data = s.data[last], s.data[:last]
	return item, nil
}

// Reset removes all items from the stack
func (s *Stack) Reset() {
	s.data = []TaskName{}
}

// Items returns all the items in the stack in order
func (s *Stack) Items() []TaskName {
	return s.data
}

// Names returns all the name of the items in the stack in order
func (s *Stack) Names() []string {
	names := []string{}
	for _, taskName := range s.data {
		names = append(names, taskName.Name())
	}
	return names
}

// NewStack returns a new empty stack
func NewStack() *Stack {
	return &Stack{data: []TaskName{}}
}
