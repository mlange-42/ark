package relations

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

var world = ecs.NewWorld()
var parent = world.NewEntity()

type Position struct {
	X float64
	Y float64
}

type Velocity struct {
	X float64
	Y float64
}

type ChildOf struct {
	ecs.RelationMarker
}

func TestNewEntity(t *testing.T) {
	// Create a component Mapper
	mapper := ecs.NewMap2[Position, ChildOf](&world)

	// Create a parent entity.
	parent := world.NewEntity()

	// Create an entity with a parent entity, the slow way.
	_ = mapper.NewEntity(&Position{}, &ChildOf{}, ecs.Rel[ChildOf](parent))
	// Create an entity with a parent entity, the fast way.
	_ = mapper.NewEntity(&Position{}, &ChildOf{}, ecs.RelIdx(1, parent))
}

func TestAdd(t *testing.T) {
	// Create a component Mapper
	mapper := ecs.NewMap2[Position, ChildOf](&world)

	// Create a parent entity.
	parent := world.NewEntity()
	// Create a child entity.
	child := world.NewEntity()

	// Add components and a relation target to it, the slow way.
	mapper.Add(child, &Position{}, &ChildOf{}, ecs.Rel[ChildOf](parent))

	// Add components and a relation target to an entity, the fast way.
	child2 := world.NewEntity()
	mapper.Add(child2, &Position{}, &ChildOf{}, ecs.RelIdx(1, parent))
}

func TestSetRelations(t *testing.T) {
	// Create a component Mapper
	mapper := ecs.NewMap2[Position, ChildOf](&world)

	// Create parent entities.
	parent1 := world.NewEntity()
	parent2 := world.NewEntity()

	// Create an entity with a parent entity.
	child := mapper.NewEntity(&Position{}, &ChildOf{}, ecs.RelIdx(1, parent1))

	// Change the entity's parent.
	mapper.SetRelations(child, ecs.RelIdx(1, parent2))

	// Change the entity's parent, the slow way.
	mapper.SetRelations(child, ecs.Rel[ChildOf](parent2))
}

func TestGetRelation(t *testing.T) {
	// Create a component Mapper
	mapper := ecs.NewMap2[Position, ChildOf](&world)

	// Create a parent entity.
	parent := world.NewEntity()

	// Create an entity with a parent entity.
	child := mapper.NewEntity(&Position{}, &ChildOf{}, ecs.RelIdx(1, parent))

	// Get a relation target by component index.
	parent = mapper.GetRelation(child, 1)
}

func TestMap(t *testing.T) {
	// Create a component Mapper
	childMap := ecs.NewMap[ChildOf](&world)

	// Create parent entities.
	parent1 := world.NewEntity()
	parent2 := world.NewEntity()
	// Create a child entity.
	child := world.NewEntity()

	// Add a component with a target parent entity.
	childMap.Add(child, &ChildOf{}, parent1)

	// Change the entity's parent.
	childMap.SetRelation(child, parent2)

	// Get the relation target.
	parent := childMap.GetRelation(child)
	_ = parent
}

func TestFilter1(t *testing.T) {
	// Create a filter with a relation target.
	filter := ecs.NewFilter2[Position, ChildOf](&world).
		Relations(ecs.Rel[ChildOf](parent))

	// Get a query for iteration.
	query := filter.Query()
	// ...
	_ = query
}
