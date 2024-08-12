package set

import "sync"

type Set[T comparable] struct {
	sync.Mutex
	Size int64
	Data map[T]struct{}
}

// NewSet initializes and returns a new Set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{Data: make(map[T]struct{})}
}

// Add inserts a new element into the Set
func (s *Set[T]) Add(element T) (bool, int64) {
	s.Lock()
	defer s.Unlock()
	if _, exists := s.Data[element]; !exists {
		s.Data[element] = struct{}{}
		s.Size++
		return true, s.Size
	}

	return false, s.Size
}

// AddSome inserts some elements into the Set
func (s *Set[T]) AddSome(elements ...T) {
	s.Lock()
	defer s.Unlock()
	for _, element := range elements {
		if _, exists := s.Data[element]; !exists {
			s.Data[element] = struct{}{}
			s.Size++
		}
	}
}

// RemoveSome deletes some elements from the Set
func (s *Set[T]) RemoveSome(elements ...T) {
	s.Lock()
	defer s.Unlock()
	for _, element := range elements {
		if _, exists := s.Data[element]; exists {
			delete(s.Data, element)
			s.Size--
		}
	}
}

// Remove deletes an element from the Set
func (s *Set[T]) Remove(element T) (bool, int64) {
	s.Lock()
	defer s.Unlock()
	if _, exists := s.Data[element]; exists {
		delete(s.Data, element)
		s.Size--
		return true, s.Size
	}

	return false, s.Size
}

// Contains checks if an element exists in the Set
func (s *Set[T]) Contains(element T) bool {
	s.Lock()
	defer s.Unlock()
	_, exists := s.Data[element]
	return exists
}

// Len returns the number of elements in the Set
func (s *Set[T]) Len() int64 {
	return s.Size
}

// IsEmpty checks if the Set is empty
func (s *Set[T]) IsEmpty() bool {
	return s.Size == 0
}

// Clear removes all elements from the Set
func (s *Set[T]) Clear() {
	s.Lock()
	defer s.Unlock()
	s.Data = make(map[T]struct{})
	s.Size = 0
}

// Union returns a new Set with elements from both Sets
func (s *Set[T]) Union(o *Set[T]) *Set[T] {
	unionSet := NewSet[T]()
	for k := range s.Data {
		unionSet.Add(k)
	}
	for k := range o.Data {
		unionSet.Add(k)
	}
	return unionSet
}

// Intersection returns a new Set with elements common to both Sets
func (s *Set[T]) Intersection(o *Set[T]) *Set[T] {
	intersectionSet := NewSet[T]()
	for k := range s.Data {
		if o.Contains(k) {
			intersectionSet.Add(k)
		}
	}
	return intersectionSet
}

// Difference returns a new Set with elements in the current Set but not in the other
func (s *Set[T]) Difference(o *Set[T]) *Set[T] {
	differenceSet := NewSet[T]()
	for k := range s.Data {
		if !o.Contains(k) {
			differenceSet.Add(k)
		}
	}
	return differenceSet
}

// SymmetricDifference returns a new Set with elements in either Set, but not both
func (s *Set[T]) SymmetricDifference(o *Set[T]) *Set[T] {
	symmetricDifferenceSet := NewSet[T]()
	for k := range s.Data {
		if !o.Contains(k) {
			symmetricDifferenceSet.Add(k)
		}
	}
	for k := range o.Data {
		if !s.Contains(k) {
			symmetricDifferenceSet.Add(k)
		}
	}
	return symmetricDifferenceSet
}

// Subset checks if the Set is a subset of another Set
func (s *Set[T]) Subset(o *Set[T]) bool {
	for k := range s.Data {
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
