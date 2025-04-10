package main

// Profiling:
// go build ./profile/exchange
// ./exchange
// go tool pprof -http=":8000" -nodefraction=0.001 ./exchange cpu.pprof

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
	count := 10
	iters := 10000
	entities := 1000

	stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	run(count, iters, entities)
	stop.Stop()
}

func run(rounds, iters, numEntities int) {
	for range rounds {
		w := ecs.NewWorld(1024)

		ex1 := ecs.NewExchange1[comp1](&w).Removes(ecs.C[comp2]())
		ex2 := ecs.NewExchange1[comp2](&w).Removes(ecs.C[comp1]())
		builder := ecs.NewMap1[comp1](&w)

		entities := make([]ecs.Entity, 0, numEntities)
		builder.NewBatchFn(numEntities, func(entity ecs.Entity, a *comp1) {
			entities = append(entities, entity)
		})

		c1 := comp1{}
		c2 := comp2{}

		for range iters {
			for _, e := range entities {
				ex2.Exchange(e, &c2)
			}
			for _, e := range entities {
				ex1.Exchange(e, &c1)
			}
		}
	}
}
