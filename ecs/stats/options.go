package stats

// Option flags type
type Option uint8

const (
	// Archetypes included in stats.
	Archetypes Option = 1 << iota
	// Tables included in stats. Requires Archetypes.
	Tables
	// None discards all stats details.
	None Option = 0
	// All includes all stats details.
	All Option = Archetypes | Tables
)
