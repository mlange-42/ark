package ecs

import (
	"fmt"
	"reflect"
	"testing"
)

func expectEqual[T comparable](t *testing.T, want, got T, msg ...any) {
	t.Helper()
	if got != want {
		base := fmt.Sprintf("expected %v, got %v", want, got)
		if len(msg) > 0 {
			t.Error(base + ": " + fmt.Sprint(msg...))
		} else {
			t.Error(base)
		}
	}
}

func expectNotEqual[T comparable](t *testing.T, notWant, got T, msg ...any) {
	t.Helper()
	if got == notWant {
		base := fmt.Sprintf("expected value not equal to %v, but got %v", notWant, got)
		if len(msg) > 0 {
			t.Error(base + ": " + fmt.Sprint(msg...))
		} else {
			t.Error(base)
		}
	}
}

func expectGreater[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](t *testing.T, got, min T, msg ...any) {
	t.Helper()
	if got <= min {
		base := fmt.Sprintf("expected value > %v, got %v", min, got)
		if len(msg) > 0 {
			t.Error(base + ": " + fmt.Sprint(msg...))
		} else {
			t.Error(base)
		}
	}
}

func expectTrue(t *testing.T, cond bool, msg ...any) {
	t.Helper()
	if !cond {
		base := "expected condition to be true, but was false"
		if len(msg) > 0 {
			t.Error(base + ": " + fmt.Sprint(msg...))
		} else {
			t.Error(base)
		}
	}
}

func expectFalse(t *testing.T, cond bool, msg ...any) {
	t.Helper()
	if cond {
		base := "expected condition to be false, but was true"
		if len(msg) > 0 {
			t.Error(base + ": " + fmt.Sprint(msg...))
		} else {
			t.Error(base)
		}
	}
}

func expectNil(t *testing.T, val interface{}, msg ...any) {
	t.Helper()
	if val != nil && !isReallyNil(val) {
		base := fmt.Sprintf("expected nil, got %v", val)
		if len(msg) > 0 {
			t.Error(base + ": " + fmt.Sprint(msg...))
		} else {
			t.Error(base)
		}
	}
}

func expectNotNil(t *testing.T, val interface{}, msg ...any) {
	t.Helper()
	if val == nil || isReallyNil(val) {
		base := "expected non-nil value, but got nil"
		if len(msg) > 0 {
			t.Error(base + ": " + fmt.Sprint(msg...))
		} else {
			t.Error(base)
		}
	}
}

func expectSlicesEqual[T comparable](t *testing.T, want, got []T, msg ...any) {
	t.Helper()
	if len(got) != len(want) {
		base := fmt.Sprintf("slice length mismatch: expected %d, got %d", len(want), len(got))
		if len(msg) > 0 {
			t.Error(base + ": " + fmt.Sprint(msg...))
		} else {
			t.Error(base)
		}
		return
	}
	for i := range got {
		if got[i] != want[i] {
			base := fmt.Sprintf("slice mismatch at index %d: expected %v, got %v", i, want[i], got[i])
			if len(msg) > 0 {
				t.Error(base + ": " + fmt.Sprint(msg...))
			} else {
				t.Error(base)
			}
			return
		}
	}
}

func expectPanicsWithValue(t *testing.T, expected interface{}, f func(), msg ...any) {
	t.Helper()
	defer func() {
		if r := recover(); r != expected {
			base := fmt.Sprintf("expected panic with %v, got %v", expected, r)
			if len(msg) > 0 {
				t.Error(base + ": " + fmt.Sprint(msg...))
			} else {
				t.Error(base)
			}
		}
	}()
	f()
}

func expectPanics(t *testing.T, f func(), msg ...any) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			base := "expected panic, but none occurred"
			if len(msg) > 0 {
				t.Error(base + ": " + fmt.Sprint(msg...))
			} else {
				t.Error(base)
			}
		}
	}()
	f()
}

func isReallyNil(val interface{}) bool {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}
