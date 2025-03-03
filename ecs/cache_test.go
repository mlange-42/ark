package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterCache(t *testing.T) {
	world := NewWorld()

	mapper2 := NewMap2[Position, Velocity](&world)
	mapper3 := NewMap3[Position, Velocity, Heading](&world)

	world.NewEntity()
	mapper2.NewEntity(&Position{}, &Velocity{})
	mapper3.NewEntity(&Position{}, &Velocity{}, &Heading{})

	filter2 := NewFilter2[Position, Velocity](&world).Register()
	filter3 := NewFilter3[Position, Velocity, Heading](&world).Register()

	assert.Equal(t, 0, int(filter2.cache.id))
	assert.Equal(t, 1, int(filter3.cache.id))

	assert.Equal(t, 2, len(world.storage.getTableIDs(filter2.Batch())))
	assert.Equal(t, 1, len(world.storage.getTableIDs(filter3.Batch())))

	assert.PanicsWithValue(t, "filter is already registered, can't register", func() { filter2.Register() })

	filter2.Unregister()
	filter3.Unregister()

	assert.PanicsWithValue(t, "no filter for id found to unregister", func() { filter2.Unregister() })
}

/*
func TestFilterCacheRelation(t *testing.T) {
	world := NewWorld()
	posID := ComponentID[Position](&world)
	rel1ID := ComponentID[testRelationA](&world)
	rel2ID := ComponentID[testRelationB](&world)

	target1 := world.NewEntity()
	target2 := world.NewEntity()
	target3 := world.NewEntity()
	target4 := world.NewEntity()

	cache := world.Cache()

	f1 := All(rel1ID)
	ff1 := cache.Register(f1)

	f2 := NewRelationFilter(f1, target1)
	ff2 := cache.Register(&f2)

	f3 := NewRelationFilter(f1, target2)
	ff3 := cache.Register(&f3)

	c1 := world.Cache().get(&ff1)
	c2 := world.Cache().get(&ff2)
	c3 := world.Cache().get(&ff3)

	NewBuilder(&world, posID).NewBatch(10)

	assert.Equal(t, int32(0), c1.Archetypes.Len())
	assert.Equal(t, int32(0), c2.Archetypes.Len())
	assert.Equal(t, int32(0), c3.Archetypes.Len())

	e1 := NewBuilder(&world, rel1ID).WithRelation(rel1ID).New(target1)
	assert.Equal(t, int32(1), c1.Archetypes.Len())
	assert.Equal(t, int32(1), c2.Archetypes.Len())

	_ = NewBuilder(&world, rel1ID).WithRelation(rel1ID).New(target3)
	assert.Equal(t, int32(2), c1.Archetypes.Len())
	assert.Equal(t, int32(1), c2.Archetypes.Len())

	_ = NewBuilder(&world, rel2ID).WithRelation(rel2ID).New(target2)

	world.RemoveEntity(e1)
	world.RemoveEntity(target1)
	assert.Equal(t, int32(1), c1.Archetypes.Len())
	assert.Equal(t, int32(0), c2.Archetypes.Len())

	_ = NewBuilder(&world, rel1ID).WithRelation(rel1ID).New(target2)
	_ = NewBuilder(&world, rel1ID, posID).WithRelation(rel1ID).New(target2)
	_ = NewBuilder(&world, rel1ID, posID).WithRelation(rel1ID).New(target3)
	_ = NewBuilder(&world, rel1ID, posID).WithRelation(rel1ID).New(target4)
	assert.Equal(t, int32(5), c1.Archetypes.Len())
	assert.Equal(t, int32(2), c3.Archetypes.Len())

	world.Batch().RemoveEntities(All())
	assert.Equal(t, int32(0), c1.Archetypes.Len())
	assert.Equal(t, int32(0), c2.Archetypes.Len())
}
*/
