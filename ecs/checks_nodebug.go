//go:build !ark_debug

package ecs

const isDebug = false

func (s *storage) checkHasComponent(_ Entity, _ ID) {}
