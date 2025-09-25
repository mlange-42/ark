package ecs

import (
	"fmt"
	"testing"
)

func TestUnsafe(t *testing.T) {
	w := NewWorld(1024)
	u := w.Unsafe()

	expectEqual(t, &w, u.world)
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

	expectPanicsWithValue(t, "can't get component of a dead entity",
		func() {
			u.Get(Entity{}, posID)
		})
	expectPanicsWithValue(t, "can't check component of a dead entity",
		func() {
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

	expectEqual(t, parent1, u.GetRelation(e, childID))
	expectEqual(t, parent2, u.GetRelationUnchecked(e, child2ID))

	u.SetRelations(e, RelID(childID, parent2), RelID(child2ID, parent1))
	expectEqual(t, parent2, u.GetRelation(e, childID))
	expectEqual(t, parent1, u.GetRelationUnchecked(e, child2ID))
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
	expectEqual(t, e1, u.GetRelation(e2, childID))

	u.Remove(e1, posID)
	expectFalse(t, u.Has(e1, posID))

	expectPanicsWithValue(t, "at least one component required to add",
		func() {
			u.Add(e1)
		})
	expectPanicsWithValue(t, "at least one component required to add",
		func() {
			u.AddRel(e1, []ID{}, RelID(childID, e1))
		})
	expectPanicsWithValue(t, "at least one component required to remove",
		func() {
			u.Remove(e1)
		})

	expectPanicsWithValue(t, "can't add components to a dead entity",
		func() {
			u.Add(Entity{}, posID)
		})
	expectPanicsWithValue(t, "can't add components to a dead entity",
		func() {
			u.AddRel(Entity{}, []ID{posID, childID}, RelID(childID, e1))
		})
	expectPanicsWithValue(t, "can't remove components from a dead entity",
		func() {
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
	expectNotNil(t, child)
	expectEqual(t, parent, u.GetRelation(e, childID))

	expectPanicsWithValue(t, "can't exchange components on a dead entity",
		func() {
			u.Exchange(Entity{}, []ID{childID}, []ID{posID})
		})

	expectPanicsWithValue(t, "at least one component required to add or remove",
		func() {
			u.Exchange(e, []ID{}, []ID{})
		})
}

func TestUnsafeIDs(t *testing.T) {
	w := NewWorld(16)
	u := w.Unsafe()

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)

	e := u.NewEntity(posID, velID)
	ids := u.IDs(e)
	expectSlicesEqual(t, []ID{posID, velID}, ids.data)

	expectEqual(t, 2, ids.Len())
	expectEqual(t, posID, ids.Get(0))
	expectEqual(t, velID, ids.Get(1))
	expectEqual(t, posID.Index(), ids.Get(0).Index())
	expectEqual(t, velID.Index(), ids.Get(1).Index())

	expectPanicsWithValue(t, "can't get component IDs of a dead entity",
		func() {
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

	//expectEqual(t, w.Ids(e1), []ID{})

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
	expectEqual(t, 2, query.Count())
	query.Close()
}

func TestUnsafeEntityDumpFail(t *testing.T) {
	w := NewWorld(1024)
	_ = w.NewEntity()

	eData := w.Unsafe().DumpEntities()

	w2 := NewWorld(1024)
	e1 := w2.NewEntity()
	w2.RemoveEntity(e1)

	expectPanicsWithValue(t, "can set entity data only on a fresh or reset world",
		func() {
			w2.Unsafe().LoadEntities(&eData)
		})
}
