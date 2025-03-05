package ecs

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntity(t *testing.T) {
	e := newEntity(100)
	assert.EqualValues(t, 100, e.ID())
	assert.EqualValues(t, 0, e.Gen())
}

func TestEntityIndex(t *testing.T) {
	index := entityIndex{}
	assert.EqualValues(t, 0, index.table)
	assert.EqualValues(t, 0, index.row)
}

func TestReservedEntities(t *testing.T) {
	w := NewWorld(1024)

	zero := Entity{}
	wildcard := Entity{1, 0}

	assert.False(t, w.Alive(zero))
	assert.False(t, w.Alive(wildcard))

	assert.True(t, zero.IsZero())
	assert.False(t, wildcard.IsZero())
	assert.True(t, wildcard.isWildcard())
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

	assert.Equal(t, e2, e)

	err = e2.UnmarshalJSON([]byte("pft"))
	assert.NotNil(t, err)
}
