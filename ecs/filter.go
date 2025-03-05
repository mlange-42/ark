package ecs

// Filter is a filter for components.
//
// It is significantly slower than type-safe generic filters like [Filter2],
// and should only be used when component types are not known at compile time.
type Filter struct {
	world      *World
	mask       Mask
	without    Mask
	hasWithout bool
}

// NewFilter creates a new [Filter] matching the given components.
func NewFilter(world *World, ids ...ID) Filter {
	return Filter{
		world: world,
		mask:  NewMask(ids...),
	}
}

// Without specifies components to exclude.
// Resets previous excludes.
func (f Filter) Without(ids ...ID) Filter {
	f.without = NewMask(ids...)
	f.hasWithout = true
	return f
}

// Exclusive makes the filter exclusive in the sense that the component composition is matched exactly,
// and no other components are allowed.
func (f Filter) Exclusive() Filter {
	f.without = f.mask.Not()
	f.hasWithout = true
	return f
}

// Query returns a new query matching this filter and the given entity relation targets.
// This is a synonym for [Unsafe.Query].
func (f Filter) Query(relations ...RelationID) Query {
	return newQuery(f.world, f, relations)
}

func (f *Filter) matches(mask *Mask) bool {
	return mask.Contains(&f.mask) && (!f.hasWithout || !mask.ContainsAny(&f.without))
}
