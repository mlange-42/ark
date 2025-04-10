package main

// Profiling:
// go build ./profile/entities
// ./entities
// go tool pprof -http=":8000" -nodefraction=0.001 ./entities cpu.pprof

import (
	"github.com/mlange-42/ark/ecs"
	"github.com/pkg/profile"
)

type comp1 struct {
	V int64
	W int64
}

type comp2 struct {
	V int64
	W int64
}

func main() {
	count := 50
	iters := 10000
	entities := 1000

	stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	run(count, iters, entities)
	stop.Stop()
}

func run(rounds, iters, numEntities int) {
	for range rounds {
		w := ecs.NewWorld(1024)

		builder := ecs.NewMap1[comp1](&w)
		filter := ecs.NewFilter1[comp1](&w)

		for range iters {
			for range numEntities {
				builder.NewEntityFn(nil)
			}
			w.RemoveEntities(filter.Batch(), nil)
		}
	}
}
