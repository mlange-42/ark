// Demonstrates how to implement systems.
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

// AddSystem adds a System to the Scheduler
func (s *Scheduler) AddSystem(sys System) {
	s.systems = append(s.systems, sys)
}

// Run initializes and updates all Systems
func (s *Scheduler) Run(steps int) {
	s.initialize()
	for i := 0; i < steps; i++ {
		s.update()
	}
}

func (s *Scheduler) initialize() {
	s.world = ecs.NewWorld()

	for _, sys := range s.systems {
		sys.Initialize(&s.world)
	}
}

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
	filter *ecs.Filter2[Position, Velocity]
}

// Initialize the system
func (s *PosUpdaterSystem) Initialize(w *ecs.World) {
	s.filter = ecs.NewFilter2[Position, Velocity](w)
}

// Update the system
func (s *PosUpdaterSystem) Update(w *ecs.World) {
	query := s.filter.Query()
	for query.Next() {
		pos, vel := query.Get()
		pos.X += vel.X
		pos.Y += vel.Y
	}
}
