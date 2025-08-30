package ecs

import (
	"encoding/json"
	"testing"
)

func TestEntity(t *testing.T) {
	e := newEntity(100)
	if e.ID() != 100 {
		t.Errorf("expected ID to be 100, got %d", e.ID())
	}
	if e.Gen() != 0 {
		t.Errorf("expected Gen to be 0, got %d", e.Gen())
	}
}

func TestEntityIndex(t *testing.T) {
	index := entityIndex{}
	if index.table != 0 {
		t.Errorf("expected table to be 0, got %d", index.table)
	}
	if index.row != 0 {
		t.Errorf("expected row to be 0, got %d", index.row)
	}
}

func TestReservedEntities(t *testing.T) {
	w := NewWorld(1024)

	zero := Entity{}
	wildcard := Entity{1, 0}

	if w.Alive(zero) {
		t.Errorf("expected zero entity to be not alive")
	}
	if w.Alive(wildcard) {
		t.Errorf("expected wildcard entity to be not alive")
	}

	if !zero.IsZero() {
		t.Errorf("expected zero.IsZero() to be true")
	}
	if wildcard.IsZero() {
		t.Errorf("expected wildcard.IsZero() to be false")
	}
	if !wildcard.isWildcard() {
		t.Errorf("expected wildcard.isWildcard() to be true")
	}
}

func TestEntityMarshal(t *testing.T) {
	e := Entity{2, 3}

	jsonData, err := json.Marshal(&e)
	if err != nil {
		t.Fatalf("unexpected error during marshal: %v", err)
	}

	e2 := Entity{}
	err = json.Unmarshal(jsonData, &e2)
	if err != nil {
		t.Fatalf("unexpected error during unmarshal: %v", err)
	}

	if e2 != e {
		t.Errorf("expected unmarshaled entity to equal original, got %+v vs %+v", e2, e)
	}

	err = e2.UnmarshalJSON([]byte("pft"))
	if err == nil {
		t.Errorf("expected error from UnmarshalJSON with invalid input, got nil")
	}
}
