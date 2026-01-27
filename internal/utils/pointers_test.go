package pointers

import "testing"

func TestTo(t *testing.T) {
	// int
	x := 42
	ptr := To(x)
	if *ptr != x {
		t.Errorf("expected %v, got %v", x, *ptr)
	}

	// string
	s := "hello"
	ptrS := To(s)
	if *ptrS != s {
		t.Errorf("expected %q, got %q", s, *ptrS)
	}

	// struct
	type foo struct{ A int }
	f := foo{A: 7}
	ptrF := To(f)
	if *ptrF != f {
		t.Errorf("expected %+v, got %+v", f, *ptrF)
	}
}
