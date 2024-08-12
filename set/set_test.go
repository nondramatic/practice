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

func TestSet_AddSome(t *testing.T) {
	s := NewSet[int64]()
	s.AddSome(1, 2, 3)
	var expect int64 = 3
	if s.Len() != expect {
		t.Errorf("expect size: %d, got: %d", expect, s.Len())
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

func TestSet_RemoveSome(t *testing.T) {
	s := NewSet[int64]()
	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.RemoveSome(1, 2)
	var expect int64 = 1
	if s.Len() != expect {
		t.Errorf("expect size: %d, got: %d", expect, s.Len())
	}
}

func TestSet_IsEmpty(t *testing.T) {
	s := NewSet[int64]()
	var expect bool = true
	if s.IsEmpty() != expect {
		t.Errorf("expect %v, got: %v", expect, s.Len())
	}
}

func TestSet_Clear(t *testing.T) {
	s := NewSet[int64]()
	s.Add(1)
	var expect int64 = 1
	if s.Len() != expect {
		t.Errorf("expect %v, got: %v", expect, s.Len())
	}

	s.Clear()
	expect = 0
	if s.Len() != 0 {
		t.Errorf("expect %v, got: %v", expect, s.Len())
	}

	if len(s.Data) != 0 {
		t.Errorf("expect %v, got: %v", expect, len(s.Data))
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
	s1 := NewSet[int64]()
	s1.Add(1)
	s1.Add(2)

	s2 := NewSet[int64]()
	s2.Add(1)
	s2.Add(3)

	var expect bool = true
	if (s1.Intersection(s2)).Len() == 0 {
		t.Errorf("expect %v, got: %v", expect, false)
	}
}

func TestSet_Difference(t *testing.T) {
	s1 := NewSet[int64]()
	s1.Add(1)
	s1.Add(2)

	s2 := NewSet[int64]()
	s2.Add(1)
	s2.Add(3)

	var expect bool = true
	if (s1.Difference(s2)).Len() == 0 {
		t.Errorf("expect %v, got: %v", expect, false)
	}
}

func TestSet_SymmetricDifference(t *testing.T) {
	s1 := NewSet[int64]()
	s1.Add(1)
	s1.Add(2)

	s2 := NewSet[int64]()
	s2.Add(1)
	s2.Add(3)

	var expect bool = true
	if (s1.SymmetricDifference(s2)).Len() == 0 {
		t.Errorf("expect %v, got: %v", expect, false)
	}
}

func TestSet_Subset(t *testing.T) {
	s1 := NewSet[int64]()
	s1.Add(1)
	s1.Add(2)

	s2 := NewSet[int64]()
	s2.Add(1)
	s2.Add(2)
	s2.Add(3)

	var expect bool = true
	if !s1.Subset(s2) {
		t.Errorf("expect %v, got: %v", expect, false)
	}

	expect = false
	if s2.Subset(s1) {
		t.Errorf("expect %v, got: %v", expect, true)
	}
}

func TestSet_Superset(t *testing.T) {
	s1 := NewSet[int64]()
	s1.Add(1)
	s1.Add(2)

	s2 := NewSet[int64]()
	s2.Add(1)
	s2.Add(2)
	s2.Add(3)

	var expect bool = true
	if !s2.Superset(s1) {
		t.Errorf("expect %v, got: %v", expect, false)
	}
}

func BenchmarkSetAdd(b *testing.B) {
	s := NewSet[int]()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
}
