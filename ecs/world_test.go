package ecs

import (
	"math/rand/v2"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorld(t *testing.T) {
	w := NewWorld(1024)

	assert.Equal(t, 2, len(w.storage.entities))
	assert.Equal(t, 1, len(w.storage.tables))
	assert.Equal(t, 1, len(w.storage.archetypes))
	assert.Equal(t, 1, len(w.storage.archetypes[0].tables))
}

func TestWorldNewEntity(t *testing.T) {
	w := NewWorld(8)

	assert.False(t, w.Alive(Entity{}))
	for i := range 10 {
		e := w.NewEntity()
		assert.EqualValues(t, e.id, i+2)
		assert.EqualValues(t, e.gen, 0)
		assert.True(t, w.Alive(e))
	}
	assert.Equal(t, 12, len(w.storage.entities))

	idx := w.storage.entities[4]
	assert.EqualValues(t, 0, idx.table)
	assert.EqualValues(t, 2, idx.row)
}

func TestWorldExchange(t *testing.T) {
	w := NewWorld(2)

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)

	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()

	w.exchange(e1, []ID{posID}, nil, nil, nil)
	w.exchange(e2, []ID{posID, velID}, nil, nil, nil)
	w.exchange(e3, []ID{posID, velID}, nil, nil, nil)

	assert.True(t, w.storage.has(e1, posID))
	assert.False(t, w.storage.has(e1, velID))

	assert.True(t, w.storage.has(e2, posID))
	assert.True(t, w.storage.has(e2, velID))

	pos := (*Position)(w.storage.get(e1, posID))
	pos.X = 100

	pos = (*Position)(w.storage.get(e1, posID))
	assert.Equal(t, pos.X, 100.0)

	w.exchange(e2, nil, []ID{posID}, nil, nil)
	assert.False(t, w.storage.has(e2, posID))
	assert.True(t, w.storage.has(e2, velID))
}

func TestWorldRemoveEntity(t *testing.T) {
	w := NewWorld(32)

	mapper := NewMap2[Position, Velocity](&w)

	entities := make([]Entity, 0, 100)
	for range 100 {
		e := mapper.NewEntity(&Position{}, &Velocity{})
		assert.True(t, w.Alive(e))
		entities = append(entities, e)
	}

	filter := NewFilter0(&w)
	query := filter.Query()
	cnt := 0
	for query.Next() {
		assert.True(t, w.Alive(query.Entity()))
		cnt++
	}
	assert.Equal(t, 100, cnt)

	for _, e := range entities {
		w.RemoveEntity(e)
		assert.False(t, w.Alive(e))
	}

	query = filter.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)

	e := w.NewEntity()
	w.RemoveEntity(e)
	assert.False(t, w.Alive(e))
}

func TestWorldRelations(t *testing.T) {
	w := NewWorld(16)

	_ = ComponentID[CompA](&w)
	_ = ComponentID[CompB](&w)

	parent1 := w.NewEntity()
	parent2 := w.NewEntity()

	mapper1 := NewMap3[Position, ChildOf, ChildOf2](&w)
	assert.True(t, w.storage.registry.IsRelation[ComponentID[ChildOf](&w).id])
	assert.True(t, w.storage.registry.IsRelation[ComponentID[ChildOf2](&w).id])

	for range 10 {
		mapper1.NewEntity(&Position{}, &ChildOf{}, &ChildOf2{}, Rel(1, parent1), Rel(2, parent1))
		mapper1.NewEntity(&Position{}, &ChildOf{}, &ChildOf2{}, Rel(1, parent1), Rel(2, parent2))
		mapper1.NewEntity(&Position{}, &ChildOf{}, &ChildOf2{}, Rel(1, parent2), Rel(2, parent1))
		mapper1.NewEntity(&Position{}, &ChildOf{}, &ChildOf2{}, Rel(1, parent2), Rel(2, parent2))
	}

	filter := NewFilter3[Position, ChildOf, ChildOf2](&w)

	query := filter.Query()
	cnt := 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 40, cnt)

	query = filter.Query(Rel(1, parent1), Rel(2, parent2))
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 10, cnt)

	query = filter.Query(Rel(1, parent1))
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 20, cnt)

	mapper2 := NewMap2[Position, ChildOf](&w)
	child2Map := NewMap1[ChildOf2](&w)

	e := mapper2.NewEntity(&Position{}, &ChildOf{}, Rel(1, parent1))
	child2Map.Add(e, &ChildOf2{}, Rel(0, parent2))
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
	assert.Equal(t, parent1, map1.GetRelation(e))

	map1.SetRelation(e, parent2)
	assert.Equal(t, parent2, map1.GetRelation(e))
	assert.Equal(t, parent1, map2.GetRelation(e))
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
		e := posChildMap.NewEntity(&Position{X: -1, Y: 1}, &ChildOf{}, Rel(1, parent1))
		assert.Equal(t, parent1, childMap.GetRelation(e))
		entities = append(entities, e)
	}
	_ = posChildMap.NewEntity(&Position{}, &ChildOf{}, Rel(1, parent2))

	w.RemoveEntity(parent1)

	for _, e := range entities {
		assert.Equal(t, Entity{}, childMap.GetRelation(e))
	}

	archetype := &w.storage.archetypes[1]
	assert.Equal(t, []tableID{3, 2}, archetype.tables)
	assert.Equal(t, []tableID{1}, archetype.freeTables)

	for _, e := range entities {
		childMap.SetRelation(e, parent3)
		assert.Equal(t, parent3, childMap.GetRelation(e))
	}
	assert.Equal(t, []tableID{3, 2, 1}, archetype.tables)
	assert.Equal(t, []tableID{}, archetype.freeTables)

	filter := NewFilter2[Position, ChildOf](&w)
	query := filter.Query(Rel(1, parent3))
	cnt := 0
	for query.Next() {
		pos, _ := query.Get()
		assert.Equal(t, Position{X: -1, Y: 1}, *pos)
		cnt++
	}
	assert.Equal(t, 32, cnt)
}

func TestWorldReset(t *testing.T) {
	world := NewWorld(16)
	u := world.Unsafe()

	AddResource(&world, &Heading{100})

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)
	relID := ComponentID[ChildOf](&world)

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

	assert.Equal(t, 0, int(world.storage.tables[0].Len()))
	assert.Equal(t, 0, int(world.storage.tables[1].Len()))
	assert.Equal(t, 0, world.storage.entityPool.Len())
	assert.Equal(t, 2, len(world.storage.entities))

	query := u.Query(NewFilter())
	assert.Equal(t, 0, query.Count())
	query.Close()

	e1 = u.NewEntity(posID)
	e2 := u.NewEntity(velID)
	u.NewEntity(posID, velID)
	u.NewEntity(posID, velID)

	assert.Equal(t, Entity{2, 0}, e1)
	assert.Equal(t, Entity{3, 0}, e2)

	query = u.Query(NewFilter())
	assert.Equal(t, 4, query.Count())
	query.Close()
}

func TestWorldRemoveGC(t *testing.T) {
	w := NewWorld(128)
	mapper := NewMap[SliceComp](&w)

	runtime.GC()
	mem1 := runtime.MemStats{}
	mem2 := runtime.MemStats{}
	runtime.ReadMemStats(&mem1)

	entities := []Entity{}
	for i := 0; i < 100; i++ {
		e := mapper.NewEntity(&SliceComp{})
		ws := mapper.Get(e)
		ws.Slice = make([]int, 10000)
		entities = append(entities, e)
	}

	runtime.GC()
	runtime.ReadMemStats(&mem2)
	heap := int(mem2.HeapInuse - mem1.HeapInuse)
	assert.Greater(t, heap, 8000000)
	assert.Less(t, heap, 10000000)

	rand.Shuffle(len(entities), func(i, j int) {
		entities[i], entities[j] = entities[j], entities[i]
	})

	for _, e := range entities {
		w.RemoveEntity(e)
	}

	runtime.GC()
	runtime.ReadMemStats(&mem2)
	heap = int(mem2.HeapInuse - mem1.HeapInuse)
	assert.Less(t, heap, 800000)

	_ = mapper.NewEntity(&SliceComp{})
}

func TestWorldPointerStressTest(t *testing.T) {
	w := NewWorld(128)

	mapper := NewMap[PointerComp](&w)

	count := 0
	var entities []Entity

	for i := 0; i < 1000; i++ {
		add := rand.IntN(1000)
		count += add
		for n := 0; n < add; n++ {
			e := mapper.NewEntity(&PointerComp{})
			ptr := mapper.Get(e)
			ptr.Ptr = &PointerType{&Position{X: float64(e.id), Y: 2}}
		}

		filter := NewFilter1[PointerComp](&w)
		query := filter.Query()
		for query.Next() {
			ptr := query.Get()
			assert.EqualValues(t, ptr.Ptr.Pos.X, int(query.Entity().id))
			entities = append(entities, query.Entity())
		}
		rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

		rem := rand.IntN(count)
		count -= rem
		for n := 0; n < rem; n++ {
			w.RemoveEntity(entities[n])
		}

		entities = entities[:0]
		runtime.GC()
	}
}
