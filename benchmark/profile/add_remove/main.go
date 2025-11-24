package main

// Profiling:
// go build ./profile/add_remove
// ./add_remove
// go tool pprof -http=":8000" -nodefraction=0.001 ./add_remove cpu.pprof

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
	count := 20
	iters := 10000
	entities := 1000

	stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	run(count, iters, entities)
	stop.Stop()
}

func run(rounds, iters, numEntities int) {
	for range rounds {
		w := ecs.NewWorld(1024)

		map1 := ecs.NewMap1[comp1](w)
		map2 := ecs.NewMap1[comp2](w)

		entities := make([]ecs.Entity, 0, numEntities)
		map1.NewBatchFn(numEntities, func(entity ecs.Entity, a *comp1) {
			entities = append(entities, entity)
		})
		for _, e := range entities {
			map2.AddFn(e, nil)
		}
		for _, e := range entities {
			map2.Remove(e)
		}

		for range iters {
			for _, e := range entities {
				map2.AddFn(e, nil)
			}
			for _, e := range entities {
				map2.Remove(e)
			}
		}
	}
}
