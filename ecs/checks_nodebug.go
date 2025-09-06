//go:build !debug

package ecs

func (s *storage) checkHasComponent(entity Entity, comp ID) {}
