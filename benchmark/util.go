package benchmark

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
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
	Mem    float64
	Factor float64
	Units  string
}

// Result type
type Result struct {
	Name string
	Time float64
	Mem  float64
}

// Format for writing benchmark results.
type Format struct {
	Format func(string, []Benchmark) string
	Writer io.Writer
}

// RunBenchmarks runs the benchmarks and prints the results.
func RunBenchmarks(title string, benches []Benchmark, count int, format []Format) {
	for i := range benches {
		b := &benches[i]
		var t int64
		var n int
		var mem int
		for range count {
			res := testing.Benchmark(b.F)
			t += res.T.Nanoseconds()
			n += res.N
			mem += int(res.MemBytes)
		}
		b.T = float64(t) / float64(n*b.N)
		b.Mem = float64(mem) / float64(n*b.N)
	}
	for _, f := range format {
		_, err := fmt.Fprint(f.Writer, f.Format(title, benches))
		if err != nil {
			panic(err)
		}
	}
}

// ToMarkdown converts the benchmarks to a markdown table.
func ToMarkdown(title string, benches []Benchmark) string {
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("## %s\n\n", title))

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

// ToCSV converts the benchmarks to a CSV table.
func ToCSV(title string, benches []Benchmark) string {
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("%s;%s;%s\n", "Operation", "Time", "Alloc"))

	for i := range benches {
		bench := &benches[i]
		b.WriteString(fmt.Sprintf("%s;%0.2f;%0.2f\n", bench.Name, bench.T, bench.Mem))
	}

	return b.String()
}

// ReadCSV reade benchmark results from a CSV file.
func ReadCSV(file string) ([]Result, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'

	_, err = r.Read()
	if err != nil {
		return nil, err
	}

	var results []Result

	for {
		record, err := r.Read()
		if err != nil {
			break // EOF is expected
		}

		timeVal, _ := strconv.ParseFloat(record[1], 64)
		memVal, _ := strconv.ParseFloat(record[2], 64)

		results = append(results, Result{
			Name: record[0],
			Time: timeVal,
			Mem:  memVal,
		})
	}

	return results, nil
}
