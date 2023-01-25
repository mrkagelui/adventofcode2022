package main

import "errors"

type stack[T any] []T

func (s *stack[T]) push(v T) {
	*s = append(*s, v)
}

func (s *stack[T]) pop() (T, error) {
	if len(*s) == 0 {
		return *new(T), errors.New("empty stack")
	}
	v := []T(*s)[len(*s)-1]
	*s = []T(*s)[:len(*s)-1]
	return v, nil
}

func (s *stack[T]) peek() (T, error) {
	if len(*s) == 0 {
		return *new(T), errors.New("empty stack")
	}
	return []T(*s)[len(*s)-1], nil
}

func (s *stack[T]) isEmpty() bool {
	return len(*s) == 0
}
