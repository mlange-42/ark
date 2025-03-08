package benchmark

import (
	"fmt"
	"strings"
	"testing"
)

// Benchmark represents a benchmark to be run.
type Benchmark struct {
	Name   string
	Desc   string
	F      func(b *testing.B)
	N      int
	T      float64
	Factor float64
	Units  string
}

// RunBenchmarks runs the benchmarks and prints the results.
func RunBenchmarks(title string, benches []Benchmark, format func([]Benchmark) string) {
	for i := range benches {
		b := &benches[i]
		res := testing.Benchmark(b.F)
		b.T = float64(res.T.Nanoseconds()) / float64(res.N*b.N)
	}
	fmt.Printf("## %s\n\n%s", title, format(benches))
}

// ToMarkdown converts the benchmarks to a markdown table.
func ToMarkdown(benches []Benchmark) string {
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("| %-38s | %-12s | %-28s |\n", "Operation", "Time", "Remark"))
	b.WriteString(fmt.Sprintf("|%s|%s:|%s|\n", strings.Repeat("-", 40), strings.Repeat("-", 13), strings.Repeat("-", 30)))

	for i := range benches {
		bench := &benches[i]
		factor := bench.Factor
		if factor == 0 {
			factor = 1
		}
		units := bench.Units
		if units == "" {
			units = "ns"
		}

		t := fmt.Sprintf("%.1f %s", bench.T*factor, units)
		b.WriteString(fmt.Sprintf("| %-38s | %12s | %-28s |\n", bench.Name, t, bench.Desc))
	}
	b.WriteString("\n")

	return b.String()
}
