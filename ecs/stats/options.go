package stats

// Option flags type.
type Option uint8

const (
	// None discards all stats details.
	None Option = iota
	// Archetypes included in stats.
	Archetypes
	// Tables and archetypes included in stats.
	Tables
	// All details included in stats. Currently equivalent to Tables.
	All = Archetypes | Tables
)
