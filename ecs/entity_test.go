package ecs

import (
	"encoding/json"
	"testing"
)

func TestEntity(t *testing.T) {
	e := newEntity(100)
	expectEqual(t, 100, e.ID())
	expectEqual(t, 0, e.Gen())
}

func TestEntityIndex(t *testing.T) {
	index := entityIndex{}
	expectEqual(t, 0, index.table)
	expectEqual(t, 0, index.row)
}

func TestReservedEntities(t *testing.T) {
	w := NewWorld(1024)

	zero := Entity{}
	wildcard := Entity{1, 0}

	expectFalse(t, w.Alive(zero))
	expectFalse(t, w.Alive(wildcard))

	expectTrue(t, zero.IsZero())
	expectFalse(t, wildcard.IsZero())
	expectTrue(t, wildcard.isWildcard())
}

func TestEntityMarshal(t *testing.T) {
	e := Entity{2, 3}

	jsonData, err := json.Marshal(&e)
	if err != nil {
		t.Fatal(err)
	}

	e2 := Entity{}
	err = json.Unmarshal(jsonData, &e2)
	if err != nil {
		t.Fatal(err)
	}

	expectEqual(t, e2, e)

	err = e2.UnmarshalJSON([]byte("pft"))
	expectNotNil(t, err)
}
