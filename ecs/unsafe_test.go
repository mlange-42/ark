package ecs

import (
	"fmt"
	"testing"
)

func TestUnsafe(t *testing.T) {
	w := NewWorld(1024)
	u := w.Unsafe()

	expectEqual(t, u.world, &w)
}

func TestUnsafeNewEntity(t *testing.T) {
	w := NewWorld(16)
	u := w.Unsafe()

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)

	e := u.NewEntity(posID, velID)

	expectTrue(t, u.Has(e, posID))
	expectTrue(t, u.Has(e, velID))
}

func TestUnsafeGet(t *testing.T) {
	w := NewWorld(16)
	u := w.Unsafe()

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)

	e := u.NewEntity(posID)

	expectTrue(t, u.Has(e, posID))
	expectFalse(t, u.Has(e, velID))

	expectTrue(t, u.HasUnchecked(e, posID))
	expectFalse(t, u.HasUnchecked(e, velID))

	pos := (*Position)(u.Get(e, posID))
	pos.X = 100

	pos2 := (*Position)(u.GetUnchecked(e, posID))
	expectEqual(t, pos, pos2)

	expectPanic(t, func() {
		u.Get(Entity{}, posID)
	})
	expectPanic(t, func() {
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

	expectEqual(t, u.GetRelation(e, childID), parent1)
	expectEqual(t, u.GetRelationUnchecked(e, child2ID), parent2)

	u.SetRelations(e, RelID(childID, parent2), RelID(child2ID, parent1))
	expectEqual(t, u.GetRelation(e, childID), parent2)
	expectEqual(t, u.GetRelationUnchecked(e, child2ID), parent1)
}

func TestUnsafeAddRemove(t *testing.T) {
	w := NewWorld(16)
	u := w.Unsafe()

	posID := ComponentID[Position](&w)
	childID := ComponentID[ChildOf](&w)

	e1 := w.NewEntity()
	u.Add(e1, posID)
	expectTrue(t, u.Has(e1, posID))

	e2 := w.NewEntity()
	u.AddRel(e2, []ID{posID, childID}, RelID(childID, e1))

	expectTrue(t, u.Has(e2, posID))
	expectTrue(t, u.Has(e2, childID))
	expectEqual(t, u.GetRelation(e2, childID), e1)

	u.Remove(e1, posID)
	expectFalse(t, u.Has(e1, posID))

	expectPanic(t, func() {
		u.Add(Entity{}, posID)
	})
	expectPanic(t, func() {
		u.AddRel(Entity{}, []ID{posID, childID}, RelID(childID, e1))
	})
	expectPanic(t, func() {
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
	expectFalse(t, u.Has(e, posID))
	expectTrue(t, u.Has(e, childID))

	child := (*ChildOf)(u.Get(e, childID))
	if child == nil {
		t.Errorf("expected non-nil child component")
	}
	expectEqual(t, u.GetRelation(e, childID), parent)

	expectPanic(t, func() {
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

	expectEqual(t, ids.data, []ID{posID, velID})
	expectEqual(t, ids.Len(), 2)
	expectEqual(t, ids.Get(0), posID)
	expectEqual(t, ids.Get(1), velID)
	expectEqual(t, ids.Get(0).Index(), posID.Index())
	expectEqual(t, ids.Get(1).Index(), velID.Index())

	expectPanic(t, func() {
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

	expectTrue(t, w2.Alive(e1))
	expectTrue(t, w2.Alive(e4))
	expectTrue(t, w2.Alive(e5))

	expectFalse(t, w2.Alive(e2))
	expectFalse(t, w2.Alive(e3))

	query := NewUnsafeFilter(&w2).Query()
	expectEqual(t, query.Count(), 3)
	query.Close()
}

func TestUnsafeEntityDumpEmpty(t *testing.T) {
	w := NewWorld(1024)

	eData := w.Unsafe().DumpEntities()

	w2 := NewWorld(1024)
	w2.Unsafe().LoadEntities(&eData)

	e1 := w2.NewEntity()
	e2 := w2.NewEntity()

	expectTrue(t, w2.Alive(e1))
	expectTrue(t, w2.Alive(e2))

	query := NewUnsafeFilter(&w2).Query()
	expectEqual(t, query.Count(), 2)
	query.Close()
}

func TestUnsafeEntityDumpFail(t *testing.T) {
	w := NewWorld(1024)
	_ = w.NewEntity()

	eData := w.Unsafe().DumpEntities()

	w2 := NewWorld(1024)
	e1 := w2.NewEntity()
	w2.RemoveEntity(e1)

	expectPanicWithValue(t, "can set entity data only on a fresh or reset world", func() {
		w2.Unsafe().LoadEntities(&eData)
	})
}
