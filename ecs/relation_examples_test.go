package ecs_test

import (
	"fmt"

	"github.com/mlange-42/ark/ecs"
)

// ChildOf demonstrates how to define a relation component.
// There may be more fields _after_ the embed.
type ChildOf struct {
	ecs.RelationMarker
}

func ExampleRelationMarker() {
	world := ecs.NewWorld()
	// Create a targets/parents entity for the relation.
	parent1 := world.NewEntity()
	parent2 := world.NewEntity()

	// Create a mapper for one or more components.
	mapper := ecs.NewMap1[ChildOf](&world)

	// Create a child entity with a relation to a parent.
	child1 := mapper.NewEntity(&ChildOf{}, ecs.RelIdx(0, parent1))
	// Create another child entity with a relation to a parent.
	// This version is slower than the above variant, but safer.
	_ = mapper.NewEntity(&ChildOf{}, ecs.Rel[ChildOf](parent2))

	// Set the relation target.
	mapper.SetRelations(child1, ecs.RelIdx(0, parent2))

	// Filter for the relation with a given target.
	filter := ecs.NewFilter1[ChildOf](&world)
	query := filter.Query(ecs.RelIdx(0, parent2))
	for query.Next() {
		fmt.Println(
			query.Entity(),
			query.GetRelation(0), // Get the relation target in a query
		)
	}
	// Output: {5 0} {3 0}
	// {4 0} {3 0}
}
