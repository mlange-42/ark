package ecs

func removeBatch(world *World, batch *Batch, ids []ID, fn func(entity Entity)) {
	var process func(tableID tableID, start, len int)
	if fn != nil {
		process = func(tableID tableID, start, len int) {
			table := world.storage.tables[tableID]

			lock := world.lock()
			for i := range len {
				index := uintptr(start + i)
				fn(table.GetEntity(index))
			}
			world.unlock(lock)
		}
	}
	world.exchangeBatch(batch, nil, ids, nil, process)
}

func setRelationsBatch(world *World, batch *Batch, fn func(entity Entity), relations []RelationID) {
	var process func(tableID tableID, start, len int)
	if fn != nil {
		process = func(tableID tableID, start, len int) {
			table := world.storage.tables[tableID]

			lock := world.lock()
			for i := range len {
				index := uintptr(start + i)
				fn(table.GetEntity(index))
			}
			world.unlock(lock)
		}
	}
	world.setRelationsBatch(batch, relations, process)
}
