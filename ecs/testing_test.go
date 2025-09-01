package ecs

import (
	"fmt"
	"reflect"
	"testing"
)

func expectEqual[T comparable](t *testing.T, want, got T, msgAndArgs ...interface{}) {
	t.Helper()
	if got != want {
		base := fmt.Sprintf("expected %v, got %v", want, got)
		t.Error(base + formatMsg(msgAndArgs...))
	}
}

func expectNotEqual[T comparable](t *testing.T, notWant, got T, msgAndArgs ...interface{}) {
	t.Helper()
	if got == notWant {
		base := fmt.Sprintf("expected value not equal to %v, but got %v", notWant, got)
		t.Error(base + formatMsg(msgAndArgs...))
	}
}

func expectGreater[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](t *testing.T, got, min T, msgAndArgs ...interface{}) {
	t.Helper()
	if got <= min {
		base := fmt.Sprintf("expected value > %v, got %v", min, got)
		t.Error(base + formatMsg(msgAndArgs...))
	}
}

func expectLess[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](t *testing.T, got, max T, msgAndArgs ...interface{}) {
	t.Helper()
	if got >= max {
		base := fmt.Sprintf("expected value < %v, got %v", max, got)
		t.Error(base + formatMsg(msgAndArgs...))
	}
}

func expectTrue(t *testing.T, cond bool, msgAndArgs ...interface{}) {
	t.Helper()
	if !cond {
		base := "expected condition to be true, but was false"
		t.Error(base + formatMsg(msgAndArgs...))
	}
}

func expectFalse(t *testing.T, cond bool, msgAndArgs ...interface{}) {
	t.Helper()
	if cond {
		base := "expected condition to be false, but was true"
		t.Error(base + formatMsg(msgAndArgs...))
	}
}

func expectNil(t *testing.T, val interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if val != nil && !isReallyNil(val) {
		base := fmt.Sprintf("expected nil, got %v", val)
		t.Error(base + formatMsg(msgAndArgs...))
	}
}

func expectNotNil(t *testing.T, val interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if val == nil || isReallyNil(val) {
		base := "expected non-nil value, but got nil"
		t.Error(base + formatMsg(msgAndArgs...))
	}
}

func expectSlicesEqual[T comparable](t *testing.T, want, got []T, msgAndArgs ...interface{}) {
	t.Helper()
	if len(got) != len(want) {
		base := fmt.Sprintf("slice length mismatch: expected %d, got %d", len(want), len(got))
		t.Error(base + formatMsg(msgAndArgs...))
		return
	}
	for i := range got {
		if got[i] != want[i] {
			base := fmt.Sprintf("slice mismatch at index %d: expected %v, got %v", i, want[i], got[i])
			t.Error(base + formatMsg(msgAndArgs...))
			return
		}
	}
}

func expectPanicsWithValue(t *testing.T, expected interface{}, f func(), msgAndArgs ...interface{}) {
	t.Helper()
	defer func() {
		if r := recover(); r != expected {
			base := fmt.Sprintf("expected panic with %v, got %v", expected, r)
			t.Error(base + formatMsg(msgAndArgs...))
		}
	}()
	f()
}

func expectPanics(t *testing.T, f func(), msgAndArgs ...interface{}) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			base := "expected panic, but none occurred"
			t.Error(base + formatMsg(msgAndArgs...))
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

func formatMsg(msgAndArgs ...interface{}) string {
	msg := messageFromMsgAndArgs(msgAndArgs...)
	if msg == "" {
		return ""
	}
	return ": " + msg
}

func messageFromMsgAndArgs(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			return msgAsStr
		}
		return fmt.Sprintf("%+v", msg)
	}
	if len(msgAndArgs) > 1 {
		if format, ok := msgAndArgs[0].(string); ok {
			return fmt.Sprintf(format, msgAndArgs[1:]...)
		}
		return fmt.Sprintf("%+v", msgAndArgs...)
	}
	return ""
}
