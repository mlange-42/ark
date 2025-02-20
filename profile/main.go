package main

// Profiling:
// go build ./profile
// ./profile
// go tool pprof -http=":8000" -nodefraction=0.001 ./profile cpu.pprof

import (
	"github.com/mlange-42/ark/ecs"
	"github.com/pkg/profile"
)

type position struct {
	X float64
	Y float64
}

type velocity struct {
	X float64
	Y float64
}

func main() {

	count := 10
	iters := 10000
	entities := 100000

	stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	run(count, iters, entities)
	stop.Stop()
}

func run(rounds, iters, entities int) {
	for i := 0; i < rounds; i++ {
		world := ecs.NewWorld(1024)

		posMap := ecs.NewMap[position](&world)
		velMap := ecs.NewMap[velocity](&world)

		for j := 0; j < entities; j++ {
			e := world.NewEntity()
			posMap.Add(e)
			velMap.Add(e)
		}

		query := ecs.NewQuery2[position, velocity](&world).Build()
		for j := 0; j < iters; j++ {
			for query.Next() {
				pos, vel := query.Get()
				pos.X += vel.X
				pos.Y += vel.Y
			}
		}
	}
}
