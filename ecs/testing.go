package ecs

import "testing"

func expectPanicWithValue(t *testing.T, expected interface{}, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != expected {
			t.Errorf("expected panic with %v, got %v", expected, r)
		}
	}()
	f()
}

func expectPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but none occurred")
		}
	}()
	f()
}

func expectEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func expectTrue(t *testing.T, cond bool) {
	t.Helper()
	if !cond {
		t.Errorf("expected condition to be true, but was false")
	}
}

func expectFalse(t *testing.T, cond bool) {
	t.Helper()
	if cond {
		t.Errorf("expected condition to be false, but was true")
	}
}

func equalSlices[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
