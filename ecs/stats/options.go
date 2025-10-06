package stats

// Option flags type
type Option uint8

const (
	// Archetypes included in stats.
	Archetypes Option = 1 << iota
	// Tables included in stats. Requires Archetypes.
	Tables
	// Filters included in stats.
	Filters
	// Observers included in stats.
	Observers
)

const (
	// None discards all stats details.
	None Option = iota
	// All includes all stats details.
	All Option = Archetypes | Tables | Filters | Observers
)
