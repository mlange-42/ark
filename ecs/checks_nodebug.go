//go:build !debug

package ecs

func (c *cursor) checkQueryNext() {}

func (c *cursor) checkQueryGet() {}

func (s *storage) checkHasComponent(entity Entity, comp ID) {}
