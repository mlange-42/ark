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

func equalEntities(a, b []Entity) bool {
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

func equalRelations(a, b []RelationID) bool {
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
