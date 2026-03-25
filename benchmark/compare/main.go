package main

import (
	"os"

	"github.com/mlange-42/ark/benchmark"
)

func main() {
	repetitions := 1

	f, err := os.Create("benches.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	formats := []benchmark.Format{
		{Format: benchmark.ToCSV, Writer: f},
	}

	benchmark.RunBenchmarks("Compare", benchesCompare(), repetitions, formats)
}

func benchesCompare() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "Pos/Vel query entities n=10", Desc: "", F: posVelQuery10, N: 10},
		{Name: "Pos/Vel query entities n=1000", Desc: "", F: posVelQuery1000, N: 1000},
		{Name: "Pos/Vel query entities n=100000", Desc: "", F: posVelQuery100000, N: 100_000},

		{Name: "Pos/Vel mapper n=10", Desc: "", F: posVelMap10, N: 10},
		{Name: "Pos/Vel mapper n=1000", Desc: "", F: posVelMap1000, N: 1000},
		{Name: "Pos/Vel mapper n=100000", Desc: "", F: posVelMap100000, N: 100_000},
	}
}
