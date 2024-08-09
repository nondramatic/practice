package set

import (
	"fmt"
	"testing"
)

func TestNewSet(t *testing.T) {
	fmt.Println(NewSet[int64]())
	fmt.Println(NewSet[string]())
	fmt.Println(NewSet[float64]())
}

func TestSet_Add(t *testing.T) {
	s := NewSet[int64]()
	s.Add(1)
	s.Add(2)
	s.Add(3)
	fmt.Printf("datq: %v, size: %d \n", s.Data, s.Size)
}

func TestSet_Remove(t *testing.T) {
	s := NewSet[int64]()
	s.Add(1)
	s.Add(2)
	s.Add(3)
	fmt.Printf("datq: %v, size: %d \n", s.Data, s.Size)
	s.Remove(2)
	fmt.Printf("datq: %v, size: %d \n", s.Data, s.Size)
}
