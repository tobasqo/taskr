package main

import (
	"fmt"
	"iter"
	"runtime"
)

type Set[T comparable] struct {
	elements map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{make(map[T]struct{})}
}

func NewSetFrom[T comparable, U any](values iter.Seq[U], getter func(U) T) *Set[T] {
	set := NewSet[T]()

	for value := range values {
		set.Add(getter(value))
	}

	return set
}

func (s *Set[T]) Contains(elem T) bool {
	_, exists := s.elements[elem]
	return exists
}

func (s *Set[T]) Add(elem T) {
	s.elements[elem] = struct{}{}
}

func (s *Set[T]) Items() iter.Seq[T] {
	return func(yield func(T) bool) {
		for elem := range s.elements {
			if !yield(elem) {
				return
			}
		}
	}
}

func (s Set[T]) Len() int {
	return len(s.elements)
}

func (s1 Set[T]) Difference(s2 Set[T]) *Set[T] {
	diff := NewSet[T]()

	for elem := range s1.elements {
		if !s2.Contains(elem) {
			diff.Add(elem)
		}
	}

	return diff
}

func SliceToSeq[E any](slice []E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, v := range slice {
			if !yield(v) {
				return
			}
		}
	}
}

func MapKeysToSeq[K comparable, V any](m map[K]V) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}

func MapValuesToSeq[K comparable, V any](m map[K]V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m {
			if !yield(v) {
				return
			}
		}
	}
}

func NotImplemented() {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Println("could not recover runtime information")
		return
	}

	funcName := runtime.FuncForPC(pc).Name()

	panic(fmt.Sprintf("%s:%d `%s` not implemented", file, line, funcName))
}
