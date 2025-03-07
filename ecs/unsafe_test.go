package ecs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsafe(t *testing.T) {
	w := NewWorld(1024)
	u := w.Unsafe()

	assert.Equal(t, &w, u.world)
}

func TestUnsafeNewEntity(t *testing.T) {
	w := NewWorld(16)
	u := w.Unsafe()

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)

	e := u.NewEntity(posID, velID)

	assert.True(t, u.Has(e, posID))
	assert.True(t, u.Has(e, velID))
}

func TestUnsafeGet(t *testing.T) {
	w := NewWorld(16)
	u := w.Unsafe()

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)

	e := u.NewEntity(posID)

	assert.True(t, u.Has(e, posID))
	assert.False(t, u.Has(e, velID))

	assert.True(t, u.HasUnchecked(e, posID))
	assert.False(t, u.HasUnchecked(e, velID))

	pos := (*Position)(u.Get(e, posID))
	pos.X = 100

	pos2 := (*Position)(u.GetUnchecked(e, posID))

	assert.Equal(t, pos, pos2)

	assert.Panics(t, func() {
		u.Get(Entity{}, posID)
	})
	assert.Panics(t, func() {
		u.Has(Entity{}, posID)
	})
}

func TestUnsafeRelations(t *testing.T) {
	w := NewWorld(16)
	u := w.Unsafe()

	posID := ComponentID[Position](&w)
	childID := ComponentID[ChildOf](&w)
	child2ID := ComponentID[ChildOf2](&w)

	parent1 := w.NewEntity()
	parent2 := w.NewEntity()

	e := u.NewEntityRel([]ID{posID, childID, child2ID}, RelID(childID, parent1), RelID(child2ID, parent2))

	assert.Equal(t, parent1, u.GetRelation(e, childID))
	assert.Equal(t, parent2, u.GetRelationUnchecked(e, child2ID))

	u.SetRelations(e, RelID(childID, parent2), RelID(child2ID, parent1))
	assert.Equal(t, parent2, u.GetRelation(e, childID))
	assert.Equal(t, parent1, u.GetRelationUnchecked(e, child2ID))
}

func TestUnsafeAddRemove(t *testing.T) {
	w := NewWorld(16)
	u := w.Unsafe()

	posID := ComponentID[Position](&w)
	childID := ComponentID[ChildOf](&w)

	e1 := w.NewEntity()
	u.Add(e1, posID)

	assert.True(t, u.Has(e1, posID))

	e2 := w.NewEntity()
	u.AddRel(e2, []ID{posID, childID}, RelID(childID, e1))

	assert.True(t, u.Has(e2, posID))
	assert.True(t, u.Has(e2, childID))
	assert.Equal(t, e1, u.GetRelation(e2, childID))

	u.Remove(e1, posID)
	assert.False(t, u.Has(e1, posID))

	assert.Panics(t, func() {
		u.Add(Entity{}, posID)
	})
	assert.Panics(t, func() {
		u.AddRel(Entity{}, []ID{posID, childID}, RelID(childID, e1))
	})
	assert.Panics(t, func() {
		u.Remove(Entity{}, posID)
	})
}

func TestUnsafeExchange(t *testing.T) {
	w := NewWorld(16)
	u := w.Unsafe()

	posID := ComponentID[Position](&w)
	childID := ComponentID[ChildOf](&w)

	parent := u.NewEntity()
	e := u.NewEntity(posID)

	u.Exchange(e, []ID{childID}, []ID{posID}, RelID(childID, parent))
	assert.False(t, u.Has(e, posID))
	assert.True(t, u.Has(e, childID))

	child := (*ChildOf)(u.Get(e, childID))
	assert.NotNil(t, child)
	assert.Equal(t, parent, u.GetRelation(e, childID))

	assert.Panics(t, func() {
		u.Exchange(Entity{}, []ID{childID}, []ID{posID})
	})
}

func TestUnsafeIDs(t *testing.T) {
	w := NewWorld(16)
	u := w.Unsafe()

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)

	e := u.NewEntity(posID, velID)
	ids := u.IDs(e)
	assert.Equal(t, []ID{posID, velID}, ids.data)

	assert.Equal(t, 2, ids.Len())
	assert.Equal(t, posID, ids.Get(0))
	assert.Equal(t, velID, ids.Get(1))

	assert.Panics(t, func() {
		u.IDs(Entity{})
	})
}

func TestUnsafeEntityDump(t *testing.T) {
	w := NewWorld(1024)

	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()
	e4 := w.NewEntity()

	w.RemoveEntity(e2)
	w.RemoveEntity(e3)
	e5 := w.NewEntity()

	eData := w.Unsafe().DumpEntities()
	fmt.Println(eData)

	w2 := NewWorld(1024)
	w2.Unsafe().LoadEntities(&eData)

	assert.True(t, w2.Alive(e1))
	assert.True(t, w2.Alive(e4))
	assert.True(t, w2.Alive(e5))

	assert.False(t, w2.Alive(e2))
	assert.False(t, w2.Alive(e3))

	//assert.Equal(t, w.Ids(e1), []ID{})

	query := NewFilter(&w2).Query()
	assert.Equal(t, query.Count(), 3)
	query.Close()
}

func TestUnsafeEntityDumpEmpty(t *testing.T) {
	w := NewWorld(1024)

	eData := w.Unsafe().DumpEntities()

	w2 := NewWorld(1024)
	w2.Unsafe().LoadEntities(&eData)

	e1 := w2.NewEntity()
	e2 := w2.NewEntity()

	assert.True(t, w2.Alive(e1))
	assert.True(t, w2.Alive(e2))

	query := NewFilter(&w2).Query()
	assert.Equal(t, 2, query.Count())
	query.Close()
}

func TestUnsafeEntityDumpFail(t *testing.T) {
	w := NewWorld(1024)
	_ = w.NewEntity()

	eData := w.Unsafe().DumpEntities()

	w2 := NewWorld(1024)
	e1 := w2.NewEntity()
	w2.RemoveEntity(e1)

	assert.PanicsWithValue(t, "can set entity data only on a fresh or reset world",
		func() {
			w2.Unsafe().LoadEntities(&eData)
		})
}
