// Demonstrates how to run queries in parallel.
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/mlange-42/ark/ecs"
)

// Position component
type Position struct {
	X float64
	Y float64
}

// Velocity component
type Velocity struct {
	X float64
	Y float64
}

// InProcess component for assignment to parallel queries
// Uses entity relationships
type InProcess struct {
	ecs.RelationMarker
}

// Total number of entities
const numEntities = 1_000_000

// Number of parallel CPU processes to use
const numProc = 4

func main() {
	// Create a world
	world := ecs.NewWorld()

	// Create a builder for entities
	// Entities have a relation InProcess to assign them to processes/queries
	builder := ecs.NewMap3[Position, Velocity, InProcess](&world)

	// Create entities
	processes := []ecs.Entity{}
	for range numProc {
		// One "parent" entity per process
		procEntity := world.NewEntity()

		// Create entities and assign them to the current "parent" process entity
		builder.NewBatch(numEntities/numProc,
			&Position{}, &Velocity{X: 1, Y: 0}, // Usual components
			&InProcess{},                   // Component for the relation
			ecs.Rel[InProcess](procEntity)) // relation target, i.e. process

		// Store process entities in a slice
		processes = append(processes, procEntity)
	}

	// Create a separate filter for each process
	// Each one filters for a particular "parent" process entity
	filters := []*ecs.Filter2[Position, Velocity]{}
	for _, procEntity := range processes {
		filters = append(filters,
			ecs.NewFilter2[Position, Velocity](&world). // Filter for the usual components
									With(ecs.C[InProcess]()).                  // Relation required, but not accessed
									Relations(ecs.Rel[InProcess](procEntity)), // Relation target
		)
	}

	// Take starting time
	start := time.Now()

	// Time loop
	iterations := 1000
	for range iterations {
		// Set up a WaitGroup to wait for queries to complete
		var wg sync.WaitGroup
		wg.Add(numProc)

		// Start a goroutine for each process, passing the resp. filter
		for _, f := range filters {
			// Actual query iteration, see below
			go runQuery(f, &wg)
		}

		// Wait for the queries to complete
		wg.Wait()
	}

	// Print elapsed time
	fmt.Printf("%s per iteration with %d entities", time.Since(start)/time.Duration(iterations), numEntities)
}

// The actual query iteration, executed numProc times in parallel
func runQuery(f *ecs.Filter2[Position, Velocity], wg *sync.WaitGroup) {
	// Defer signalling completion to the WaitGroup
	defer wg.Done()

	// Get a fresh query from the filter
	query := f.Query()
	// Do the usual iteration
	for query.Next() {
		pos, vel := query.Get()
		pos.X += vel.X
		pos.Y += vel.Y
	}
}
