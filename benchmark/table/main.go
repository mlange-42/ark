package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/klauspost/cpuid/v2"
	"github.com/mlange-42/ark/benchmark"
)

const version = "v0.7.1"
const goVersion = "1.25.4"

func main() {
	testing.Init()
	flag.Parse()

	repetitions, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	fmt.Printf("Last run: %s  \n", time.Now().Format(time.RFC1123))
	fmt.Printf("Version: Ark %s  \n", version)
	fmt.Printf("Go version: %s  \n", goVersion)
	fmt.Printf("CPU: %s\n\n", cpuid.CPU.BrandName)

	formats := []benchmark.Format{
		{Format: benchmark.ToMarkdown, Writer: os.Stdout},
	}

	benchmark.RunBenchmarks("Query", benchesQuery(), repetitions, formats)
	benchmark.RunBenchmarks("World access", benchesWorld(), repetitions, formats)
	benchmark.RunBenchmarks("Entities", benchesEntities(), repetitions, formats)
	benchmark.RunBenchmarks("Entities, batched", benchesEntitiesBatch(), repetitions, formats)
	benchmark.RunBenchmarks("Components", benchesComponents(), repetitions, formats)
	benchmark.RunBenchmarks("Components, batched", benchesComponentsBatch(), repetitions, formats)
	benchmark.RunBenchmarks("Other", benchesOther(), repetitions, formats)

	fmt.Print("\n\n")
}
