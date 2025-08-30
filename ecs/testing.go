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
