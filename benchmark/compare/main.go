package main

import (
	"flag"
	"os"
	"testing"

	"github.com/mlange-42/ark/benchmark"
)

func main() {
	testing.Init()
	flag.Parse()

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
		{Name: "Query create + close", Desc: "", F: queryCreateClose, N: 1},
		{Name: "Query(reg) create + close", Desc: "", F: queryCreateCloseRegistered, N: 1},

		{Name: "Pos/Vel query entities", Desc: "", F: posVelQuery10, N: 10},
		{Name: "Pos/Vel query entities", Desc: "", F: posVelQuery1000, N: 1000},
		{Name: "Pos/Vel query entities", Desc: "", F: posVelQuery100000, N: 100_000},

		{Name: "Pos/Vel query tables", Desc: "", F: posVelQueryTables10, N: 10},
		{Name: "Pos/Vel query tables", Desc: "", F: posVelQueryTables1000, N: 1000},
		{Name: "Pos/Vel query tables", Desc: "", F: posVelQueryTables100000, N: 100_000},

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

		{Name: "RemoveEntity 2 Comp", Desc: "", F: removeEntity2Comp10, N: 10},
		{Name: "RemoveEntity 2 Comp", Desc: "", F: removeEntity2Comp1000, N: 1000},

		{Name: "RemoveEntity 5 Comp", Desc: "", F: removeEntity5Comp10, N: 10},
		{Name: "RemoveEntity 5 Comp", Desc: "", F: removeEntity5Comp1000, N: 1000},

		{Name: "RemoveEntityBatch 2 Comp", Desc: "", F: removeEntityBatch2Comp10, N: 10},
		{Name: "RemoveEntityBatch 2 Comp", Desc: "", F: removeEntityBatch2Comp1000, N: 1000},

		{Name: "RemoveEntityBatch 5 Comp", Desc: "", F: removeEntityBatch5Comp10, N: 10},
		{Name: "RemoveEntityBatch 5 Comp", Desc: "", F: removeEntityBatch5Comp1000, N: 1000},

		{Name: "Add/Remove 1 of 2", Desc: "", F: addRemove10, N: 10},
		{Name: "Add/Remove 1 of 2", Desc: "", F: addRemove1000, N: 1000},

		{Name: "Add/Remove 1 of 2 non-trivial", Desc: "", F: addRemoveNonTrivial10, N: 10},
		{Name: "Add/Remove 1 of 2 non-trivial", Desc: "", F: addRemoveNonTrivial1000, N: 1000},

		{Name: "Add/Remove batch 1 of 2", Desc: "", F: addRemoveBatch10, N: 10},
		{Name: "Add/Remove batch 1 of 2", Desc: "", F: addRemoveBatch1000, N: 1000},
		{Name: "Add/Remove batch 1 of 2", Desc: "", F: addRemoveBatch100000, N: 100000},
	}
}
