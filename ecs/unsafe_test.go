package ecs

import (
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
