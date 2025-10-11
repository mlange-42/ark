package ecs

// removeBatch performs batch-removal of components.
func removeBatch(world *World, batch *Batch, ids []ID, fn func(entity Entity)) {
	var process func(tableID tableID, start, len uint32)
	if fn != nil {
		process = func(tableID tableID, start, len uint32) {
			table := world.storage.tables[tableID]

			for i := range len {
				index := uintptr(start + i)
				fn(table.GetEntity(index))
			}
		}
	}
	world.exchangeBatch(batch, nil, ids, nil, process)
}

// setRelationsBatch performs batch relation changes.
func setRelationsBatch(world *World, batch *Batch, fn func(entity Entity), relations []relationID) {
	var process func(tableID tableID, start, len int)
	if fn != nil {
		process = func(tableID tableID, start, len int) {
			table := world.storage.tables[tableID]

			for i := range len {
				index := uintptr(start + i)
				fn(table.GetEntity(index))
			}
		}
	}
	world.setRelationsBatch(batch, relations, process)
}
