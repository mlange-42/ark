package stats

// Option flags type
type Option uint8

const (
	// Archetypes included in stats.
	Archetypes Option = iota
	// Tables included in stats. Requires Archetypes.
	Tables
	// Filters included in stats.
	Filters
	// Observers included in stats.
	Observers
)
