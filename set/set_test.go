package set

import (
	"testing"
)

func TestSet_Add(t *testing.T) {
	s := NewSet[int64]()
	s.Add(1)
	s.Add(2)
	s.Add(3)
	if s.Size != 3 {
		t.Errorf("data: %v, size: %d \n", s.Data, s.Size)
	}

}

func TestSet_Remove(t *testing.T) {
	s := NewSet[int64]()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	if s.Size != 3 {
		t.Errorf("expect size: %d, got: %d", 3, s.Size)
	}
	s.Remove(2)
	if s.Size != 2 {
		t.Errorf("expect size: %d, got: %d", 2, s.Size)
	}
}

func TestSet_Contains(t *testing.T) {
	s := NewSet[int64]()
	s.Add(1)
	s.Add(2)
	if s.Contains(3) {
		t.Errorf("expect %v, got: %v", false, true)
	}

	if !s.Contains(2) {
		t.Errorf("expect %v, got: %v", true, false)
	}
}

func TestSet_Size(t *testing.T) {
	s := NewSet[int64]()
	s.Add(1)
	s.Add(2)
	var expect int64 = 2
	if s.Len() != expect {
		t.Errorf("expect %v, got: %v", expect, s.Len())
	}
}

func TestSet_Union(t *testing.T) {
	s1 := NewSet[int64]()
	s1.Add(1)
	s1.Add(2)

	s2 := NewSet[int64]()
	s2.Add(4)
	s2.Add(5)

	s3 := s1.Union(s2)
	var expect int64 = 4
	if s3.Len() != expect {
		t.Errorf("expect %v, got: %v", expect, s3.Len())
	}
}

func TestSet_Intersection(t *testing.T) {

}

func TestSet_Difference(t *testing.T) {

}

func TestSet_SymmetricDifference(t *testing.T) {

}

func TestSet_SubsetOf(t *testing.T) {

}

func TestSet_SymmetricSubsetOf(t *testing.T) {

}
