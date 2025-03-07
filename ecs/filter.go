package ecs

// Filter is a filter for components.
//
// It is significantly slower than type-safe generic filters like [Filter2],
// and should only be used when component types are not known at compile time.
type Filter struct {
	filter
	world *World
}

// NewFilter creates a new [Filter] matching the given components.
func NewFilter(world *World, ids ...ID) Filter {
	return Filter{
		world:  world,
		filter: newFilter(ids...),
	}
}

// Without specifies components to exclude.
// Resets previous excludes.
func (f Filter) Without(ids ...ID) Filter {
	f.filter = f.filter.Without(ids...)
	return f
}

// Exclusive makes the filter exclusive in the sense that the component composition is matched exactly,
// and no other components are allowed.
func (f Filter) Exclusive() Filter {
	f.filter = f.filter.Exclusive()
	return f
}

// Query returns a new query matching this filter and the given entity relation targets.
func (f Filter) Query(relations ...RelationID) Query {
	return newQuery(f.world, f.filter, relations)
}

type filter struct {
	mask       bitMask
	without    bitMask
	hasWithout bool
}

func newFilter(ids ...ID) filter {
	return filter{
		mask: newMask(ids...),
	}
}

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
