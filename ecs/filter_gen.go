package ecs

// Code generated by go generate; DO NOT EDIT.

// Filter0 is a filter for 0 components.
type Filter0 struct {
	world         *World
	ids           []ID
	filter        Filter
	relations     []RelationID
	tempRelations []RelationID
	cache         *cacheEntry
}

// NewFilter0 creates a new [Filter0].
//
// Use [Filter0.Query] to obtain a [Query0].
func NewFilter0(world *World) *Filter0 {
	ids := []ID{}

	return &Filter0{
		world:  world,
		ids:    ids,
		filter: NewFilter(ids...),
	}
}

// With specifies additional components to filter for.
func (f *Filter0) With(comps ...Comp) *Filter0 {
	for _, c := range comps {
		id := f.world.componentID(c.tp)
		f.filter.mask.Set(id, true)
	}
	return f
}

// Without specifies components to exclude.
func (f *Filter0) Without(comps ...Comp) *Filter0 {
	for _, c := range comps {
		id := f.world.componentID(c.tp)
		f.filter.without.Set(id, true)
		f.filter.hasWithout = true
	}
	return f
}

// Exclusive makes the filter exclusive in the sense that the component composition is matched exactly,
// and no other components are allowed.
func (f *Filter0) Exclusive() *Filter0 {
	f.filter = f.filter.Exclusive()
	return f
}

// Relations sets permanent entity relation targets for this filter.
// Relation targets set here are included in filter caching.
// Contrary, relation targets specified in [Filter0.Query] or [Filter0.Batch] are not cached.
func (f *Filter0) Relations(rel ...RelationIndex) *Filter0 {
	f.relations = relations(rel).toRelations(&f.world.storage.registry, f.ids, nil, f.relations)
	return f
}

// Register this filter to the world's filter cache.
func (f *Filter0) Register() *Filter0 {
	if f.cache != nil {
		panic("filter is already registered, can't register")
	}
	f.cache = f.world.storage.registerFilter(f.Batch())
	return f
}

// Unregister this filter from the world's filter cache.
func (f *Filter0) Unregister() {
	if f.cache == nil {
		panic("filter is not registered, can't unregister")
	}
	f.world.storage.unregisterFilter(f.cache)
	f.cache = nil
}

// Query creates a [Query0] from this filter.
// This must be used each time before iterating a query.
func (f *Filter0) Query(rel ...RelationIndex) Query0 {
	if f.cache == nil {
		f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, f.relations, f.tempRelations)
	} else {
		f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, nil, f.tempRelations)
	}
	return newQuery0(f.world, f.filter, f.tempRelations, f.cache)
}

// Batch creates a [Batch] from this filter.
func (f *Filter0) Batch(rel ...RelationIndex) *Batch {
	// TODO: use cache?
	f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, f.relations, f.tempRelations)
	return &Batch{
		filter:    f.filter,
		relations: f.tempRelations,
	}
}

// Filter1 is a filter for 1 components.
type Filter1[A any] struct {
	world         *World
	ids           []ID
	filter        Filter
	relations     []RelationID
	tempRelations []RelationID
	cache         *cacheEntry
}

// NewFilter1 creates a new [Filter1].
//
// Use [Filter1.Query] to obtain a [Query1].
func NewFilter1[A any](world *World) *Filter1[A] {
	ids := []ID{
		ComponentID[A](world),
	}

	return &Filter1[A]{
		world:  world,
		ids:    ids,
		filter: NewFilter(ids...),
	}
}

// With specifies additional components to filter for.
func (f *Filter1[A]) With(comps ...Comp) *Filter1[A] {
	for _, c := range comps {
		id := f.world.componentID(c.tp)
		f.filter.mask.Set(id, true)
	}
	return f
}

// Without specifies components to exclude.
func (f *Filter1[A]) Without(comps ...Comp) *Filter1[A] {
	for _, c := range comps {
		id := f.world.componentID(c.tp)
		f.filter.without.Set(id, true)
		f.filter.hasWithout = true
	}
	return f
}

// Exclusive makes the filter exclusive in the sense that the component composition is matched exactly,
// and no other components are allowed.
func (f *Filter1[A]) Exclusive() *Filter1[A] {
	f.filter = f.filter.Exclusive()
	return f
}

// Relations sets permanent entity relation targets for this filter.
// Relation targets set here are included in filter caching.
// Contrary, relation targets specified in [Filter1.Query] or [Filter1.Batch] are not cached.
func (f *Filter1[A]) Relations(rel ...RelationIndex) *Filter1[A] {
	f.relations = relations(rel).toRelations(&f.world.storage.registry, f.ids, nil, f.relations)
	return f
}

// Register this filter to the world's filter cache.
func (f *Filter1[A]) Register() *Filter1[A] {
	if f.cache != nil {
		panic("filter is already registered, can't register")
	}
	f.cache = f.world.storage.registerFilter(f.Batch())
	return f
}

// Unregister this filter from the world's filter cache.
func (f *Filter1[A]) Unregister() {
	if f.cache == nil {
		panic("filter is not registered, can't unregister")
	}
	f.world.storage.unregisterFilter(f.cache)
	f.cache = nil
}

// Query creates a [Query1] from this filter.
// This must be used each time before iterating a query.
func (f *Filter1[A]) Query(rel ...RelationIndex) Query1[A] {
	if f.cache == nil {
		f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, f.relations, f.tempRelations)
	} else {
		f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, nil, f.tempRelations)
	}
	return newQuery1[A](f.world, f.filter, f.ids, f.tempRelations, f.cache)
}

// Batch creates a [Batch] from this filter.
func (f *Filter1[A]) Batch(rel ...RelationIndex) *Batch {
	// TODO: use cache?
	f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, f.relations, f.tempRelations)
	return &Batch{
		filter:    f.filter,
		relations: f.tempRelations,
	}
}

// Filter2 is a filter for 2 components.
type Filter2[A any, B any] struct {
	world         *World
	ids           []ID
	filter        Filter
	relations     []RelationID
	tempRelations []RelationID
	cache         *cacheEntry
}

// NewFilter2 creates a new [Filter2].
//
// Use [Filter2.Query] to obtain a [Query2].
func NewFilter2[A any, B any](world *World) *Filter2[A, B] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
	}

	return &Filter2[A, B]{
		world:  world,
		ids:    ids,
		filter: NewFilter(ids...),
	}
}

// With specifies additional components to filter for.
func (f *Filter2[A, B]) With(comps ...Comp) *Filter2[A, B] {
	for _, c := range comps {
		id := f.world.componentID(c.tp)
		f.filter.mask.Set(id, true)
	}
	return f
}

// Without specifies components to exclude.
func (f *Filter2[A, B]) Without(comps ...Comp) *Filter2[A, B] {
	for _, c := range comps {
		id := f.world.componentID(c.tp)
		f.filter.without.Set(id, true)
		f.filter.hasWithout = true
	}
	return f
}

// Exclusive makes the filter exclusive in the sense that the component composition is matched exactly,
// and no other components are allowed.
func (f *Filter2[A, B]) Exclusive() *Filter2[A, B] {
	f.filter = f.filter.Exclusive()
	return f
}

// Relations sets permanent entity relation targets for this filter.
// Relation targets set here are included in filter caching.
// Contrary, relation targets specified in [Filter2.Query] or [Filter2.Batch] are not cached.
func (f *Filter2[A, B]) Relations(rel ...RelationIndex) *Filter2[A, B] {
	f.relations = relations(rel).toRelations(&f.world.storage.registry, f.ids, nil, f.relations)
	return f
}

// Register this filter to the world's filter cache.
func (f *Filter2[A, B]) Register() *Filter2[A, B] {
	if f.cache != nil {
		panic("filter is already registered, can't register")
	}
	f.cache = f.world.storage.registerFilter(f.Batch())
	return f
}

// Unregister this filter from the world's filter cache.
func (f *Filter2[A, B]) Unregister() {
	if f.cache == nil {
		panic("filter is not registered, can't unregister")
	}
	f.world.storage.unregisterFilter(f.cache)
	f.cache = nil
}

// Query creates a [Query2] from this filter.
// This must be used each time before iterating a query.
func (f *Filter2[A, B]) Query(rel ...RelationIndex) Query2[A, B] {
	if f.cache == nil {
		f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, f.relations, f.tempRelations)
	} else {
		f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, nil, f.tempRelations)
	}
	return newQuery2[A, B](f.world, f.filter, f.ids, f.tempRelations, f.cache)
}

// Batch creates a [Batch] from this filter.
func (f *Filter2[A, B]) Batch(rel ...RelationIndex) *Batch {
	// TODO: use cache?
	f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, f.relations, f.tempRelations)
	return &Batch{
		filter:    f.filter,
		relations: f.tempRelations,
	}
}

// Filter3 is a filter for 3 components.
type Filter3[A any, B any, C any] struct {
	world         *World
	ids           []ID
	filter        Filter
	relations     []RelationID
	tempRelations []RelationID
	cache         *cacheEntry
}

// NewFilter3 creates a new [Filter3].
//
// Use [Filter3.Query] to obtain a [Query3].
func NewFilter3[A any, B any, C any](world *World) *Filter3[A, B, C] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
	}

	return &Filter3[A, B, C]{
		world:  world,
		ids:    ids,
		filter: NewFilter(ids...),
	}
}

// With specifies additional components to filter for.
func (f *Filter3[A, B, C]) With(comps ...Comp) *Filter3[A, B, C] {
	for _, c := range comps {
		id := f.world.componentID(c.tp)
		f.filter.mask.Set(id, true)
	}
	return f
}

// Without specifies components to exclude.
func (f *Filter3[A, B, C]) Without(comps ...Comp) *Filter3[A, B, C] {
	for _, c := range comps {
		id := f.world.componentID(c.tp)
		f.filter.without.Set(id, true)
		f.filter.hasWithout = true
	}
	return f
}

// Exclusive makes the filter exclusive in the sense that the component composition is matched exactly,
// and no other components are allowed.
func (f *Filter3[A, B, C]) Exclusive() *Filter3[A, B, C] {
	f.filter = f.filter.Exclusive()
	return f
}

// Relations sets permanent entity relation targets for this filter.
// Relation targets set here are included in filter caching.
// Contrary, relation targets specified in [Filter3.Query] or [Filter3.Batch] are not cached.
func (f *Filter3[A, B, C]) Relations(rel ...RelationIndex) *Filter3[A, B, C] {
	f.relations = relations(rel).toRelations(&f.world.storage.registry, f.ids, nil, f.relations)
	return f
}

// Register this filter to the world's filter cache.
func (f *Filter3[A, B, C]) Register() *Filter3[A, B, C] {
	if f.cache != nil {
		panic("filter is already registered, can't register")
	}
	f.cache = f.world.storage.registerFilter(f.Batch())
	return f
}

// Unregister this filter from the world's filter cache.
func (f *Filter3[A, B, C]) Unregister() {
	if f.cache == nil {
		panic("filter is not registered, can't unregister")
	}
	f.world.storage.unregisterFilter(f.cache)
	f.cache = nil
}

// Query creates a [Query3] from this filter.
// This must be used each time before iterating a query.
func (f *Filter3[A, B, C]) Query(rel ...RelationIndex) Query3[A, B, C] {
	if f.cache == nil {
		f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, f.relations, f.tempRelations)
	} else {
		f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, nil, f.tempRelations)
	}
	return newQuery3[A, B, C](f.world, f.filter, f.ids, f.tempRelations, f.cache)
}

// Batch creates a [Batch] from this filter.
func (f *Filter3[A, B, C]) Batch(rel ...RelationIndex) *Batch {
	// TODO: use cache?
	f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, f.relations, f.tempRelations)
	return &Batch{
		filter:    f.filter,
		relations: f.tempRelations,
	}
}

// Filter4 is a filter for 4 components.
type Filter4[A any, B any, C any, D any] struct {
	world         *World
	ids           []ID
	filter        Filter
	relations     []RelationID
	tempRelations []RelationID
	cache         *cacheEntry
}

// NewFilter4 creates a new [Filter4].
//
// Use [Filter4.Query] to obtain a [Query4].
func NewFilter4[A any, B any, C any, D any](world *World) *Filter4[A, B, C, D] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
		ComponentID[D](world),
	}

	return &Filter4[A, B, C, D]{
		world:  world,
		ids:    ids,
		filter: NewFilter(ids...),
	}
}

// With specifies additional components to filter for.
func (f *Filter4[A, B, C, D]) With(comps ...Comp) *Filter4[A, B, C, D] {
	for _, c := range comps {
		id := f.world.componentID(c.tp)
		f.filter.mask.Set(id, true)
	}
	return f
}

// Without specifies components to exclude.
func (f *Filter4[A, B, C, D]) Without(comps ...Comp) *Filter4[A, B, C, D] {
	for _, c := range comps {
		id := f.world.componentID(c.tp)
		f.filter.without.Set(id, true)
		f.filter.hasWithout = true
	}
	return f
}

// Exclusive makes the filter exclusive in the sense that the component composition is matched exactly,
// and no other components are allowed.
func (f *Filter4[A, B, C, D]) Exclusive() *Filter4[A, B, C, D] {
	f.filter = f.filter.Exclusive()
	return f
}

// Relations sets permanent entity relation targets for this filter.
// Relation targets set here are included in filter caching.
// Contrary, relation targets specified in [Filter4.Query] or [Filter4.Batch] are not cached.
func (f *Filter4[A, B, C, D]) Relations(rel ...RelationIndex) *Filter4[A, B, C, D] {
	f.relations = relations(rel).toRelations(&f.world.storage.registry, f.ids, nil, f.relations)
	return f
}

// Register this filter to the world's filter cache.
func (f *Filter4[A, B, C, D]) Register() *Filter4[A, B, C, D] {
	if f.cache != nil {
		panic("filter is already registered, can't register")
	}
	f.cache = f.world.storage.registerFilter(f.Batch())
	return f
}

// Unregister this filter from the world's filter cache.
func (f *Filter4[A, B, C, D]) Unregister() {
	if f.cache == nil {
		panic("filter is not registered, can't unregister")
	}
	f.world.storage.unregisterFilter(f.cache)
	f.cache = nil
}

// Query creates a [Query4] from this filter.
// This must be used each time before iterating a query.
func (f *Filter4[A, B, C, D]) Query(rel ...RelationIndex) Query4[A, B, C, D] {
	if f.cache == nil {
		f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, f.relations, f.tempRelations)
	} else {
		f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, nil, f.tempRelations)
	}
	return newQuery4[A, B, C, D](f.world, f.filter, f.ids, f.tempRelations, f.cache)
}

// Batch creates a [Batch] from this filter.
func (f *Filter4[A, B, C, D]) Batch(rel ...RelationIndex) *Batch {
	// TODO: use cache?
	f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, f.relations, f.tempRelations)
	return &Batch{
		filter:    f.filter,
		relations: f.tempRelations,
	}
}

// Filter5 is a filter for 5 components.
type Filter5[A any, B any, C any, D any, E any] struct {
	world         *World
	ids           []ID
	filter        Filter
	relations     []RelationID
	tempRelations []RelationID
	cache         *cacheEntry
}

// NewFilter5 creates a new [Filter5].
//
// Use [Filter5.Query] to obtain a [Query5].
func NewFilter5[A any, B any, C any, D any, E any](world *World) *Filter5[A, B, C, D, E] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
		ComponentID[D](world),
		ComponentID[E](world),
	}

	return &Filter5[A, B, C, D, E]{
		world:  world,
		ids:    ids,
		filter: NewFilter(ids...),
	}
}

// With specifies additional components to filter for.
func (f *Filter5[A, B, C, D, E]) With(comps ...Comp) *Filter5[A, B, C, D, E] {
	for _, c := range comps {
		id := f.world.componentID(c.tp)
		f.filter.mask.Set(id, true)
	}
	return f
}

// Without specifies components to exclude.
func (f *Filter5[A, B, C, D, E]) Without(comps ...Comp) *Filter5[A, B, C, D, E] {
	for _, c := range comps {
		id := f.world.componentID(c.tp)
		f.filter.without.Set(id, true)
		f.filter.hasWithout = true
	}
	return f
}

// Exclusive makes the filter exclusive in the sense that the component composition is matched exactly,
// and no other components are allowed.
func (f *Filter5[A, B, C, D, E]) Exclusive() *Filter5[A, B, C, D, E] {
	f.filter = f.filter.Exclusive()
	return f
}

// Relations sets permanent entity relation targets for this filter.
// Relation targets set here are included in filter caching.
// Contrary, relation targets specified in [Filter5.Query] or [Filter5.Batch] are not cached.
func (f *Filter5[A, B, C, D, E]) Relations(rel ...RelationIndex) *Filter5[A, B, C, D, E] {
	f.relations = relations(rel).toRelations(&f.world.storage.registry, f.ids, nil, f.relations)
	return f
}

// Register this filter to the world's filter cache.
func (f *Filter5[A, B, C, D, E]) Register() *Filter5[A, B, C, D, E] {
	if f.cache != nil {
		panic("filter is already registered, can't register")
	}
	f.cache = f.world.storage.registerFilter(f.Batch())
	return f
}

// Unregister this filter from the world's filter cache.
func (f *Filter5[A, B, C, D, E]) Unregister() {
	if f.cache == nil {
		panic("filter is not registered, can't unregister")
	}
	f.world.storage.unregisterFilter(f.cache)
	f.cache = nil
}

// Query creates a [Query5] from this filter.
// This must be used each time before iterating a query.
func (f *Filter5[A, B, C, D, E]) Query(rel ...RelationIndex) Query5[A, B, C, D, E] {
	if f.cache == nil {
		f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, f.relations, f.tempRelations)
	} else {
		f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, nil, f.tempRelations)
	}
	return newQuery5[A, B, C, D, E](f.world, f.filter, f.ids, f.tempRelations, f.cache)
}

// Batch creates a [Batch] from this filter.
func (f *Filter5[A, B, C, D, E]) Batch(rel ...RelationIndex) *Batch {
	// TODO: use cache?
	f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, f.relations, f.tempRelations)
	return &Batch{
		filter:    f.filter,
		relations: f.tempRelations,
	}
}

// Filter6 is a filter for 6 components.
type Filter6[A any, B any, C any, D any, E any, F any] struct {
	world         *World
	ids           []ID
	filter        Filter
	relations     []RelationID
	tempRelations []RelationID
	cache         *cacheEntry
}

// NewFilter6 creates a new [Filter6].
//
// Use [Filter6.Query] to obtain a [Query6].
func NewFilter6[A any, B any, C any, D any, E any, F any](world *World) *Filter6[A, B, C, D, E, F] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
		ComponentID[D](world),
		ComponentID[E](world),
		ComponentID[F](world),
	}

	return &Filter6[A, B, C, D, E, F]{
		world:  world,
		ids:    ids,
		filter: NewFilter(ids...),
	}
}

// With specifies additional components to filter for.
func (f *Filter6[A, B, C, D, E, F]) With(comps ...Comp) *Filter6[A, B, C, D, E, F] {
	for _, c := range comps {
		id := f.world.componentID(c.tp)
		f.filter.mask.Set(id, true)
	}
	return f
}

// Without specifies components to exclude.
func (f *Filter6[A, B, C, D, E, F]) Without(comps ...Comp) *Filter6[A, B, C, D, E, F] {
	for _, c := range comps {
		id := f.world.componentID(c.tp)
		f.filter.without.Set(id, true)
		f.filter.hasWithout = true
	}
	return f
}

// Exclusive makes the filter exclusive in the sense that the component composition is matched exactly,
// and no other components are allowed.
func (f *Filter6[A, B, C, D, E, F]) Exclusive() *Filter6[A, B, C, D, E, F] {
	f.filter = f.filter.Exclusive()
	return f
}

// Relations sets permanent entity relation targets for this filter.
// Relation targets set here are included in filter caching.
// Contrary, relation targets specified in [Filter6.Query] or [Filter6.Batch] are not cached.
func (f *Filter6[A, B, C, D, E, F]) Relations(rel ...RelationIndex) *Filter6[A, B, C, D, E, F] {
	f.relations = relations(rel).toRelations(&f.world.storage.registry, f.ids, nil, f.relations)
	return f
}

// Register this filter to the world's filter cache.
func (f *Filter6[A, B, C, D, E, F]) Register() *Filter6[A, B, C, D, E, F] {
	if f.cache != nil {
		panic("filter is already registered, can't register")
	}
	f.cache = f.world.storage.registerFilter(f.Batch())
	return f
}

// Unregister this filter from the world's filter cache.
func (f *Filter6[A, B, C, D, E, F]) Unregister() {
	if f.cache == nil {
		panic("filter is not registered, can't unregister")
	}
	f.world.storage.unregisterFilter(f.cache)
	f.cache = nil
}

// Query creates a [Query6] from this filter.
// This must be used each time before iterating a query.
func (f *Filter6[A, B, C, D, E, F]) Query(rel ...RelationIndex) Query6[A, B, C, D, E, F] {
	if f.cache == nil {
		f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, f.relations, f.tempRelations)
	} else {
		f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, nil, f.tempRelations)
	}
	return newQuery6[A, B, C, D, E, F](f.world, f.filter, f.ids, f.tempRelations, f.cache)
}

// Batch creates a [Batch] from this filter.
func (f *Filter6[A, B, C, D, E, F]) Batch(rel ...RelationIndex) *Batch {
	// TODO: use cache?
	f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, f.relations, f.tempRelations)
	return &Batch{
		filter:    f.filter,
		relations: f.tempRelations,
	}
}

// Filter7 is a filter for 7 components.
type Filter7[A any, B any, C any, D any, E any, F any, G any] struct {
	world         *World
	ids           []ID
	filter        Filter
	relations     []RelationID
	tempRelations []RelationID
	cache         *cacheEntry
}

// NewFilter7 creates a new [Filter7].
//
// Use [Filter7.Query] to obtain a [Query7].
func NewFilter7[A any, B any, C any, D any, E any, F any, G any](world *World) *Filter7[A, B, C, D, E, F, G] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
		ComponentID[D](world),
		ComponentID[E](world),
		ComponentID[F](world),
		ComponentID[G](world),
	}

	return &Filter7[A, B, C, D, E, F, G]{
		world:  world,
		ids:    ids,
		filter: NewFilter(ids...),
	}
}

// With specifies additional components to filter for.
func (f *Filter7[A, B, C, D, E, F, G]) With(comps ...Comp) *Filter7[A, B, C, D, E, F, G] {
	for _, c := range comps {
		id := f.world.componentID(c.tp)
		f.filter.mask.Set(id, true)
	}
	return f
}

// Without specifies components to exclude.
func (f *Filter7[A, B, C, D, E, F, G]) Without(comps ...Comp) *Filter7[A, B, C, D, E, F, G] {
	for _, c := range comps {
		id := f.world.componentID(c.tp)
		f.filter.without.Set(id, true)
		f.filter.hasWithout = true
	}
	return f
}

// Exclusive makes the filter exclusive in the sense that the component composition is matched exactly,
// and no other components are allowed.
func (f *Filter7[A, B, C, D, E, F, G]) Exclusive() *Filter7[A, B, C, D, E, F, G] {
	f.filter = f.filter.Exclusive()
	return f
}

// Relations sets permanent entity relation targets for this filter.
// Relation targets set here are included in filter caching.
// Contrary, relation targets specified in [Filter7.Query] or [Filter7.Batch] are not cached.
func (f *Filter7[A, B, C, D, E, F, G]) Relations(rel ...RelationIndex) *Filter7[A, B, C, D, E, F, G] {
	f.relations = relations(rel).toRelations(&f.world.storage.registry, f.ids, nil, f.relations)
	return f
}

// Register this filter to the world's filter cache.
func (f *Filter7[A, B, C, D, E, F, G]) Register() *Filter7[A, B, C, D, E, F, G] {
	if f.cache != nil {
		panic("filter is already registered, can't register")
	}
	f.cache = f.world.storage.registerFilter(f.Batch())
	return f
}

// Unregister this filter from the world's filter cache.
func (f *Filter7[A, B, C, D, E, F, G]) Unregister() {
	if f.cache == nil {
		panic("filter is not registered, can't unregister")
	}
	f.world.storage.unregisterFilter(f.cache)
	f.cache = nil
}

// Query creates a [Query7] from this filter.
// This must be used each time before iterating a query.
func (f *Filter7[A, B, C, D, E, F, G]) Query(rel ...RelationIndex) Query7[A, B, C, D, E, F, G] {
	if f.cache == nil {
		f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, f.relations, f.tempRelations)
	} else {
		f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, nil, f.tempRelations)
	}
	return newQuery7[A, B, C, D, E, F, G](f.world, f.filter, f.ids, f.tempRelations, f.cache)
}

// Batch creates a [Batch] from this filter.
func (f *Filter7[A, B, C, D, E, F, G]) Batch(rel ...RelationIndex) *Batch {
	// TODO: use cache?
	f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, f.relations, f.tempRelations)
	return &Batch{
		filter:    f.filter,
		relations: f.tempRelations,
	}
}

// Filter8 is a filter for 8 components.
type Filter8[A any, B any, C any, D any, E any, F any, G any, H any] struct {
	world         *World
	ids           []ID
	filter        Filter
	relations     []RelationID
	tempRelations []RelationID
	cache         *cacheEntry
}

// NewFilter8 creates a new [Filter8].
//
// Use [Filter8.Query] to obtain a [Query8].
func NewFilter8[A any, B any, C any, D any, E any, F any, G any, H any](world *World) *Filter8[A, B, C, D, E, F, G, H] {
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

	return &Filter8[A, B, C, D, E, F, G, H]{
		world:  world,
		ids:    ids,
		filter: NewFilter(ids...),
	}
}

// With specifies additional components to filter for.
func (f *Filter8[A, B, C, D, E, F, G, H]) With(comps ...Comp) *Filter8[A, B, C, D, E, F, G, H] {
	for _, c := range comps {
		id := f.world.componentID(c.tp)
		f.filter.mask.Set(id, true)
	}
	return f
}

// Without specifies components to exclude.
func (f *Filter8[A, B, C, D, E, F, G, H]) Without(comps ...Comp) *Filter8[A, B, C, D, E, F, G, H] {
	for _, c := range comps {
		id := f.world.componentID(c.tp)
		f.filter.without.Set(id, true)
		f.filter.hasWithout = true
	}
	return f
}

// Exclusive makes the filter exclusive in the sense that the component composition is matched exactly,
// and no other components are allowed.
func (f *Filter8[A, B, C, D, E, F, G, H]) Exclusive() *Filter8[A, B, C, D, E, F, G, H] {
	f.filter = f.filter.Exclusive()
	return f
}

// Relations sets permanent entity relation targets for this filter.
// Relation targets set here are included in filter caching.
// Contrary, relation targets specified in [Filter8.Query] or [Filter8.Batch] are not cached.
func (f *Filter8[A, B, C, D, E, F, G, H]) Relations(rel ...RelationIndex) *Filter8[A, B, C, D, E, F, G, H] {
	f.relations = relations(rel).toRelations(&f.world.storage.registry, f.ids, nil, f.relations)
	return f
}

// Register this filter to the world's filter cache.
func (f *Filter8[A, B, C, D, E, F, G, H]) Register() *Filter8[A, B, C, D, E, F, G, H] {
	if f.cache != nil {
		panic("filter is already registered, can't register")
	}
	f.cache = f.world.storage.registerFilter(f.Batch())
	return f
}

// Unregister this filter from the world's filter cache.
func (f *Filter8[A, B, C, D, E, F, G, H]) Unregister() {
	if f.cache == nil {
		panic("filter is not registered, can't unregister")
	}
	f.world.storage.unregisterFilter(f.cache)
	f.cache = nil
}

// Query creates a [Query8] from this filter.
// This must be used each time before iterating a query.
func (f *Filter8[A, B, C, D, E, F, G, H]) Query(rel ...RelationIndex) Query8[A, B, C, D, E, F, G, H] {
	if f.cache == nil {
		f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, f.relations, f.tempRelations)
	} else {
		f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, nil, f.tempRelations)
	}
	return newQuery8[A, B, C, D, E, F, G, H](f.world, f.filter, f.ids, f.tempRelations, f.cache)
}

// Batch creates a [Batch] from this filter.
func (f *Filter8[A, B, C, D, E, F, G, H]) Batch(rel ...RelationIndex) *Batch {
	// TODO: use cache?
	f.tempRelations = relations(rel).toRelations(&f.world.storage.registry, f.ids, f.relations, f.tempRelations)
	return &Batch{
		filter:    f.filter,
		relations: f.tempRelations,
	}
}
