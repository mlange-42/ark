// Demonstrates how to run multiple simulations in parallel.
package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/mlange-42/ark/ecs"
)

func main() {
	totalRuns := 100
	start := time.Now()

	// As many workers as processors available.
	workers := runtime.NumCPU()
	// Channel for sending jobs to workers (buffered!).
	jobs := make(chan int, totalRuns)
	// Channel for retrieving results / done messages (buffered!).
	results := make(chan int, totalRuns)

	// Start the workers.
	for w := 0; w < workers; w++ {
		go worker(jobs, results)
	}

	// Send the jobs. Does not block due to buffered channel.
	for j := 0; j < totalRuns; j++ {
		jobs <- j
	}
	close(jobs)

	// Collect and print done messages.
	for j := 0; j < totalRuns; j++ {
		job := <-results
		fmt.Printf("Job %4d done after %6.1fms\n", job, float64(time.Since(start).Microseconds())/1000.0)
	}

	fmt.Printf("Parallel (%d): %s\n", workers, time.Since(start))
}

// Worker for running simulations on a reused world.
// Each worker needs its own world.
func worker(jobs <-chan int, results chan<- int) {
	// Create the worker's world. Will be reused for all jobs of the worker.
	w := ecs.NewWorld()

	// Process incoming jobs.
	for j := range jobs {
		// Run the model.
		runSimulation(&w)
		// Send done message. Does not block due to buffered channel.
		results <- j
	}
}

// Countdown component
type Countdown struct {
	Remaining int // Remaining ticks.
}

// A simulation that creates 10k entities, and then removes them after 1000 steps, 10 each step.
// The argument is an ECS world that is reused for multiple simulations.
func runSimulation(w *ecs.World) {
	// Reset the world.
	w.Reset()

	// Create 10k entities, each with a countdown.
	builder := ecs.NewMap1[Countdown](w)
	counter := 0
	builder.NewBatchFn(10_000, func(entity ecs.Entity, countdown *Countdown) {
		countdown.Remaining = counter/10 + 1000
		counter++
	})

	// List of entities to remove in each step.
	toRemove := []ecs.Entity{}

	// Filter for querying entities with Countdown.
	filter := ecs.NewFilter1[Countdown](w)

	// Run until there are no more entities.
	for {
		// Get a fresh query from the filter.
		query := filter.Query()

		// If no more entities, break out of the loop.
		if query.Count() == 0 {
			query.Close()
			break
		}

		// Iterate entities.
		for query.Next() {
			countdown := query.Get()
			// Countdown.
			countdown.Remaining--
			// List entity for removal if countdown hits zero.
			if countdown.Remaining <= 0 {
				toRemove = append(toRemove, query.Entity())
			}
		}

		// Remove all entities with countdown zero.
		for _, e := range toRemove {
			w.RemoveEntity(e)
		}

		// Clear list of entities to remove.
		toRemove = toRemove[:0]
	}
}
