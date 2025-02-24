//go:build !debug

package ecs

func checkQueryNext(cursor *cursor) {}

func checkQueryGet(cursor *cursor) {}

func checkMapHasComponent(comp *componentStorage, table tableID) {}
