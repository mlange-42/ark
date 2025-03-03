package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	s := newStorage()
	assert.Equal(t, 1, len(s.archetypes))
	assert.Equal(t, 1, len(s.tables))

	s.AddComponent(0)
	s.AddComponent(1)

	assert.Panics(t, func() { s.AddComponent(3) })
}

func TestStorageGetExchangeMask(t *testing.T) {
	s := newStorage()

	id1, _ := s.registry.ComponentID(typeOf[CompA]())
	id2, _ := s.registry.ComponentID(typeOf[CompB]())
	id3, _ := s.registry.ComponentID(typeOf[CompC]())

	mask := NewMask(ID{id1}, ID{id2})

	s.getExchangeMask(&mask, []ID{{id3}}, []ID{{id2}})
	assert.Equal(t, NewMask(ID{id1}, ID{id3}), mask)

	mask = NewMask(ID{id1}, ID{id2})
	assert.Panics(t, func() {
		s.getExchangeMask(&mask, []ID{{id1}}, []ID{})
	})
	mask = NewMask(ID{id1}, ID{id2})
	assert.Panics(t, func() {
		s.getExchangeMask(&mask, []ID{}, []ID{{id3}})
	})
}
