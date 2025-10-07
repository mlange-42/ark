package cheatsheet

import (
	"math/rand/v2"
	"testing"

	"github.com/mlange-42/ark/ecs"
)

type Position struct {
	X float64
	Y float64
}

type Velocity struct {
	X float64
	Y float64
}

type Altitude struct {
	Z float64
}

type Grid struct{}

func NewGrid(sx, sy int) Grid {
	return Grid{}
}

var world = ecs.NewWorld()
var mapper = ecs.NewMap2[Position, Velocity](&world)
var filter = ecs.NewFilter2[Position, Velocity](&world).Exclusive()
var entity = world.NewEntity()

func TestCreateWorld(t *testing.T) {
	world := ecs.NewWorld()
	_ = &world
}

func TestCreateWorldConfig(t *testing.T) {
	world := ecs.NewWorld(1024)
	_ = &world
}

func TestCreateEmpty(t *testing.T) {
	e := world.NewEntity()
	_ = e
}

func TestCreateMapper(t *testing.T) {
	mapper := ecs.NewMap2[Position, Velocity](&world)
	_ = mapper
}

func TestCreateEntity(t *testing.T) {
	e := mapper.NewEntity(
		&Position{X: 100, Y: 100},
		&Velocity{X: 1, Y: -1},
	)
	_ = e
}

func TestCreateEntityFn(t *testing.T) {
	e := mapper.NewEntityFn(func(pos *Position, vel *Velocity) {
		pos.X, pos.Y = 100, 100
		vel.X, vel.Y = 1, -1
	})
	_ = e
}

func TestCreateBatch(t *testing.T) {
	mapper.NewBatch(
		100,
		&Position{X: 100, Y: 100},
		&Velocity{X: 1, Y: -1},
	)
}

func TestCreateBatchFn(t *testing.T) {
	mapper.NewBatchFn(
		100,
		func(e ecs.Entity, pos *Position, vel *Velocity) {
			pos.X, pos.Y = rand.Float64()*100, rand.Float64()*100
			vel.X, vel.Y = rand.NormFloat64(), rand.NormFloat64()
		})
}

func TestRemoveEntity(t *testing.T) {
	world.RemoveEntity(entity)
}

func TestRemoveEntities(t *testing.T) {
	filter := ecs.NewFilter2[Position, Velocity](&world).Exclusive()

	world.RemoveEntities(filter.Batch(), nil)
}

func TestRemoveEntitiesFn(t *testing.T) {
	world.RemoveEntities(filter.Batch(), func(entity ecs.Entity) {
		// Do something before removal
	})
}

func TestAddRemoveComponents(t *testing.T) {
	entity := world.NewEntity()

	mapper.Add(
		entity,
		&Position{X: 100, Y: 100},
		&Velocity{X: 1, Y: -1},
	)

	mapper.Remove(entity)
}

func TestAddBatch(t *testing.T) {
	filter := ecs.NewFilter1[Altitude](&world).Exclusive()

	mapper.AddBatch(
		filter.Batch(),
		&Position{X: 100, Y: 100},
		&Velocity{X: 1, Y: -1},
	)
}

func TestAddBatchFn(t *testing.T) {
	filter := ecs.NewFilter1[Altitude](&world).Exclusive()

	mapper.AddBatchFn(
		filter.Batch(),
		func(e ecs.Entity, pos *Position, vel *Velocity) {
			pos.X, pos.Y = rand.Float64()*100, rand.Float64()*100
			vel.X, vel.Y = rand.NormFloat64(), rand.NormFloat64()
		})
}

func TestRemoveBatch(t *testing.T) {
	filter := ecs.NewFilter2[Position, Velocity](&world)

	mapper.RemoveBatch(
		filter.Batch(),
		func(entity ecs.Entity) { /* ... */ })
}

func TestResourcesQuick(t *testing.T) {
	grid := NewGrid(100, 100)

	ecs.AddResource(&world, &grid)
	_ = ecs.GetResource[Grid](&world)
}

func TestResources(t *testing.T) {
	gridResource := ecs.NewResource[Grid](&world)

	grid := gridResource.Get()
	_ = grid
}
