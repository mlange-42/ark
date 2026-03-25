package main

import (
	"os"

	"github.com/mlange-42/ark/benchmark"
)

func main() {
	repetitions := 1

	f, err := os.Create("bench.csv")
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
		{Name: "Pos/Vel query entities", Desc: "", F: posVelQuery10, N: 10},
		{Name: "Pos/Vel query entities", Desc: "", F: posVelQuery1000, N: 1000},
		{Name: "Pos/Vel query entities", Desc: "", F: posVelQuery100000, N: 100_000},

		{Name: "Pos/Vel mapper", Desc: "", F: posVelMap10, N: 10},
		{Name: "Pos/Vel mapper", Desc: "", F: posVelMap1000, N: 1000},
		{Name: "Pos/Vel mapper", Desc: "", F: posVelMap100000, N: 100_000},

		{Name: "NewEntity 2 Comp", Desc: "", F: newEntity2Comp10, N: 10},
		{Name: "NewEntity 2 Comp", Desc: "", F: newEntity2Comp1000, N: 1000},

		{Name: "NewEntity 5 Comp", Desc: "", F: newEntity5Comp10, N: 10},
		{Name: "NewEntity 5 Comp", Desc: "", F: newEntity5Comp1000, N: 1000},

		{Name: "NewBatch 2 Comp", Desc: "", F: newBatch2Comp10, N: 10},
		{Name: "NewBatch 2 Comp", Desc: "", F: newBatch2Comp1000, N: 1000},

		{Name: "NewBatch 5 Comp", Desc: "", F: newBatch5Comp10, N: 10},
		{Name: "NewBatch 5 Comp", Desc: "", F: newBatch5Comp1000, N: 1000},
	}
}
