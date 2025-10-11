package ecs

// UnsafeFilter is a filter for components.
//
// It is significantly slower than type-safe generic filters like [Filter2],
// and should only be used when component types are not known at compile time.
type UnsafeFilter struct {
	filter
	world *World
}

// NewUnsafeFilter creates a new [UnsafeFilter] matching the given components.
func NewUnsafeFilter(world *World, ids ...ID) UnsafeFilter {
	return UnsafeFilter{
		world:  world,
		filter: newFilter(ids...),
	}
}

// Without specifies components to exclude.
// Resets previous excludes.
func (f UnsafeFilter) Without(ids ...ID) UnsafeFilter {
	if len(ids) == 0 {
		return f
	}
	f.filter = f.filter.Without(ids...)
	return f
}

// Exclusive makes the filter exclusive in the sense that the component composition is matched exactly,
// and no other components are allowed.
//
// Overwrites components set via [Filter.Without].
func (f UnsafeFilter) Exclusive() UnsafeFilter {
	f.filter = f.filter.Exclusive()
	return f
}

// Query returns a new query matching this filter and the given entity relation targets.
func (f UnsafeFilter) Query(relations ...Relation) UnsafeQuery {
	f.world.storage.mu.Lock()
	lock := f.world.lock()
	rel := f.world.storage.slices.relationsPool.Get()
	f.world.storage.mu.Unlock()

	rel = relationSlice(relations).ToRelationIDsForUnsafe(f.world, rel)
	return UnsafeQuery{
		world:     f.world,
		filter:    f.filter,
		relations: rel,
		lock:      lock,
		cursor: cursor{
			archetype: -1,
			table:     -1,
			index:     0,
			maxIndex:  -1,
		},
	}
}

// filter is an mask filter for component presence and optional absence.
type filter struct {
	mask       bitMask
	without    bitMask
	cache      cacheID
	hasWithout bool
}

// newFilter creates a new filter for presence of the given components.
func newFilter(ids ...ID) filter {
	return filter{
		mask:  newMask(ids...),
		cache: maxCacheID,
	}
}

// matches this filter against a (archetype) mask.
func (f *filter) matches(mask *bitMask) bool {
	return mask.Contains(&f.mask) && (!f.hasWithout || !mask.ContainsAny(&f.without))
}

// Without specifies components to exclude.
// Resets previous excludes.
func (f filter) Without(ids ...ID) filter {
	f.without = newMask(ids...)
	f.hasWithout = true
	return f
}

// Exclusive makes the filter exclusive in the sense that the component composition is matched exactly,
// and no other components are allowed.
func (f filter) Exclusive() filter {
	f.without = f.mask.Not()
	f.hasWithout = true
	return f
}
