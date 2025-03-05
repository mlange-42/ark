package main

// Profiling:
// go build ./profile/newquery
// ./newquery
// go tool pprof -http=":8000" -nodefraction=0.001 ./newquery cpu.pprof

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

	count := 1000
	iters := 100000

	stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	run(count, iters)
	stop.Stop()
}

func run(rounds, iters int) {
	for i := 0; i < rounds; i++ {
		world := ecs.NewWorld(1024)

		mapper := ecs.NewMap2[position, velocity](&world)
		mapper.NewBatchFn(100, nil)

		filter := ecs.NewFilter2[position, velocity](&world)
		for j := 0; j < iters; j++ {
			query := filter.Query()
			query.Close()
		}
	}
}
