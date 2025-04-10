package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeOf(t *testing.T) {
	posType := typeOf[Position]()
	assert.Equal(t, "Position", posType.Name())
}

func TestPagedSlice(t *testing.T) {
	a := pagedSlice[int32]{}

	var i int32
	for i = range 66 {
		a.Add(i)
		assert.Equal(t, i, *a.Get(i))
		assert.Equal(t, i+1, a.Len())
	}

	a.Set(3, 100)
	assert.Equal(t, int32(100), *a.Get(3))
}

func BenchmarkTypeOf(b *testing.B) {
	for b.Loop() {
		_ = typeOf[Position]()
	}
}

func BenchmarkSizeOf(b *testing.B) {
	tp := typeOf[Position]()
	for b.Loop() {
		_ = sizeOf(tp)
	}
}
