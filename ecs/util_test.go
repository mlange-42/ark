package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeOf(t *testing.T) {
	posType := typeOf[Position]()
	assert.Equal(t, "Position", posType.Name())
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
