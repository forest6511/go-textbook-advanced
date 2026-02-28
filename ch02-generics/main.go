package main

import (
	"fmt"
)

// MapKeys returns all keys of a map
func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Stack is a generic type-safe stack
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	n := len(s.items) - 1
	item := s.items[n]
	s.items = s.items[:n]
	return item, true
}

// Adder demonstrates self-referential generic constraints (Go 1.26)
type Adder[A Adder[A]] interface {
	Add(A) A
}

func algo[A Adder[A]](x, y A) A {
	return x.Add(y)
}

// Vec2 implements Adder[Vec2]
type Vec2 struct{ X, Y float64 }

func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{v.X + other.X, v.Y + other.Y}
}

func main() {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := MapKeys(m)
	fmt.Println("keys count:", len(keys))

	s := &Stack[int]{}
	s.Push(1)
	s.Push(2)
	s.Push(3)
	if v, ok := s.Pop(); ok {
		fmt.Println("popped:", v)
	}

	a := Vec2{1, 2}
	b := Vec2{3, 4}
	result := algo(a, b)
	fmt.Printf("Vec2 add: {%.1f, %.1f}\n", result.X, result.Y)
}
