package main

import (
	"fmt"

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

	compare(dataOld, dataNew)
}

func compare(dataOld, dataNew [][]benchmark.Result) {}
