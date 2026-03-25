package main

import (
	"fmt"
	"math"
	"sort"

	"github.com/mlange-42/ark/benchmark"
)

func main() {
	n := 3

	dataOld := [][]benchmark.Result{}
	dataNew := [][]benchmark.Result{}

	for i := range n {
		d, err := benchmark.ReadCSV(fmt.Sprintf("bench_main_%d.csv", i+1))
		if err != nil {
			panic(err)
		}
		dataOld = append(dataOld, d)

		d, err = benchmark.ReadCSV(fmt.Sprintf("bench_current_%d.csv", i+1))
		if err != nil {
			panic(err)
		}
		dataNew = append(dataNew, d)
	}

	result := compare(dataOld, dataNew)
	fmt.Println(result)
}

func compareTables(dataOld, dataNew []benchmark.Result) []benchmark.CompResult {
	dictA := map[string]benchmark.Result{}
	dictB := map[string]benchmark.Result{}

	for _, b := range dataOld {
		dictA[fmt.Sprintf("%s %07d", b.Name, b.N)] = b
	}
	for _, b := range dataNew {
		dictB[fmt.Sprintf("%s %07d", b.Name, b.N)] = b
	}

	keysMap := map[string]struct{}{}
	for key := range dictA {
		keysMap[key] = struct{}{}
	}
	for key := range dictB {
		keysMap[key] = struct{}{}
	}

	keys := make([]string, 0, len(keysMap))
	for k := range keysMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := []benchmark.CompResult{}
	for _, key := range keys {
		out := benchmark.CompResult{}
		out.TimeMain = math.NaN()
		out.TimeCurr = math.NaN()
		out.Factor = math.NaN()
		out.Allocs = -1
		out.Bytes = -1

		if b, ok := dictA[key]; ok {
			out.Name = b.Name
			out.N = b.N
			out.TimeMain = b.Time
		}

		if b, ok := dictB[key]; ok {
			out.Name = b.Name
			out.N = b.N
			out.TimeCurr = b.Time
			out.Allocs = b.Allocs
			out.Bytes = b.Bytes
		}
		out.Factor = out.TimeCurr / out.TimeMain

		result = append(result, out)
	}

	return result
}

func compare(dataOld, dataNew [][]benchmark.Result) []benchmark.CompResult {
	comp := [][]benchmark.CompResult{}

	for i := range dataOld {
		result := compareTables(dataOld[i], dataNew[i])
		comp = append(comp, result)
	}

	count := len(comp)
	result := []benchmark.CompResult{}
	for i := range comp[0] {
		out := benchmark.CompResult{}
		for _, c := range comp {
			row := c[i]
			out.Name = row.Name
			out.N = row.N
			out.TimeMain += row.TimeMain
			out.TimeCurr += row.TimeCurr
			out.Factor += row.Factor
			out.Allocs += row.Allocs
			out.Bytes += row.Bytes
		}
		out.TimeMain /= float64(count)
		out.TimeCurr /= float64(count)
		out.Factor /= float64(count)
		out.Allocs /= float64(count)
		out.Bytes /= float64(count)

		result = append(result, out)
	}

	return result
}
