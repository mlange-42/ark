package inlining

import (
	"runtime"

	"github.com/mlange-42/ark/ecs"
)

type position struct {
	X, Y float64
}

func main() {
	w := ecs.NewWorld()

	builder := ecs.NewMap1[position](&w)
	e := builder.NewEntity(&position{})
	runtime.KeepAlive(e)
}
