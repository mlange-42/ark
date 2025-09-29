//go:build !ark_debug

package ecs

func (s *storage) checkHasComponent(_ Entity, _ ID) {}
