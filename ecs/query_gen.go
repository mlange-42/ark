package ecs

// Code generated by go generate; DO NOT EDIT.

type cursor struct {
	table    int
	index    uintptr
	maxIndex int64
}

// Query0 is a filter for two components.
type Query0 struct {
	world      *World
	mask       Mask
	cursor     cursor
	table      *table
	components []*componentStorage
}

// Next advances the query's cursor to the next entity.
func (q *Query0) Next() bool {
	if int64(q.cursor.index) < q.cursor.maxIndex {
		q.cursor.index++
		return true
	}
	return q.nextTable()
}

// Entity returns the current entity.
func (q *Query0) Entity() Entity {
	return q.table.GetEntity(q.cursor.index)
}

// Get returns the queries components of the current entity.
func (q *Query0) Get() {
	return
}

func (q *Query0) nextTable() bool {
	maxTableIndex := len(q.world.storage.tables) - 1
	for q.cursor.table < maxTableIndex {
		q.cursor.table++
		q.table = &q.world.storage.tables[q.cursor.table]
		archetype := &q.world.storage.archetypes[q.table.archetype]
		if !archetype.mask.Contains(&q.mask) || q.table.entities.Len() == 0 {
			continue
		}

		q.cursor.index = 0
		q.cursor.maxIndex = int64(q.table.entities.Len() - 1)
		return true
	}
	q.cursor.table = -1
	q.cursor.index = 0
	q.cursor.maxIndex = -1
	q.table = nil
	return false
}

// Query1 is a filter for two components.
type Query1[A any] struct {
	world      *World
	mask       Mask
	cursor     cursor
	table      *table
	components []*componentStorage
	columnA    *column
}

// Next advances the query's cursor to the next entity.
func (q *Query1[A]) Next() bool {
	if int64(q.cursor.index) < q.cursor.maxIndex {
		q.cursor.index++
		return true
	}
	return q.nextTable()
}

// Entity returns the current entity.
func (q *Query1[A]) Entity() Entity {
	return q.table.GetEntity(q.cursor.index)
}

// Get returns the queries components of the current entity.
func (q *Query1[A]) Get() *A {
	return (*A)(q.columnA.Get(q.cursor.index))
}

func (q *Query1[A]) nextTable() bool {
	maxTableIndex := len(q.world.storage.tables) - 1
	for q.cursor.table < maxTableIndex {
		q.cursor.table++
		q.table = &q.world.storage.tables[q.cursor.table]
		archetype := &q.world.storage.archetypes[q.table.archetype]
		if !archetype.mask.Contains(&q.mask) || q.table.entities.Len() == 0 {
			continue
		}
		q.columnA = q.components[0].GetColumn(tableID(q.cursor.table))

		q.cursor.index = 0
		q.cursor.maxIndex = int64(q.table.entities.Len() - 1)
		return true
	}
	q.cursor.table = -1
	q.cursor.index = 0
	q.cursor.maxIndex = -1
	q.table = nil
	return false
}

// Query2 is a filter for two components.
type Query2[A any, B any] struct {
	world      *World
	mask       Mask
	cursor     cursor
	table      *table
	components []*componentStorage
	columnA    *column
	columnB    *column
}

// Next advances the query's cursor to the next entity.
func (q *Query2[A, B]) Next() bool {
	if int64(q.cursor.index) < q.cursor.maxIndex {
		q.cursor.index++
		return true
	}
	return q.nextTable()
}

// Entity returns the current entity.
func (q *Query2[A, B]) Entity() Entity {
	return q.table.GetEntity(q.cursor.index)
}

// Get returns the queries components of the current entity.
func (q *Query2[A, B]) Get() (*A, *B) {
	return (*A)(q.columnA.Get(q.cursor.index)),
		(*B)(q.columnB.Get(q.cursor.index))
}

func (q *Query2[A, B]) nextTable() bool {
	maxTableIndex := len(q.world.storage.tables) - 1
	for q.cursor.table < maxTableIndex {
		q.cursor.table++
		q.table = &q.world.storage.tables[q.cursor.table]
		archetype := &q.world.storage.archetypes[q.table.archetype]
		if !archetype.mask.Contains(&q.mask) || q.table.entities.Len() == 0 {
			continue
		}
		q.columnA = q.components[0].GetColumn(tableID(q.cursor.table))
		q.columnB = q.components[1].GetColumn(tableID(q.cursor.table))

		q.cursor.index = 0
		q.cursor.maxIndex = int64(q.table.entities.Len() - 1)
		return true
	}
	q.cursor.table = -1
	q.cursor.index = 0
	q.cursor.maxIndex = -1
	q.table = nil
	return false
}

// Query3 is a filter for two components.
type Query3[A any, B any, C any] struct {
	world      *World
	mask       Mask
	cursor     cursor
	table      *table
	components []*componentStorage
	columnA    *column
	columnB    *column
	columnC    *column
}

// Next advances the query's cursor to the next entity.
func (q *Query3[A, B, C]) Next() bool {
	if int64(q.cursor.index) < q.cursor.maxIndex {
		q.cursor.index++
		return true
	}
	return q.nextTable()
}

// Entity returns the current entity.
func (q *Query3[A, B, C]) Entity() Entity {
	return q.table.GetEntity(q.cursor.index)
}

// Get returns the queries components of the current entity.
func (q *Query3[A, B, C]) Get() (*A, *B, *C) {
	return (*A)(q.columnA.Get(q.cursor.index)),
		(*B)(q.columnB.Get(q.cursor.index)),
		(*C)(q.columnC.Get(q.cursor.index))
}

func (q *Query3[A, B, C]) nextTable() bool {
	maxTableIndex := len(q.world.storage.tables) - 1
	for q.cursor.table < maxTableIndex {
		q.cursor.table++
		q.table = &q.world.storage.tables[q.cursor.table]
		archetype := &q.world.storage.archetypes[q.table.archetype]
		if !archetype.mask.Contains(&q.mask) || q.table.entities.Len() == 0 {
			continue
		}
		q.columnA = q.components[0].GetColumn(tableID(q.cursor.table))
		q.columnB = q.components[1].GetColumn(tableID(q.cursor.table))
		q.columnC = q.components[2].GetColumn(tableID(q.cursor.table))

		q.cursor.index = 0
		q.cursor.maxIndex = int64(q.table.entities.Len() - 1)
		return true
	}
	q.cursor.table = -1
	q.cursor.index = 0
	q.cursor.maxIndex = -1
	q.table = nil
	return false
}

// Query4 is a filter for two components.
type Query4[A any, B any, C any, D any] struct {
	world      *World
	mask       Mask
	cursor     cursor
	table      *table
	components []*componentStorage
	columnA    *column
	columnB    *column
	columnC    *column
	columnD    *column
}

// Next advances the query's cursor to the next entity.
func (q *Query4[A, B, C, D]) Next() bool {
	if int64(q.cursor.index) < q.cursor.maxIndex {
		q.cursor.index++
		return true
	}
	return q.nextTable()
}

// Entity returns the current entity.
func (q *Query4[A, B, C, D]) Entity() Entity {
	return q.table.GetEntity(q.cursor.index)
}

// Get returns the queries components of the current entity.
func (q *Query4[A, B, C, D]) Get() (*A, *B, *C, *D) {
	return (*A)(q.columnA.Get(q.cursor.index)),
		(*B)(q.columnB.Get(q.cursor.index)),
		(*C)(q.columnC.Get(q.cursor.index)),
		(*D)(q.columnD.Get(q.cursor.index))
}

func (q *Query4[A, B, C, D]) nextTable() bool {
	maxTableIndex := len(q.world.storage.tables) - 1
	for q.cursor.table < maxTableIndex {
		q.cursor.table++
		q.table = &q.world.storage.tables[q.cursor.table]
		archetype := &q.world.storage.archetypes[q.table.archetype]
		if !archetype.mask.Contains(&q.mask) || q.table.entities.Len() == 0 {
			continue
		}
		q.columnA = q.components[0].GetColumn(tableID(q.cursor.table))
		q.columnB = q.components[1].GetColumn(tableID(q.cursor.table))
		q.columnC = q.components[2].GetColumn(tableID(q.cursor.table))
		q.columnD = q.components[3].GetColumn(tableID(q.cursor.table))

		q.cursor.index = 0
		q.cursor.maxIndex = int64(q.table.entities.Len() - 1)
		return true
	}
	q.cursor.table = -1
	q.cursor.index = 0
	q.cursor.maxIndex = -1
	q.table = nil
	return false
}

// Query5 is a filter for two components.
type Query5[A any, B any, C any, D any, E any] struct {
	world      *World
	mask       Mask
	cursor     cursor
	table      *table
	components []*componentStorage
	columnA    *column
	columnB    *column
	columnC    *column
	columnD    *column
	columnE    *column
}

// Next advances the query's cursor to the next entity.
func (q *Query5[A, B, C, D, E]) Next() bool {
	if int64(q.cursor.index) < q.cursor.maxIndex {
		q.cursor.index++
		return true
	}
	return q.nextTable()
}

// Entity returns the current entity.
func (q *Query5[A, B, C, D, E]) Entity() Entity {
	return q.table.GetEntity(q.cursor.index)
}

// Get returns the queries components of the current entity.
func (q *Query5[A, B, C, D, E]) Get() (*A, *B, *C, *D, *E) {
	return (*A)(q.columnA.Get(q.cursor.index)),
		(*B)(q.columnB.Get(q.cursor.index)),
		(*C)(q.columnC.Get(q.cursor.index)),
		(*D)(q.columnD.Get(q.cursor.index)),
		(*E)(q.columnE.Get(q.cursor.index))
}

func (q *Query5[A, B, C, D, E]) nextTable() bool {
	maxTableIndex := len(q.world.storage.tables) - 1
	for q.cursor.table < maxTableIndex {
		q.cursor.table++
		q.table = &q.world.storage.tables[q.cursor.table]
		archetype := &q.world.storage.archetypes[q.table.archetype]
		if !archetype.mask.Contains(&q.mask) || q.table.entities.Len() == 0 {
			continue
		}
		q.columnA = q.components[0].GetColumn(tableID(q.cursor.table))
		q.columnB = q.components[1].GetColumn(tableID(q.cursor.table))
		q.columnC = q.components[2].GetColumn(tableID(q.cursor.table))
		q.columnD = q.components[3].GetColumn(tableID(q.cursor.table))
		q.columnE = q.components[4].GetColumn(tableID(q.cursor.table))

		q.cursor.index = 0
		q.cursor.maxIndex = int64(q.table.entities.Len() - 1)
		return true
	}
	q.cursor.table = -1
	q.cursor.index = 0
	q.cursor.maxIndex = -1
	q.table = nil
	return false
}

// Query6 is a filter for two components.
type Query6[A any, B any, C any, D any, E any, F any] struct {
	world      *World
	mask       Mask
	cursor     cursor
	table      *table
	components []*componentStorage
	columnA    *column
	columnB    *column
	columnC    *column
	columnD    *column
	columnE    *column
	columnF    *column
}

// Next advances the query's cursor to the next entity.
func (q *Query6[A, B, C, D, E, F]) Next() bool {
	if int64(q.cursor.index) < q.cursor.maxIndex {
		q.cursor.index++
		return true
	}
	return q.nextTable()
}

// Entity returns the current entity.
func (q *Query6[A, B, C, D, E, F]) Entity() Entity {
	return q.table.GetEntity(q.cursor.index)
}

// Get returns the queries components of the current entity.
func (q *Query6[A, B, C, D, E, F]) Get() (*A, *B, *C, *D, *E, *F) {
	return (*A)(q.columnA.Get(q.cursor.index)),
		(*B)(q.columnB.Get(q.cursor.index)),
		(*C)(q.columnC.Get(q.cursor.index)),
		(*D)(q.columnD.Get(q.cursor.index)),
		(*E)(q.columnE.Get(q.cursor.index)),
		(*F)(q.columnF.Get(q.cursor.index))
}

func (q *Query6[A, B, C, D, E, F]) nextTable() bool {
	maxTableIndex := len(q.world.storage.tables) - 1
	for q.cursor.table < maxTableIndex {
		q.cursor.table++
		q.table = &q.world.storage.tables[q.cursor.table]
		archetype := &q.world.storage.archetypes[q.table.archetype]
		if !archetype.mask.Contains(&q.mask) || q.table.entities.Len() == 0 {
			continue
		}
		q.columnA = q.components[0].GetColumn(tableID(q.cursor.table))
		q.columnB = q.components[1].GetColumn(tableID(q.cursor.table))
		q.columnC = q.components[2].GetColumn(tableID(q.cursor.table))
		q.columnD = q.components[3].GetColumn(tableID(q.cursor.table))
		q.columnE = q.components[4].GetColumn(tableID(q.cursor.table))
		q.columnF = q.components[5].GetColumn(tableID(q.cursor.table))

		q.cursor.index = 0
		q.cursor.maxIndex = int64(q.table.entities.Len() - 1)
		return true
	}
	q.cursor.table = -1
	q.cursor.index = 0
	q.cursor.maxIndex = -1
	q.table = nil
	return false
}

// Query7 is a filter for two components.
type Query7[A any, B any, C any, D any, E any, F any, G any] struct {
	world      *World
	mask       Mask
	cursor     cursor
	table      *table
	components []*componentStorage
	columnA    *column
	columnB    *column
	columnC    *column
	columnD    *column
	columnE    *column
	columnF    *column
	columnG    *column
}

// Next advances the query's cursor to the next entity.
func (q *Query7[A, B, C, D, E, F, G]) Next() bool {
	if int64(q.cursor.index) < q.cursor.maxIndex {
		q.cursor.index++
		return true
	}
	return q.nextTable()
}

// Entity returns the current entity.
func (q *Query7[A, B, C, D, E, F, G]) Entity() Entity {
	return q.table.GetEntity(q.cursor.index)
}

// Get returns the queries components of the current entity.
func (q *Query7[A, B, C, D, E, F, G]) Get() (*A, *B, *C, *D, *E, *F, *G) {
	return (*A)(q.columnA.Get(q.cursor.index)),
		(*B)(q.columnB.Get(q.cursor.index)),
		(*C)(q.columnC.Get(q.cursor.index)),
		(*D)(q.columnD.Get(q.cursor.index)),
		(*E)(q.columnE.Get(q.cursor.index)),
		(*F)(q.columnF.Get(q.cursor.index)),
		(*G)(q.columnG.Get(q.cursor.index))
}

func (q *Query7[A, B, C, D, E, F, G]) nextTable() bool {
	maxTableIndex := len(q.world.storage.tables) - 1
	for q.cursor.table < maxTableIndex {
		q.cursor.table++
		q.table = &q.world.storage.tables[q.cursor.table]
		archetype := &q.world.storage.archetypes[q.table.archetype]
		if !archetype.mask.Contains(&q.mask) || q.table.entities.Len() == 0 {
			continue
		}
		q.columnA = q.components[0].GetColumn(tableID(q.cursor.table))
		q.columnB = q.components[1].GetColumn(tableID(q.cursor.table))
		q.columnC = q.components[2].GetColumn(tableID(q.cursor.table))
		q.columnD = q.components[3].GetColumn(tableID(q.cursor.table))
		q.columnE = q.components[4].GetColumn(tableID(q.cursor.table))
		q.columnF = q.components[5].GetColumn(tableID(q.cursor.table))
		q.columnG = q.components[6].GetColumn(tableID(q.cursor.table))

		q.cursor.index = 0
		q.cursor.maxIndex = int64(q.table.entities.Len() - 1)
		return true
	}
	q.cursor.table = -1
	q.cursor.index = 0
	q.cursor.maxIndex = -1
	q.table = nil
	return false
}

// Query8 is a filter for two components.
type Query8[A any, B any, C any, D any, E any, F any, G any, H any] struct {
	world      *World
	mask       Mask
	cursor     cursor
	table      *table
	components []*componentStorage
	columnA    *column
	columnB    *column
	columnC    *column
	columnD    *column
	columnE    *column
	columnF    *column
	columnG    *column
	columnH    *column
}

// Next advances the query's cursor to the next entity.
func (q *Query8[A, B, C, D, E, F, G, H]) Next() bool {
	if int64(q.cursor.index) < q.cursor.maxIndex {
		q.cursor.index++
		return true
	}
	return q.nextTable()
}

// Entity returns the current entity.
func (q *Query8[A, B, C, D, E, F, G, H]) Entity() Entity {
	return q.table.GetEntity(q.cursor.index)
}

// Get returns the queries components of the current entity.
func (q *Query8[A, B, C, D, E, F, G, H]) Get() (*A, *B, *C, *D, *E, *F, *G, *H) {
	return (*A)(q.columnA.Get(q.cursor.index)),
		(*B)(q.columnB.Get(q.cursor.index)),
		(*C)(q.columnC.Get(q.cursor.index)),
		(*D)(q.columnD.Get(q.cursor.index)),
		(*E)(q.columnE.Get(q.cursor.index)),
		(*F)(q.columnF.Get(q.cursor.index)),
		(*G)(q.columnG.Get(q.cursor.index)),
		(*H)(q.columnH.Get(q.cursor.index))
}

func (q *Query8[A, B, C, D, E, F, G, H]) nextTable() bool {
	maxTableIndex := len(q.world.storage.tables) - 1
	for q.cursor.table < maxTableIndex {
		q.cursor.table++
		q.table = &q.world.storage.tables[q.cursor.table]
		archetype := &q.world.storage.archetypes[q.table.archetype]
		if !archetype.mask.Contains(&q.mask) || q.table.entities.Len() == 0 {
			continue
		}
		q.columnA = q.components[0].GetColumn(tableID(q.cursor.table))
		q.columnB = q.components[1].GetColumn(tableID(q.cursor.table))
		q.columnC = q.components[2].GetColumn(tableID(q.cursor.table))
		q.columnD = q.components[3].GetColumn(tableID(q.cursor.table))
		q.columnE = q.components[4].GetColumn(tableID(q.cursor.table))
		q.columnF = q.components[5].GetColumn(tableID(q.cursor.table))
		q.columnG = q.components[6].GetColumn(tableID(q.cursor.table))
		q.columnH = q.components[7].GetColumn(tableID(q.cursor.table))

		q.cursor.index = 0
		q.cursor.maxIndex = int64(q.table.entities.Len() - 1)
		return true
	}
	q.cursor.table = -1
	q.cursor.index = 0
	q.cursor.maxIndex = -1
	q.table = nil
	return false
}
