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
	Allocs float64
	Bytes  float64
	Factor float64
	Units  string
}

// Result type
type Result struct {
	Name   string
	N      int
	Time   float64
	Allocs float64
	Bytes  float64
}

// CompResult type
type CompResult struct {
	Name     string
	N        int
	TimeMain float64
	TimeCurr float64
	Factor   float64
	Allocs   float64
	Bytes    float64
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
		var allocs int
		for range count {
			res := testing.Benchmark(b.F)
			t += res.T.Nanoseconds()
			n += res.N
			mem += int(res.MemBytes)
			allocs += int(res.MemAllocs)
		}
		b.T = float64(t) / float64(n*b.N)
		b.Bytes = float64(mem) / float64(n*b.N)
		b.Allocs = float64(allocs) / float64(n*b.N)
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

	b.WriteString("Operation;N;Time;Allocs;Bytes\n")

	for i := range benches {
		bench := &benches[i]
		b.WriteString(fmt.Sprintf("%s;%d;%0.2f;%0.2f;%0.2f\n", bench.Name, bench.N, bench.T, bench.Allocs, bench.Bytes))
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

		nVal, _ := strconv.Atoi(record[1])
		timeVal, _ := strconv.ParseFloat(record[2], 64)
		allocsVal, _ := strconv.ParseFloat(record[3], 64)
		bytesVal, _ := strconv.ParseFloat(record[4], 64)

		results = append(results, Result{
			Name:   record[0],
			N:      nVal,
			Time:   timeVal,
			Allocs: allocsVal,
			Bytes:  bytesVal,
		})
	}

	return results, nil
}

// TableToHTML convert benchmark comparison results to HTML.
func TableToHTML(data []CompResult) string {
	html := `
<details>
<summary>Click to expand benchmark results</summary>
<p>
Time is per entity/N, allocations are totals.
Allocations are only shown for current.
</p>
<table>
	<thead>
	<tr>
		<th align="center">N</th>
		<th align="center">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Time main&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</th>
		<th align="center">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Time curr&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</th>
		<th align="center">&nbsp;&nbsp;&nbsp;&nbsp;Factor&nbsp;&nbsp;&nbsp;&nbsp;</th>
		<th align="center">&nbsp;&nbsp;&nbsp;&nbsp;Allocs&nbsp;&nbsp;&nbsp;&nbsp;</th>
		<th align="center">&nbsp;&nbsp;&nbsp;&nbsp;Bytes&nbsp;&nbsp;&nbsp;&nbsp;</th>
	</tr>
	</thead>
	<tbody>`

	improved := 0
	regressed := 0

	name := ""
	for _, r := range data {
		emoji := ""
		if r.Factor <= 0.9 {
			improved++
			emoji = "🚀"
		} else if r.Factor >= 1.1 {
			regressed++
			emoji = "⚠️"
		}

		if name != r.Name {
			html += fmt.Sprintf(`<tr><th colspan="6" align="center">%s</th></tr>`, r.Name) + "\n"
		}

		html += fmt.Sprintf(`            <tr>
            <td align="right">%d</td>
            <td align="right">%.2fns</td>
            <td align="right">%.2fns</td>
            <td align="right">%s %.2f</td>
            <td align="right">%d</td>
            <td align="right">%d</td>
            </tr>`, r.N, r.TimeMain, r.TimeCurr, emoji, r.Factor, int(r.Allocs), int(r.Bytes))

		name = r.Name
	}

	html += `      </tbody>
    </table>
    </details>`

	if regressed == 0 && improved == 0 {
		html = "<p>✅ Benchmarks are stable!</p>\n" + html
	} else {
		if regressed > 0 {
			html = "<p>⚠️ $regressed benchmark regressions detected!</p>\n" + html
		}
		if improved > 0 {
			html = "<p>🚀 $improved benchmark improvements detected!</p>\n" + html
		}
	}

	return html
}
