// Demonstrates how to implement systems and a scheduler.
package main

import (
	"math/rand"

	"github.com/mlange-42/ark/ecs"
)

func main() {
	// Create a new Scheduler
	scheduler := Scheduler{}

	// Parametrize and add Systems
	scheduler.AddSystem(
		&InitializerSystem{Count: 100},
	)
	scheduler.AddSystem(
		&PosUpdaterSystem{},
	)

	// Run the model
	scheduler.Run(100)
}

// System interface
type System interface {
	Initialize(w *ecs.World)
	Update(w *ecs.World)
}

// Scheduler for updating systems
type Scheduler struct {
	world   ecs.World
	systems []System
}

// AddSystem adds a System to the scheduler
func (s *Scheduler) AddSystem(sys System) {
	s.systems = append(s.systems, sys)
}

// Run initializes and updates all systems
func (s *Scheduler) Run(steps int) {
	s.initialize()
	for range steps {
		s.update()
	}
}

// initialize a new world and all systems
func (s *Scheduler) initialize() {
	s.world = ecs.NewWorld()

	for _, sys := range s.systems {
		sys.Initialize(&s.world)
	}
}

// update all systems
func (s *Scheduler) update() {
	for _, sys := range s.systems {
		sys.Update(&s.world)
	}
}

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

// InitializerSystem to create entities
type InitializerSystem struct {
	Count int
}

// Initialize the system
func (s *InitializerSystem) Initialize(w *ecs.World) {
	// Initialize entities with Position and Velocity
	mapper := ecs.NewMap2[Position, Velocity](w)
	mapper.NewBatchFn(s.Count, func(_ ecs.Entity, pos *Position, vel *Velocity) {
		pos.X = rand.Float64() * 100
		pos.Y = rand.Float64() * 100
		vel.X = rand.NormFloat64()
		vel.Y = rand.NormFloat64()
	})
}

// Update the system
func (s *InitializerSystem) Update(w *ecs.World) {}

// PosUpdaterSystem updates entity positions
type PosUpdaterSystem struct {
	// Store filters permanently, for efficiency
	filter *ecs.Filter2[Position, Velocity]
}

// Initialize the system
func (s *PosUpdaterSystem) Initialize(w *ecs.World) {
	// Initialize the filter, using a shortcut to avoid listing generics again
	s.filter = s.filter.New(w)
}

// Update the system
func (s *PosUpdaterSystem) Update(w *ecs.World) {
	// Perform system update logic
	query := s.filter.Query()
	for query.Next() {
		pos, vel := query.Get()
		pos.X += vel.X
		pos.Y += vel.Y
	}
}
