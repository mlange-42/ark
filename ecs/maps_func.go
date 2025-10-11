package ecs

func removeBatch(world *World, batch *Batch, ids []ID, fn func(entity Entity)) {
	var process func(table *table, start, len uint32)
	if fn != nil {
		process = func(table *table, start, len uint32) {
			for i := range len {
				index := uintptr(start + i)
				fn(table.GetEntity(index))
			}
		}
	}
	world.exchangeBatch(batch, nil, ids, nil, process)
}

func setRelationsBatch(world *World, batch *Batch, fn func(entity Entity), relations []relationID) {
	var process func(table *table, start, len int)
	if fn != nil {
		process = func(table *table, start, len int) {
			for i := range len {
				index := uintptr(start + i)
				fn(table.GetEntity(index))
			}
		}
	}
	world.setRelationsBatch(batch, relations, process)
}
