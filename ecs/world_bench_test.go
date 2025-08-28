package ecs_test

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

func BenchmarkCreateEntity1Comp_1000(b *testing.B) {

	w := ecs.NewWorld()
	builder := ecs.NewMap1[Position](&w)
	filter := ecs.NewFilter0(&w)

	builder.NewBatchFn(1000, nil)
	w.RemoveEntities(filter.Batch(), nil)

	for b.Loop() {
		b.StartTimer()
		for range 1000 {
			_ = builder.NewEntityFn(nil)
		}
		b.StopTimer()
		w.RemoveEntities(filter.Batch(), nil)
	}
}
