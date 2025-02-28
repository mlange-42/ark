//go:build !debug

package ecs

func (w *World) checkQueryNext(cursor *cursor) {}

func (w *World) checkQueryGet(cursor *cursor) {}

func (s *storage) checkHasComponent(entity Entity, comp ID) {}
