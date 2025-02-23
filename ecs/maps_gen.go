package ecs

// Code generated by go generate; DO NOT EDIT.

import "unsafe"

// Map1 is a mapper to access 1 components of an entity.
type Map1[A any] struct {
	world     *World
	ids       []ID
	storageA  *componentStorage
	relations []relationID
}

// NewMap1 creates a new [Map1].
func NewMap1[A any](world *World) Map1[A] {
	ids := []ID{
		ComponentID[A](world),
	}
	return Map1[A]{
		world:    world,
		ids:      ids,
		storageA: &world.storage.components[ids[0].id],
	}
}

// NewEntity creates a new entity with the mapped components.
func (m *Map1[A]) NewEntity(a *A, rel ...RelationIndex) Entity {
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	return m.world.newEntityWith(m.ids, []unsafe.Pointer{
		unsafe.Pointer(a),
	}, m.relations)
}

// Get returns the mapped components for the given entity.
func (m *Map1[A]) Get(entity Entity) *A {
	if !m.world.Alive(entity) {
		panic("can't get components of a dead entity")
	}
	return m.GetUnchecked(entity)
}

// GetUnchecked returns the mapped components for the given entity.
// It does not check whether the entity is alive.
// Can be used as an optimization when it is certain that the entity is alive.
func (m *Map1[A]) GetUnchecked(entity Entity) *A {
	index := m.world.entities[entity.id]
	row := uintptr(index.row)
	return (*A)(m.storageA.columns[index.table].Get(row))
}

// HasAll return whether the given entity has all mapped components.
func (m *Map1[A]) HasAll(entity Entity) bool {
	if !m.world.Alive(entity) {
		panic("can't check components of a dead entity")
	}
	index := m.world.entities[entity.id]
	return m.storageA.columns[index.table] != nil
}

// Add the mapped components to the given entity.
func (m *Map1[A]) Add(entity Entity, a *A, rel ...RelationIndex) {
	if !m.world.Alive(entity) {
		panic("can't add components to a dead entity")
	}
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	m.world.exchange(entity, m.ids, nil, []unsafe.Pointer{
		unsafe.Pointer(a),
	}, m.relations)
}

// Remove the mapped components from the given entity.
func (m *Map1[A]) Remove(entity Entity) {
	if !m.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	m.world.exchange(entity, nil, m.ids, nil, nil)
}

// Remove the mapped components from the given entity.
func (m *Map1[A]) SetRelations(entity Entity, rel ...RelationIndex) {
	if !m.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	m.world.setRelations(entity, m.relations)
}

// Map2 is a mapper to access 2 components of an entity.
type Map2[A any, B any] struct {
	world     *World
	ids       []ID
	storageA  *componentStorage
	storageB  *componentStorage
	relations []relationID
}

// NewMap2 creates a new [Map2].
func NewMap2[A any, B any](world *World) Map2[A, B] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
	}
	return Map2[A, B]{
		world:    world,
		ids:      ids,
		storageA: &world.storage.components[ids[0].id],
		storageB: &world.storage.components[ids[1].id],
	}
}

// NewEntity creates a new entity with the mapped components.
func (m *Map2[A, B]) NewEntity(a *A, b *B, rel ...RelationIndex) Entity {
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	return m.world.newEntityWith(m.ids, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
	}, m.relations)
}

// Get returns the mapped components for the given entity.
func (m *Map2[A, B]) Get(entity Entity) (*A, *B) {
	if !m.world.Alive(entity) {
		panic("can't get components of a dead entity")
	}
	return m.GetUnchecked(entity)
}

// GetUnchecked returns the mapped components for the given entity.
// It does not check whether the entity is alive.
// Can be used as an optimization when it is certain that the entity is alive.
func (m *Map2[A, B]) GetUnchecked(entity Entity) (*A, *B) {
	index := m.world.entities[entity.id]
	row := uintptr(index.row)
	return (*A)(m.storageA.columns[index.table].Get(row)),
		(*B)(m.storageB.columns[index.table].Get(row))
}

// HasAll return whether the given entity has all mapped components.
func (m *Map2[A, B]) HasAll(entity Entity) bool {
	if !m.world.Alive(entity) {
		panic("can't check components of a dead entity")
	}
	index := m.world.entities[entity.id]
	return m.storageA.columns[index.table] != nil &&
		m.storageB.columns[index.table] != nil
}

// Add the mapped components to the given entity.
func (m *Map2[A, B]) Add(entity Entity, a *A, b *B, rel ...RelationIndex) {
	if !m.world.Alive(entity) {
		panic("can't add components to a dead entity")
	}
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	m.world.exchange(entity, m.ids, nil, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
	}, m.relations)
}

// Remove the mapped components from the given entity.
func (m *Map2[A, B]) Remove(entity Entity) {
	if !m.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	m.world.exchange(entity, nil, m.ids, nil, nil)
}

// Remove the mapped components from the given entity.
func (m *Map2[A, B]) SetRelations(entity Entity, rel ...RelationIndex) {
	if !m.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	m.world.setRelations(entity, m.relations)
}

// Map3 is a mapper to access 3 components of an entity.
type Map3[A any, B any, C any] struct {
	world     *World
	ids       []ID
	storageA  *componentStorage
	storageB  *componentStorage
	storageC  *componentStorage
	relations []relationID
}

// NewMap3 creates a new [Map3].
func NewMap3[A any, B any, C any](world *World) Map3[A, B, C] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
	}
	return Map3[A, B, C]{
		world:    world,
		ids:      ids,
		storageA: &world.storage.components[ids[0].id],
		storageB: &world.storage.components[ids[1].id],
		storageC: &world.storage.components[ids[2].id],
	}
}

// NewEntity creates a new entity with the mapped components.
func (m *Map3[A, B, C]) NewEntity(a *A, b *B, c *C, rel ...RelationIndex) Entity {
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	return m.world.newEntityWith(m.ids, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
		unsafe.Pointer(c),
	}, m.relations)
}

// Get returns the mapped components for the given entity.
func (m *Map3[A, B, C]) Get(entity Entity) (*A, *B, *C) {
	if !m.world.Alive(entity) {
		panic("can't get components of a dead entity")
	}
	return m.GetUnchecked(entity)
}

// GetUnchecked returns the mapped components for the given entity.
// It does not check whether the entity is alive.
// Can be used as an optimization when it is certain that the entity is alive.
func (m *Map3[A, B, C]) GetUnchecked(entity Entity) (*A, *B, *C) {
	index := m.world.entities[entity.id]
	row := uintptr(index.row)
	return (*A)(m.storageA.columns[index.table].Get(row)),
		(*B)(m.storageB.columns[index.table].Get(row)),
		(*C)(m.storageC.columns[index.table].Get(row))
}

// HasAll return whether the given entity has all mapped components.
func (m *Map3[A, B, C]) HasAll(entity Entity) bool {
	if !m.world.Alive(entity) {
		panic("can't check components of a dead entity")
	}
	index := m.world.entities[entity.id]
	return m.storageA.columns[index.table] != nil &&
		m.storageB.columns[index.table] != nil &&
		m.storageC.columns[index.table] != nil
}

// Add the mapped components to the given entity.
func (m *Map3[A, B, C]) Add(entity Entity, a *A, b *B, c *C, rel ...RelationIndex) {
	if !m.world.Alive(entity) {
		panic("can't add components to a dead entity")
	}
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	m.world.exchange(entity, m.ids, nil, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
		unsafe.Pointer(c),
	}, m.relations)
}

// Remove the mapped components from the given entity.
func (m *Map3[A, B, C]) Remove(entity Entity) {
	if !m.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	m.world.exchange(entity, nil, m.ids, nil, nil)
}

// Remove the mapped components from the given entity.
func (m *Map3[A, B, C]) SetRelations(entity Entity, rel ...RelationIndex) {
	if !m.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	m.world.setRelations(entity, m.relations)
}

// Map4 is a mapper to access 4 components of an entity.
type Map4[A any, B any, C any, D any] struct {
	world     *World
	ids       []ID
	storageA  *componentStorage
	storageB  *componentStorage
	storageC  *componentStorage
	storageD  *componentStorage
	relations []relationID
}

// NewMap4 creates a new [Map4].
func NewMap4[A any, B any, C any, D any](world *World) Map4[A, B, C, D] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
		ComponentID[D](world),
	}
	return Map4[A, B, C, D]{
		world:    world,
		ids:      ids,
		storageA: &world.storage.components[ids[0].id],
		storageB: &world.storage.components[ids[1].id],
		storageC: &world.storage.components[ids[2].id],
		storageD: &world.storage.components[ids[3].id],
	}
}

// NewEntity creates a new entity with the mapped components.
func (m *Map4[A, B, C, D]) NewEntity(a *A, b *B, c *C, d *D, rel ...RelationIndex) Entity {
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	return m.world.newEntityWith(m.ids, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
		unsafe.Pointer(c),
		unsafe.Pointer(d),
	}, m.relations)
}

// Get returns the mapped components for the given entity.
func (m *Map4[A, B, C, D]) Get(entity Entity) (*A, *B, *C, *D) {
	if !m.world.Alive(entity) {
		panic("can't get components of a dead entity")
	}
	return m.GetUnchecked(entity)
}

// GetUnchecked returns the mapped components for the given entity.
// It does not check whether the entity is alive.
// Can be used as an optimization when it is certain that the entity is alive.
func (m *Map4[A, B, C, D]) GetUnchecked(entity Entity) (*A, *B, *C, *D) {
	index := m.world.entities[entity.id]
	row := uintptr(index.row)
	return (*A)(m.storageA.columns[index.table].Get(row)),
		(*B)(m.storageB.columns[index.table].Get(row)),
		(*C)(m.storageC.columns[index.table].Get(row)),
		(*D)(m.storageD.columns[index.table].Get(row))
}

// HasAll return whether the given entity has all mapped components.
func (m *Map4[A, B, C, D]) HasAll(entity Entity) bool {
	if !m.world.Alive(entity) {
		panic("can't check components of a dead entity")
	}
	index := m.world.entities[entity.id]
	return m.storageA.columns[index.table] != nil &&
		m.storageB.columns[index.table] != nil &&
		m.storageC.columns[index.table] != nil &&
		m.storageD.columns[index.table] != nil
}

// Add the mapped components to the given entity.
func (m *Map4[A, B, C, D]) Add(entity Entity, a *A, b *B, c *C, d *D, rel ...RelationIndex) {
	if !m.world.Alive(entity) {
		panic("can't add components to a dead entity")
	}
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	m.world.exchange(entity, m.ids, nil, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
		unsafe.Pointer(c),
		unsafe.Pointer(d),
	}, m.relations)
}

// Remove the mapped components from the given entity.
func (m *Map4[A, B, C, D]) Remove(entity Entity) {
	if !m.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	m.world.exchange(entity, nil, m.ids, nil, nil)
}

// Remove the mapped components from the given entity.
func (m *Map4[A, B, C, D]) SetRelations(entity Entity, rel ...RelationIndex) {
	if !m.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	m.world.setRelations(entity, m.relations)
}

// Map5 is a mapper to access 5 components of an entity.
type Map5[A any, B any, C any, D any, E any] struct {
	world     *World
	ids       []ID
	storageA  *componentStorage
	storageB  *componentStorage
	storageC  *componentStorage
	storageD  *componentStorage
	storageE  *componentStorage
	relations []relationID
}

// NewMap5 creates a new [Map5].
func NewMap5[A any, B any, C any, D any, E any](world *World) Map5[A, B, C, D, E] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
		ComponentID[D](world),
		ComponentID[E](world),
	}
	return Map5[A, B, C, D, E]{
		world:    world,
		ids:      ids,
		storageA: &world.storage.components[ids[0].id],
		storageB: &world.storage.components[ids[1].id],
		storageC: &world.storage.components[ids[2].id],
		storageD: &world.storage.components[ids[3].id],
		storageE: &world.storage.components[ids[4].id],
	}
}

// NewEntity creates a new entity with the mapped components.
func (m *Map5[A, B, C, D, E]) NewEntity(a *A, b *B, c *C, d *D, e *E, rel ...RelationIndex) Entity {
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	return m.world.newEntityWith(m.ids, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
		unsafe.Pointer(c),
		unsafe.Pointer(d),
		unsafe.Pointer(e),
	}, m.relations)
}

// Get returns the mapped components for the given entity.
func (m *Map5[A, B, C, D, E]) Get(entity Entity) (*A, *B, *C, *D, *E) {
	if !m.world.Alive(entity) {
		panic("can't get components of a dead entity")
	}
	return m.GetUnchecked(entity)
}

// GetUnchecked returns the mapped components for the given entity.
// It does not check whether the entity is alive.
// Can be used as an optimization when it is certain that the entity is alive.
func (m *Map5[A, B, C, D, E]) GetUnchecked(entity Entity) (*A, *B, *C, *D, *E) {
	index := m.world.entities[entity.id]
	row := uintptr(index.row)
	return (*A)(m.storageA.columns[index.table].Get(row)),
		(*B)(m.storageB.columns[index.table].Get(row)),
		(*C)(m.storageC.columns[index.table].Get(row)),
		(*D)(m.storageD.columns[index.table].Get(row)),
		(*E)(m.storageE.columns[index.table].Get(row))
}

// HasAll return whether the given entity has all mapped components.
func (m *Map5[A, B, C, D, E]) HasAll(entity Entity) bool {
	if !m.world.Alive(entity) {
		panic("can't check components of a dead entity")
	}
	index := m.world.entities[entity.id]
	return m.storageA.columns[index.table] != nil &&
		m.storageB.columns[index.table] != nil &&
		m.storageC.columns[index.table] != nil &&
		m.storageD.columns[index.table] != nil &&
		m.storageE.columns[index.table] != nil
}

// Add the mapped components to the given entity.
func (m *Map5[A, B, C, D, E]) Add(entity Entity, a *A, b *B, c *C, d *D, e *E, rel ...RelationIndex) {
	if !m.world.Alive(entity) {
		panic("can't add components to a dead entity")
	}
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	m.world.exchange(entity, m.ids, nil, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
		unsafe.Pointer(c),
		unsafe.Pointer(d),
		unsafe.Pointer(e),
	}, m.relations)
}

// Remove the mapped components from the given entity.
func (m *Map5[A, B, C, D, E]) Remove(entity Entity) {
	if !m.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	m.world.exchange(entity, nil, m.ids, nil, nil)
}

// Remove the mapped components from the given entity.
func (m *Map5[A, B, C, D, E]) SetRelations(entity Entity, rel ...RelationIndex) {
	if !m.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	m.world.setRelations(entity, m.relations)
}

// Map6 is a mapper to access 6 components of an entity.
type Map6[A any, B any, C any, D any, E any, F any] struct {
	world     *World
	ids       []ID
	storageA  *componentStorage
	storageB  *componentStorage
	storageC  *componentStorage
	storageD  *componentStorage
	storageE  *componentStorage
	storageF  *componentStorage
	relations []relationID
}

// NewMap6 creates a new [Map6].
func NewMap6[A any, B any, C any, D any, E any, F any](world *World) Map6[A, B, C, D, E, F] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
		ComponentID[D](world),
		ComponentID[E](world),
		ComponentID[F](world),
	}
	return Map6[A, B, C, D, E, F]{
		world:    world,
		ids:      ids,
		storageA: &world.storage.components[ids[0].id],
		storageB: &world.storage.components[ids[1].id],
		storageC: &world.storage.components[ids[2].id],
		storageD: &world.storage.components[ids[3].id],
		storageE: &world.storage.components[ids[4].id],
		storageF: &world.storage.components[ids[5].id],
	}
}

// NewEntity creates a new entity with the mapped components.
func (m *Map6[A, B, C, D, E, F]) NewEntity(a *A, b *B, c *C, d *D, e *E, f *F, rel ...RelationIndex) Entity {
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	return m.world.newEntityWith(m.ids, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
		unsafe.Pointer(c),
		unsafe.Pointer(d),
		unsafe.Pointer(e),
		unsafe.Pointer(f),
	}, m.relations)
}

// Get returns the mapped components for the given entity.
func (m *Map6[A, B, C, D, E, F]) Get(entity Entity) (*A, *B, *C, *D, *E, *F) {
	if !m.world.Alive(entity) {
		panic("can't get components of a dead entity")
	}
	return m.GetUnchecked(entity)
}

// GetUnchecked returns the mapped components for the given entity.
// It does not check whether the entity is alive.
// Can be used as an optimization when it is certain that the entity is alive.
func (m *Map6[A, B, C, D, E, F]) GetUnchecked(entity Entity) (*A, *B, *C, *D, *E, *F) {
	index := m.world.entities[entity.id]
	row := uintptr(index.row)
	return (*A)(m.storageA.columns[index.table].Get(row)),
		(*B)(m.storageB.columns[index.table].Get(row)),
		(*C)(m.storageC.columns[index.table].Get(row)),
		(*D)(m.storageD.columns[index.table].Get(row)),
		(*E)(m.storageE.columns[index.table].Get(row)),
		(*F)(m.storageF.columns[index.table].Get(row))
}

// HasAll return whether the given entity has all mapped components.
func (m *Map6[A, B, C, D, E, F]) HasAll(entity Entity) bool {
	if !m.world.Alive(entity) {
		panic("can't check components of a dead entity")
	}
	index := m.world.entities[entity.id]
	return m.storageA.columns[index.table] != nil &&
		m.storageB.columns[index.table] != nil &&
		m.storageC.columns[index.table] != nil &&
		m.storageD.columns[index.table] != nil &&
		m.storageE.columns[index.table] != nil &&
		m.storageF.columns[index.table] != nil
}

// Add the mapped components to the given entity.
func (m *Map6[A, B, C, D, E, F]) Add(entity Entity, a *A, b *B, c *C, d *D, e *E, f *F, rel ...RelationIndex) {
	if !m.world.Alive(entity) {
		panic("can't add components to a dead entity")
	}
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	m.world.exchange(entity, m.ids, nil, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
		unsafe.Pointer(c),
		unsafe.Pointer(d),
		unsafe.Pointer(e),
		unsafe.Pointer(f),
	}, m.relations)
}

// Remove the mapped components from the given entity.
func (m *Map6[A, B, C, D, E, F]) Remove(entity Entity) {
	if !m.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	m.world.exchange(entity, nil, m.ids, nil, nil)
}

// Remove the mapped components from the given entity.
func (m *Map6[A, B, C, D, E, F]) SetRelations(entity Entity, rel ...RelationIndex) {
	if !m.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	m.world.setRelations(entity, m.relations)
}

// Map7 is a mapper to access 7 components of an entity.
type Map7[A any, B any, C any, D any, E any, F any, G any] struct {
	world     *World
	ids       []ID
	storageA  *componentStorage
	storageB  *componentStorage
	storageC  *componentStorage
	storageD  *componentStorage
	storageE  *componentStorage
	storageF  *componentStorage
	storageG  *componentStorage
	relations []relationID
}

// NewMap7 creates a new [Map7].
func NewMap7[A any, B any, C any, D any, E any, F any, G any](world *World) Map7[A, B, C, D, E, F, G] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
		ComponentID[D](world),
		ComponentID[E](world),
		ComponentID[F](world),
		ComponentID[G](world),
	}
	return Map7[A, B, C, D, E, F, G]{
		world:    world,
		ids:      ids,
		storageA: &world.storage.components[ids[0].id],
		storageB: &world.storage.components[ids[1].id],
		storageC: &world.storage.components[ids[2].id],
		storageD: &world.storage.components[ids[3].id],
		storageE: &world.storage.components[ids[4].id],
		storageF: &world.storage.components[ids[5].id],
		storageG: &world.storage.components[ids[6].id],
	}
}

// NewEntity creates a new entity with the mapped components.
func (m *Map7[A, B, C, D, E, F, G]) NewEntity(a *A, b *B, c *C, d *D, e *E, f *F, g *G, rel ...RelationIndex) Entity {
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	return m.world.newEntityWith(m.ids, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
		unsafe.Pointer(c),
		unsafe.Pointer(d),
		unsafe.Pointer(e),
		unsafe.Pointer(f),
		unsafe.Pointer(g),
	}, m.relations)
}

// Get returns the mapped components for the given entity.
func (m *Map7[A, B, C, D, E, F, G]) Get(entity Entity) (*A, *B, *C, *D, *E, *F, *G) {
	if !m.world.Alive(entity) {
		panic("can't get components of a dead entity")
	}
	return m.GetUnchecked(entity)
}

// GetUnchecked returns the mapped components for the given entity.
// It does not check whether the entity is alive.
// Can be used as an optimization when it is certain that the entity is alive.
func (m *Map7[A, B, C, D, E, F, G]) GetUnchecked(entity Entity) (*A, *B, *C, *D, *E, *F, *G) {
	index := m.world.entities[entity.id]
	row := uintptr(index.row)
	return (*A)(m.storageA.columns[index.table].Get(row)),
		(*B)(m.storageB.columns[index.table].Get(row)),
		(*C)(m.storageC.columns[index.table].Get(row)),
		(*D)(m.storageD.columns[index.table].Get(row)),
		(*E)(m.storageE.columns[index.table].Get(row)),
		(*F)(m.storageF.columns[index.table].Get(row)),
		(*G)(m.storageG.columns[index.table].Get(row))
}

// HasAll return whether the given entity has all mapped components.
func (m *Map7[A, B, C, D, E, F, G]) HasAll(entity Entity) bool {
	if !m.world.Alive(entity) {
		panic("can't check components of a dead entity")
	}
	index := m.world.entities[entity.id]
	return m.storageA.columns[index.table] != nil &&
		m.storageB.columns[index.table] != nil &&
		m.storageC.columns[index.table] != nil &&
		m.storageD.columns[index.table] != nil &&
		m.storageE.columns[index.table] != nil &&
		m.storageF.columns[index.table] != nil &&
		m.storageG.columns[index.table] != nil
}

// Add the mapped components to the given entity.
func (m *Map7[A, B, C, D, E, F, G]) Add(entity Entity, a *A, b *B, c *C, d *D, e *E, f *F, g *G, rel ...RelationIndex) {
	if !m.world.Alive(entity) {
		panic("can't add components to a dead entity")
	}
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	m.world.exchange(entity, m.ids, nil, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
		unsafe.Pointer(c),
		unsafe.Pointer(d),
		unsafe.Pointer(e),
		unsafe.Pointer(f),
		unsafe.Pointer(g),
	}, m.relations)
}

// Remove the mapped components from the given entity.
func (m *Map7[A, B, C, D, E, F, G]) Remove(entity Entity) {
	if !m.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	m.world.exchange(entity, nil, m.ids, nil, nil)
}

// Remove the mapped components from the given entity.
func (m *Map7[A, B, C, D, E, F, G]) SetRelations(entity Entity, rel ...RelationIndex) {
	if !m.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	m.world.setRelations(entity, m.relations)
}

// Map8 is a mapper to access 8 components of an entity.
type Map8[A any, B any, C any, D any, E any, F any, G any, H any] struct {
	world     *World
	ids       []ID
	storageA  *componentStorage
	storageB  *componentStorage
	storageC  *componentStorage
	storageD  *componentStorage
	storageE  *componentStorage
	storageF  *componentStorage
	storageG  *componentStorage
	storageH  *componentStorage
	relations []relationID
}

// NewMap8 creates a new [Map8].
func NewMap8[A any, B any, C any, D any, E any, F any, G any, H any](world *World) Map8[A, B, C, D, E, F, G, H] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
		ComponentID[D](world),
		ComponentID[E](world),
		ComponentID[F](world),
		ComponentID[G](world),
		ComponentID[H](world),
	}
	return Map8[A, B, C, D, E, F, G, H]{
		world:    world,
		ids:      ids,
		storageA: &world.storage.components[ids[0].id],
		storageB: &world.storage.components[ids[1].id],
		storageC: &world.storage.components[ids[2].id],
		storageD: &world.storage.components[ids[3].id],
		storageE: &world.storage.components[ids[4].id],
		storageF: &world.storage.components[ids[5].id],
		storageG: &world.storage.components[ids[6].id],
		storageH: &world.storage.components[ids[7].id],
	}
}

// NewEntity creates a new entity with the mapped components.
func (m *Map8[A, B, C, D, E, F, G, H]) NewEntity(a *A, b *B, c *C, d *D, e *E, f *F, g *G, h *H, rel ...RelationIndex) Entity {
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	return m.world.newEntityWith(m.ids, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
		unsafe.Pointer(c),
		unsafe.Pointer(d),
		unsafe.Pointer(e),
		unsafe.Pointer(f),
		unsafe.Pointer(g),
		unsafe.Pointer(h),
	}, m.relations)
}

// Get returns the mapped components for the given entity.
func (m *Map8[A, B, C, D, E, F, G, H]) Get(entity Entity) (*A, *B, *C, *D, *E, *F, *G, *H) {
	if !m.world.Alive(entity) {
		panic("can't get components of a dead entity")
	}
	return m.GetUnchecked(entity)
}

// GetUnchecked returns the mapped components for the given entity.
// It does not check whether the entity is alive.
// Can be used as an optimization when it is certain that the entity is alive.
func (m *Map8[A, B, C, D, E, F, G, H]) GetUnchecked(entity Entity) (*A, *B, *C, *D, *E, *F, *G, *H) {
	index := m.world.entities[entity.id]
	row := uintptr(index.row)
	return (*A)(m.storageA.columns[index.table].Get(row)),
		(*B)(m.storageB.columns[index.table].Get(row)),
		(*C)(m.storageC.columns[index.table].Get(row)),
		(*D)(m.storageD.columns[index.table].Get(row)),
		(*E)(m.storageE.columns[index.table].Get(row)),
		(*F)(m.storageF.columns[index.table].Get(row)),
		(*G)(m.storageG.columns[index.table].Get(row)),
		(*H)(m.storageH.columns[index.table].Get(row))
}

// HasAll return whether the given entity has all mapped components.
func (m *Map8[A, B, C, D, E, F, G, H]) HasAll(entity Entity) bool {
	if !m.world.Alive(entity) {
		panic("can't check components of a dead entity")
	}
	index := m.world.entities[entity.id]
	return m.storageA.columns[index.table] != nil &&
		m.storageB.columns[index.table] != nil &&
		m.storageC.columns[index.table] != nil &&
		m.storageD.columns[index.table] != nil &&
		m.storageE.columns[index.table] != nil &&
		m.storageF.columns[index.table] != nil &&
		m.storageG.columns[index.table] != nil &&
		m.storageH.columns[index.table] != nil
}

// Add the mapped components to the given entity.
func (m *Map8[A, B, C, D, E, F, G, H]) Add(entity Entity, a *A, b *B, c *C, d *D, e *E, f *F, g *G, h *H, rel ...RelationIndex) {
	if !m.world.Alive(entity) {
		panic("can't add components to a dead entity")
	}
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	m.world.exchange(entity, m.ids, nil, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
		unsafe.Pointer(c),
		unsafe.Pointer(d),
		unsafe.Pointer(e),
		unsafe.Pointer(f),
		unsafe.Pointer(g),
		unsafe.Pointer(h),
	}, m.relations)
}

// Remove the mapped components from the given entity.
func (m *Map8[A, B, C, D, E, F, G, H]) Remove(entity Entity) {
	if !m.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	m.world.exchange(entity, nil, m.ids, nil, nil)
}

// Remove the mapped components from the given entity.
func (m *Map8[A, B, C, D, E, F, G, H]) SetRelations(entity Entity, rel ...RelationIndex) {
	if !m.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	m.relations = relations(rel).toRelations(m.ids, m.relations)
	m.world.setRelations(entity, m.relations)
}
