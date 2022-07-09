// Package stack is a tiny package implementing a generic stack.
package stack

// A Stack.
type Stack[T any] []T

// Push v unto the stack.
func (s *Stack[T]) Push(v T) {
	*s = append(*s, v)
}

// Pop returns the topmost element from the stack and removes it. Panics if the stack is empty.
func (s *Stack[T]) Pop() T {
	var v T
	*s, v = (*s)[:len(*s)-1], (*s)[len(*s)-1]
	return v
}

// Top returns the topmost element of the stack without removing it. Panics if the stack is empty.
func (s *Stack[T]) Top() T {
	return (*s)[len(*s)-1]
}

// Len returns the number of elements on the stack.
func (s *Stack[T]) Len() int {
	return len(*s)
}
