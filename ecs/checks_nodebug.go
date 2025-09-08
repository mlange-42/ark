//go:build !ark_debug

package ecs

func (s *storage) checkHasComponent(entity Entity, comp ID) {}
