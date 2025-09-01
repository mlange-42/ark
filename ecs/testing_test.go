package ecs

import (
	"fmt"
	"testing"
)

func expectEqual[T comparable](t *testing.T, got, want T, msg ...any) {
	t.Helper()
	if got != want {
		if len(msg) > 0 {
			t.Errorf("expected %v, got %v: %s", want, got, fmt.Sprint(msg...))
		} else {
			t.Errorf("expected %v, got %v", want, got)
		}
	}
}

func expectTrue(t *testing.T, cond bool, msg ...any) {
	t.Helper()
	if !cond {
		if len(msg) > 0 {
			t.Errorf("expected condition to be true, but was false: %s", fmt.Sprint(msg...))
		} else {
			t.Errorf("expected condition to be true, but was false")
		}
	}
}

func expectFalse(t *testing.T, cond bool, msg ...any) {
	t.Helper()
	if cond {
		if len(msg) > 0 {
			t.Errorf("expected condition to be false, but was true: %s", fmt.Sprint(msg...))
		} else {
			t.Errorf("expected condition to be false, but was true")
		}
	}
}

func expectSlicesEqual[T comparable](t *testing.T, got, want []T, msg ...any) {
	t.Helper()
	if len(got) != len(want) {
		if len(msg) > 0 {
			t.Errorf("slice length mismatch: expected %d, got %d: %s", len(want), len(got), fmt.Sprint(msg...))
		} else {
			t.Errorf("slice length mismatch: expected %d, got %d", len(want), len(got))
		}
		return
	}
	for i := range got {
		if got[i] != want[i] {
			if len(msg) > 0 {
				t.Errorf("slice mismatch at index %d: expected %v, got %v: %s", i, want[i], got[i], fmt.Sprint(msg...))
			} else {
				t.Errorf("slice mismatch at index %d: expected %v, got %v", i, want[i], got[i])
			}
			return
		}
	}
}

func expectPanicWithValue(t *testing.T, expected interface{}, f func(), msg ...any) {
	t.Helper()
	defer func() {
		if r := recover(); r != expected {
			if len(msg) > 0 {
				t.Errorf("expected panic with %v, got %v: %s", expected, r, fmt.Sprint(msg...))
			} else {
				t.Errorf("expected panic with %v, got %v", expected, r)
			}
		}
	}()
	f()
}

func expectPanic(t *testing.T, f func(), msg ...any) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			if len(msg) > 0 {
				t.Errorf("expected panic, but none occurred: %s", fmt.Sprint(msg...))
			} else {
				t.Errorf("expected panic, but none occurred")
			}
		}
	}()
	f()
}
