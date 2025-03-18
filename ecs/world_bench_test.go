package ecs_test

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

func BenchmarkCreateEntity1Comp_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	builder := ecs.NewMap1[Position](&w)
	filter := ecs.NewFilter0(&w)

	builder.NewBatchFn(1000, nil)
	w.RemoveEntities(filter.Batch(), nil)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for j := 0; j < 1000; j++ {
			_ = builder.NewEntityFn(nil)
		}
		b.StopTimer()
		w.RemoveEntities(filter.Batch(), nil)
	}
}
