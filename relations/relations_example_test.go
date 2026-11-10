package relations

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/mlange-42/ark/ecs"
)

// ####################### Components ###########################

// Farm component.
type Farm struct{ ID int }

// Weight component for animals.
type Weight struct{ Kilograms float64 }

// IsInFarm component for animals.
type IsInFarm struct{ ecs.RelationMarker }

// IsOfSpecies component for animals.
type IsOfSpecies struct{ ecs.RelationMarker }

func TestMain(t *testing.T) {

	// ####################### Preparations ###########################

	world := ecs.NewWorld()

	// Create a component mapper for farms.
	farmMap := ecs.NewMap1[Farm](world)
	// Create a component mapper for farm animals.
	animalMap := ecs.NewMap3[Weight, IsInFarm, IsOfSpecies](world)

	// Create a filter for farms.
	farmFilter := ecs.NewFilter1[Farm](world)
	// Create a filter for farm animals.
	animalFilter := ecs.NewFilter3[Weight, IsInFarm, IsOfSpecies](world)

	// ####################### Initialization ###########################

	// Create species.
	cow := world.NewEntity()
	pig := world.NewEntity()

	// Create farms.
	farms := []ecs.Entity{}
	for i := range 10 {
		farm := farmMap.NewEntity(&Farm{i})
		farms = append(farms, farm)
	}

	// Populate farms.
	for _, farm := range farms {
		// How many animals?
		numCows := rand.IntN(50)
		numPigs := rand.IntN(200)

		// Create cows.
		animalMap.NewBatch(numCows, // How many?
			&Weight{500}, &IsInFarm{}, &IsOfSpecies{}, // Initial values.
			ecs.Rel[IsInFarm](farm),   // This farm.
			ecs.Rel[IsOfSpecies](cow), // Species cow.
		)
		// Create pigs.
		animalMap.NewBatch(numPigs, // How many?
			&Weight{100}, &IsInFarm{}, &IsOfSpecies{}, // Initial values.
			ecs.Rel[IsInFarm](farm),   // This farm.
			ecs.Rel[IsOfSpecies](pig), // Species pig.
		)
	}

	// ####################### Logic in systems ###########################

	// Do something with all pigs.
	query := animalFilter.Query(ecs.Rel[IsOfSpecies](pig))
	for query.Next() {
		weight, _, _ := query.Get()
		weight.Kilograms += rand.Float64() * 10
	}

	// Print total weight of the pigs in each farm.
	// Iterate farms.
	farmQuery := farmFilter.Query()
	for farmQuery.Next() {
		farm := farmQuery.Get()
		farmEntity := farmQuery.Entity()

		totalWeight := 0.0

		// Iterate pigs in the farm.
		animalQuery := animalFilter.Query(
			ecs.Rel[IsInFarm](farmEntity), // This farm.
			ecs.Rel[IsOfSpecies](pig),     // Pigs only.
		)
		for animalQuery.Next() {
			weight, _, _ := animalQuery.Get()
			totalWeight += weight.Kilograms
		}
		// Print the farm's result.
		fmt.Printf("Farm %d: %.0fkg\n", farm.ID, totalWeight)
	}
}
