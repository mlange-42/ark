package ecs

import (
	"fmt"
	"math/rand/v2"
	"runtime"
	"testing"
)

func TestNewWorld(t *testing.T) {
	w := NewWorld(1024)

	expectEqual(t, len(w.storage.entities), 2)
	expectEqual(t, len(w.storage.tables), 1)
	expectEqual(t, len(w.storage.archetypes), 1)
	expectEqual(t, len(w.storage.archetypes[0].tables.tables), 1)
}

func TestWorldNewEntity(t *testing.T) {
	w := NewWorld(8)

	expectFalse(t, w.Alive(Entity{}))
	for i := range 10 {
		e := w.NewEntity()
		expectEqual(t, e.id, entityID(i+2))
		expectEqual(t, e.gen, 0)
		expectTrue(t, w.Alive(e))
	}
	expectEqual(t, len(w.storage.entities), 12)

	idx := w.storage.entities[4]
	expectEqual(t, idx.table, 0)
	expectEqual(t, idx.row, 2)
}

func TestWorldExchange(t *testing.T) {
	w := NewWorld(2)

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)

	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()

	w.exchange(e1, []ID{posID}, nil, nil)
	w.exchange(e2, []ID{posID, velID}, nil, nil)
	w.exchange(e3, []ID{posID, velID}, nil, nil)

	expectTrue(t, w.storage.has(e1, posID))
	expectFalse(t, w.storage.has(e1, velID))

	expectTrue(t, w.storage.has(e2, posID))
	expectTrue(t, w.storage.has(e2, velID))

	pos := (*Position)(w.storage.get(e1, posID))
	pos.X = 100

	pos = (*Position)(w.storage.get(e1, posID))
	expectEqual(t, pos.X, 100.0)

	w.exchange(e2, nil, []ID{posID}, nil)
	expectFalse(t, w.storage.has(e2, posID))
	expectTrue(t, w.storage.has(e2, velID))
}

func TestWorldRemoveEntity(t *testing.T) {
	w := NewWorld(32)
	mapper := NewMap2[Position, Velocity](&w)

	entities := make([]Entity, 0, 100)
	for range 100 {
		e := mapper.NewEntity(&Position{}, &Velocity{})
		expectTrue(t, w.Alive(e))
		entities = append(entities, e)
	}

	filter := NewFilter0(&w)
	query := filter.Query()
	cnt := 0
	for query.Next() {
		expectTrue(t, w.Alive(query.Entity()))
		cnt++
	}
	expectEqual(t, cnt, 100)

	for _, e := range entities {
		w.RemoveEntity(e)
		expectFalse(t, w.Alive(e))
	}

	query = filter.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	expectEqual(t, cnt, 0)

	e := w.NewEntity()
	w.RemoveEntity(e)
	expectFalse(t, w.Alive(e))

	expectPanicWithValue(t, "can't remove a dead entity", func() {
		w.RemoveEntity(e)
	})
}
func TestWorldNewEntities(t *testing.T) {
	n := 100
	w := NewWorld(16)

	cnt := 0
	w.NewEntities(n, func(entity Entity) {
		expectEqual(t, entity.ID(), uint32(cnt+2))
		cnt++
	})
	expectEqual(t, cnt, n)

	w.NewEntities(n, nil)

	filter := NewUnsafeFilter(&w)
	query := filter.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	expectEqual(t, cnt, 2*n)
}

func TestWorldRemoveEntities(t *testing.T) {
	n := 12
	w := NewWorld(16)

	posMap := NewMap1[Position](&w)
	velMap := NewMap1[Velocity](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 0
	posMap.NewBatchFn(n, func(entity Entity, _ *Position) { cnt++ })
	velMap.NewBatchFn(n, func(entity Entity, _ *Velocity) { cnt++ })
	posVelMap.NewBatchFn(n, func(entity Entity, _ *Position, _ *Velocity) { cnt++ })
	expectEqual(t, cnt, n*3)

	filter := NewFilter2[Position, Velocity](&w)
	cnt = 0
	w.RemoveEntities(filter.Batch(), func(entity Entity) { cnt++ })
	expectEqual(t, cnt, n)

	filter2 := NewFilter1[Position](&w).Register()
	cnt = 0
	w.RemoveEntities(filter2.Batch(), func(entity Entity) { cnt++ })
	expectEqual(t, cnt, n)

	filter3 := NewFilter0(&w)
	query := filter3.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	expectEqual(t, cnt, n)
}

func TestWorldRelations(t *testing.T) {
	w := NewWorld(16)
	_ = ComponentID[CompA](&w)
	_ = ComponentID[CompB](&w)

	parent1 := w.NewEntity()
	parent2 := w.NewEntity()
	parent3 := w.NewEntity()

	mapper1 := NewMap3[Position, ChildOf, ChildOf2](&w)
	expectTrue(t, w.storage.registry.IsRelation[ComponentID[ChildOf](&w).id])
	expectTrue(t, w.storage.registry.IsRelation[ComponentID[ChildOf2](&w).id])

	for i := 0; i < 10; i++ {
		mapper1.NewEntity(&Position{}, &ChildOf{}, &ChildOf2{}, Rel[ChildOf](parent1), RelIdx(2, parent1))
		mapper1.NewEntity(&Position{}, &ChildOf{}, &ChildOf2{}, RelIdx(1, parent1), RelIdx(2, parent2))
		mapper1.NewEntity(&Position{}, &ChildOf{}, &ChildOf2{}, RelIdx(1, parent2), RelIdx(2, parent1))
		mapper1.NewEntity(&Position{}, &ChildOf{}, &ChildOf2{}, Rel[ChildOf](parent2), Rel[ChildOf2](parent2))
		mapper1.NewEntity(&Position{}, &ChildOf{}, &ChildOf2{}, Rel[ChildOf](parent1), Rel[ChildOf2](parent3))
	}

	filter := NewFilter3[Position, ChildOf, ChildOf2](&w)

	query := filter.Query()
	expectEqual(t, query.Count(), 50)
	cnt := 0
	for query.Next() {
		cnt++
	}
	expectEqual(t, cnt, 50)

	query = filter.Query(RelIdx(1, parent1), RelIdx(2, parent2))
	expectEqual(t, query.Count(), 10)
	cnt = 0
	for query.Next() {
		cnt++
	}
	expectEqual(t, cnt, 10)

	query = filter.Query(RelIdx(1, parent1))
	expectEqual(t, query.Count(), 30)
	cnt = 0
	for query.Next() {
		cnt++
	}
	expectEqual(t, cnt, 30)

	mapper2 := NewMap2[Position, ChildOf](&w)
	child2Map := NewMap1[ChildOf2](&w)

	e := mapper2.NewEntity(&Position{}, &ChildOf{}, RelIdx(1, parent1))
	child2Map.Add(e, &ChildOf2{}, RelIdx(0, parent2))

	child2Map.SetRelations(e, RelIdx(0, parent1))
	expectEqual(t, child2Map.GetRelation(e, 0), parent1)

	child2Map.SetRelations(e, RelIdx(0, parent1))
	expectEqual(t, child2Map.GetRelation(e, 0), parent1)

	child2Map.Remove(e)

	expectPanicWithValue(t, "entity has no component of type ChildOf2 to set relation target for", func() {
		child2Map.SetRelations(e, RelIdx(0, parent2))
	})
}

func TestWorldSetRelations(t *testing.T) {
	w := NewWorld(16)
	_ = ComponentID[CompA](&w)
	_ = ComponentID[CompB](&w)

	parent1 := w.NewEntity()
	parent2 := w.NewEntity()

	map1 := NewMap[ChildOf](&w)
	map2 := NewMap[ChildOf2](&w)

	e := map1.NewEntity(&ChildOf{}, parent1)
	map2.Add(e, &ChildOf2{}, parent1)
	expectEqual(t, map1.GetRelation(e), parent1)

	map1.SetRelation(e, parent2)
	expectEqual(t, map1.GetRelation(e), parent2)
	expectEqual(t, map2.GetRelation(e), parent1)
}

func TestWorldRelationRemoveTarget(t *testing.T) {
	w := NewWorld(16)
	_ = ComponentID[CompA](&w)
	_ = ComponentID[CompB](&w)

	parent1 := w.NewEntity()
	parent2 := w.NewEntity()
	parent3 := w.NewEntity()

	childMap := NewMap[ChildOf](&w)
	posChildMap := NewMap2[Position, ChildOf](&w)

	entities := []Entity{}
	for range 32 {
		e := posChildMap.NewEntity(&Position{X: -1, Y: 1}, &ChildOf{}, RelIdx(1, parent1))
		expectEqual(t, childMap.GetRelation(e), parent1)
		entities = append(entities, e)
	}
	_ = posChildMap.NewEntity(&Position{}, &ChildOf{}, RelIdx(1, parent2))

	w.RemoveEntity(parent1)

	for _, e := range entities {
		expectEqual(t, childMap.GetRelation(e), Entity{})
	}

	archetype := &w.storage.archetypes[1]
	if !equalSlices(archetype.tables.tables, []tableID{3, 2}) {
		t.Errorf("expected %v, got %v", []tableID{3, 2}, archetype.tables.tables)
	}
	if !equalSlices(archetype.freeTables, []tableID{1}) {
		t.Errorf("expected %v, got %v", []tableID{1}, archetype.freeTables)
	}

	for _, e := range entities {
		childMap.SetRelation(e, parent3)
		expectEqual(t, childMap.GetRelation(e), parent3)
	}
	if !equalSlices(archetype.tables.tables, []tableID{3, 2, 1}) {
		t.Errorf("expected %v, got %v", []tableID{3, 2, 1}, archetype.tables.tables)
	}
	if !equalSlices(archetype.freeTables, []tableID{}) {
		t.Errorf("expected %v, got %v", []tableID{}, archetype.freeTables)
	}

	filter := NewFilter2[Position, ChildOf](&w)
	query := filter.Query(RelIdx(1, parent3))
	cnt := 0
	for query.Next() {
		pos, _ := query.Get()
		expectEqual(t, *pos, Position{X: -1, Y: 1})
		cnt++
	}
	expectEqual(t, cnt, 32)
}

func TestWorldReset(t *testing.T) {
	world := NewWorld(16)
	u := world.Unsafe()

	AddResource(&world, &Heading{100})

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)
	relID := ComponentID[ChildOf](&world)

	_ = NewFilter2[Position, Velocity](&world).Register()

	target1 := world.NewEntity()
	target2 := world.NewEntity()

	u.NewEntity(velID)
	u.NewEntity(posID, velID)
	u.NewEntity(posID, velID)
	e1 := u.NewEntityRel([]ID{posID, relID}, RelID(relID, target1))
	_ = u.NewEntityRel([]ID{posID, relID}, RelID(relID, target2))

	world.RemoveEntity(e1)
	world.RemoveEntity(target1)

	world.Reset()

	expectEqual(t, int(world.storage.tables[0].Len()), 0)
	expectEqual(t, int(world.storage.tables[1].Len()), 0)
	expectEqual(t, world.storage.entityPool.Len(), 0)
	expectEqual(t, len(world.storage.entities), 2)
	expectEqual(t, len(world.storage.isTarget), 2)

	query := NewUnsafeFilter(&world).Query()
	expectEqual(t, query.Count(), 0)
	query.Close()

	e1 = u.NewEntity(posID)
	e2 := u.NewEntity(velID)
	u.NewEntity(posID, velID)
	u.NewEntity(posID, velID)

	expectEqual(t, e1, Entity{2, 0})
	expectEqual(t, e2, Entity{3, 0})

	query = NewUnsafeFilter(&world).Query()
	expectEqual(t, query.Count(), 4)
	query.Close()
}
func TestWorldLock(t *testing.T) {
	w := NewWorld()

	l1 := w.lock()
	expectTrue(t, w.IsLocked())
	expectPanicWithValue(t, "attempt to modify a locked world", func() { w.checkLocked() })

	l2 := w.lock()
	expectTrue(t, w.IsLocked())
	w.unlock(l1)
	expectTrue(t, w.IsLocked())
	w.unlock(l2)
	expectFalse(t, w.IsLocked())
}

func TestWorldRemoveGC(t *testing.T) {
	w := NewWorld(128)
	mapper := NewMap[SliceComp](&w)

	runtime.GC()
	mem1 := runtime.MemStats{}
	mem2 := runtime.MemStats{}
	runtime.ReadMemStats(&mem1)

	entities := []Entity{}
	for range 100 {
		e := mapper.NewEntity(&SliceComp{})
		ws := mapper.Get(e)
		ws.Slice = make([]int, 10000)
		entities = append(entities, e)
	}

	runtime.GC()
	runtime.ReadMemStats(&mem2)
	heap := int(mem2.HeapInuse - mem1.HeapInuse)
	expectTrue(t, heap > 8000000)
	expectTrue(t, heap < 10000000)

	rand.Shuffle(len(entities), func(i, j int) {
		entities[i], entities[j] = entities[j], entities[i]
	})

	for _, e := range entities {
		w.RemoveEntity(e)
	}

	runtime.GC()
	runtime.ReadMemStats(&mem2)
	heap = int(mem2.HeapInuse - mem1.HeapInuse)
	expectTrue(t, heap < 800000)

	_ = mapper.NewEntity(&SliceComp{})
}

func TestWorldPointerStressTest(t *testing.T) {
	w := NewWorld(128)
	mapper := NewMap[PointerComp](&w)

	count := 1
	var entities []Entity

	for range 1000 {
		add := rand.IntN(1000)
		count += add
		for j := 0; j < add; j++ {
			e := mapper.NewEntity(&PointerComp{})
			ptr := mapper.Get(e)
			ptr.Ptr = &PointerType{&Position{X: float64(e.id), Y: 2}}
		}

		filter := NewFilter1[PointerComp](&w)
		query := filter.Query()
		for query.Next() {
			ptr := query.Get()
			expectEqual(t, ptr.Ptr.Pos.X, float64(query.Entity().id))
			entities = append(entities, query.Entity())
		}
		rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

		rem := rand.IntN(count)
		count -= rem
		for _, e := range entities {
			w.RemoveEntity(e)
		}

		entities = entities[:0]
		runtime.GC()
	}
}

func TestWorldPanics(t *testing.T) {
	w := NewWorld(128, 32)
	_ = ComponentID[Position](&w)
	_ = ComponentID[Velocity](&w)
	childID := ComponentID[ChildOf](&w)

	e := w.NewEntity()
	w.exchange(e, nil, nil, nil)
	w.RemoveEntity(e)

	expectPanicWithValue(t, "exchange operation has no effect, but relations were specified. Use SetRelation(s) instead", func() {
		e := w.NewEntity()
		w.exchange(e, nil, nil, []RelationID{RelID(childID, e)})
		w.RemoveEntity(e)
	})

	e = w.NewEntity()
	w.exchangeBatch(nil, nil, nil, nil, nil)
	w.RemoveEntity(e)

	expectPanicWithValue(t, "exchange operation has no effect, but relations were specified. Use SetRelationBatch instead", func() {
		e := w.NewEntity()
		w.exchangeBatch(nil, nil, nil, []RelationID{RelID(childID, e)}, nil)
		w.RemoveEntity(e)
	})
}

func TestWorldStats(t *testing.T) {
	w := NewWorld(128, 32)

	posVelMap := NewMap2[Position, Velocity](&w)
	posVelHeadMap := NewMap3[Position, Velocity, Heading](&w)
	posChildMap := NewMap3[Position, ChildOf, ChildOf2](&w)
	filter := NewFilter0(&w)

	p1 := w.NewEntity()
	p2 := w.NewEntity()
	p3 := w.NewEntity()

	posVelMap.NewBatchFn(100, nil)
	posChildMap.NewBatchFn(50, nil, RelIdx(1, p1), RelIdx(2, p2))
	posChildMap.NewBatchFn(50, nil, RelIdx(1, p3), RelIdx(2, p2))

	w.RemoveEntities(filter.Batch(), nil)
	_ = w.Stats()

	p1 = w.NewEntity()
	p2 = w.NewEntity()
	p3 = w.NewEntity()

	posVelMap.NewBatchFn(100, nil)
	posChildMap.NewBatchFn(50, nil, RelIdx(1, p1), RelIdx(2, p2))
	posChildMap.NewBatchFn(50, nil, RelIdx(1, p3), RelIdx(2, p2))

	posVelHeadMap.NewBatchFn(250, nil)
	posChildMap.NewBatchFn(50, nil, RelIdx(1, p2), RelIdx(2, p3))
	_ = w.Stats()

	stats := w.Stats()
	fmt.Println(stats.String())

	w.RemoveEntities(filter.Batch(), nil)
	stats = w.Stats()
	fmt.Println(stats.String())
}

func TestWorldCreateManyTables(t *testing.T) {
	n := 1000
	w := NewWorld()
	dataMap := NewMap1[Position](&w)

	id := ComponentID[Position](&w)
	expectTrue(t, w.storage.registry.IsTrivial[id.id])

	entities := make([]Entity, 0)
	for i := range n {
		entities = append(entities, dataMap.NewEntity(&Position{X: float64(i)}))
	}

	filter := NewFilter1[Position](&w)
	q := filter.Query()
	expectEqual(t, q.Count(), n)
	q.Close()

	relMap := NewMap1[ChildOf](&w)
	for i := range n {
		relMap.Add(entities[i], &ChildOf{}, Rel[ChildOf](entities[(i+1)%n]))
	}

	q = filter.Query()
	expectEqual(t, q.Count(), n)
	q.Close()
}

func TestWorldCreateManyTablesSlice(t *testing.T) {
	n := 1000
	w := NewWorld()
	dataMap := NewMap1[SliceComp](&w)

	id := ComponentID[SliceComp](&w)
	expectFalse(t, w.storage.registry.IsTrivial[id.id])

	entities := make([]Entity, 0)
	for range n {
		e := dataMap.NewEntity(&SliceComp{Slice: nil})
		sl := dataMap.Get(e)
		sl.Slice = []int{int(e.id) + 1, int(e.id) + 2, int(e.id) + 3, int(e.id) + 4}
		entities = append(entities, e)
	}

	filter := NewFilter1[SliceComp](&w)
	q := filter.Query()
	expectEqual(t, q.Count(), n)
	q.Close()

	velMap := NewMap1[Velocity](&w)
	for i := range n {
		velMap.Add(entities[i], &Velocity{})
	}

	relMap := NewMap1[ChildOf](&w)
	for i := range n {
		relMap.Add(entities[i], &ChildOf{}, Rel[ChildOf](entities[(i+1)%n]))
	}

	q = filter.Query()
	expectEqual(t, q.Count(), n)
	for q.Next() {
		sl := q.Get()
		e := q.Entity()

		want := []int{int(e.id) + 1, int(e.id) + 2, int(e.id) + 3, int(e.id) + 4}
		if !equalSlices(sl.Slice, want) {
			t.Errorf("expected %v, got %v", want, sl.Slice)
		}
	}
}
