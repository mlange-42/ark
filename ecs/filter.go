package ecs

// Filter is a filter for components.
type Filter struct {
	mask       Mask
	without    Mask
	hasWithout bool
}

// NewFilter creates a new [Filter] matching the given components.
func NewFilter(ids ...ID) Filter {
	return Filter{
		mask: All(ids...),
	}
}

// Without specifies components to exclude.
// Resets previous excludes.
func (f Filter) Without(ids ...ID) Filter {
	f.without = All(ids...)
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

func (f *Filter) matches(mask *Mask) bool {
	return mask.Contains(&f.mask) && (!f.hasWithout || !mask.ContainsAny(&f.without))
}
