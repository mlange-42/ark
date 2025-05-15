package main

import (
	"fmt"

	"github.com/mlange-42/ark/ecs"
)

type Data struct {
	id    int
	array []int
}

type Other struct {
	ecs.RelationMarker
}

func main() {
	nbEntities := 20

	w := ecs.NewWorld()
	dataMap := ecs.NewMap1[Data](&w)

	// Creating some entities
	entities := make([]ecs.Entity, 0)
	for i := range nbEntities {
		entities = append(entities, dataMap.NewEntity(&Data{id: i, array: []int{1, 2, 3, 4}}))
	}

	// Making sure I created the correct amount
	query := ecs.NewFilter1[Data](&w)
	q := query.Query()
	fmt.Printf("Entities created: %d\n", q.Count())
	q.Close()

	// Adding some relations between entities
	relMap := ecs.NewMap1[Other](&w)
	for i := range nbEntities {
		relMap.Add(entities[i], &Other{}, ecs.Rel[Other](entities[(i+1)%nbEntities]))
	}

	// Here, if nbEntities is large enough, new entities appear
	q = query.Query()
	fmt.Printf("Entities after relations: %d\n", q.Count())

	// Displaying the content of entities: newly created entities are empty
	i := 0
	for q.Next() {
		d := q.Get()

		fmt.Printf("i = %d, id = %d, array = %v, entity = %v\n", i, d.id, d.array, q.Entity())
		i++
	}
}
