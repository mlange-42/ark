//go:build !debug

package ecs

func (w *World) checkQueryNext(cursor *cursor) {}

func (w *World) checkQueryGet(cursor *cursor) {}

func (w *World) checkMapHasComponent(comp *componentStorage, table tableID) {}
