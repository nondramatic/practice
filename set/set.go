package set

import "sync"

type Set[T comparable] struct {
	mu   sync.RWMutex
	size int64
	data map[T]struct{}
}

// NewSet initializes and returns a new Set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{data: make(map[T]struct{})}
}

// Add inserts a new element into the Set
func (s *Set[T]) Add(element T) (bool, int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.data[element]; !exists {
		s.data[element] = struct{}{}
		s.size++
		return true, s.size
	}

	return false, s.size
}

// AddSome inserts some elements into the Set
func (s *Set[T]) AddSome(elements ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, element := range elements {
		if _, exists := s.data[element]; !exists {
			s.data[element] = struct{}{}
			s.size++
		}
	}
}

// RemoveSome deletes some elements from the Set
func (s *Set[T]) RemoveSome(elements ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, element := range elements {
		if _, exists := s.data[element]; exists {
			delete(s.data, element)
			s.size--
		}
	}
}

// Remove deletes an element from the Set
func (s *Set[T]) Remove(element T) (bool, int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.data[element]; exists {
		delete(s.data, element)
		s.size--
		return true, s.size
	}

	return false, s.size
}

// Contains checks if an element exists in the Set
func (s *Set[T]) Contains(element T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.data[element]
	return exists
}

// Len returns the number of elements in the Set
func (s *Set[T]) Len() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.size
}

// IsEmpty checks if the Set is empty
func (s *Set[T]) IsEmpty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.size == 0
}

// Clear removes all elements from the Set
func (s *Set[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[T]struct{})
	s.size = 0
}

// Union returns a new Set with elements from both Sets
func (s *Set[T]) Union(o *Set[T]) *Set[T] {
	unionSet := NewSet[T]()
	for k := range s.data {
		unionSet.Add(k)
	}
	for k := range o.data {
		unionSet.Add(k)
	}
	return unionSet
}

// Intersection returns a new Set with elements common to both Sets
func (s *Set[T]) Intersection(o *Set[T]) *Set[T] {
	intersectionSet := NewSet[T]()
	for k := range s.data {
		if o.Contains(k) {
			intersectionSet.Add(k)
		}
	}
	return intersectionSet
}

// Difference returns a new Set with elements in the current Set but not in the other
func (s *Set[T]) Difference(o *Set[T]) *Set[T] {
	differenceSet := NewSet[T]()
	for k := range s.data {
		if !o.Contains(k) {
			differenceSet.Add(k)
		}
	}
	return differenceSet
}

// SymmetricDifference returns a new Set with elements in either Set, but not both
func (s *Set[T]) SymmetricDifference(o *Set[T]) *Set[T] {
	symmetricDifferenceSet := NewSet[T]()
	for k := range s.data {
		if !o.Contains(k) {
			symmetricDifferenceSet.Add(k)
		}
	}
	for k := range o.data {
		if !s.Contains(k) {
			symmetricDifferenceSet.Add(k)
		}
	}
	return symmetricDifferenceSet
}

// Subset checks if the Set is a subset of another Set
func (s *Set[T]) Subset(o *Set[T]) bool {

	if s.size > o.size {
		return false
	}

	for k := range s.data {
		if !o.Contains(k) {
			return false
		}
	}
	return true
}

// Superset checks if the Set is a superset of another Set
func (s *Set[T]) Superset(o *Set[T]) bool {
	return o.Subset(s)
}
