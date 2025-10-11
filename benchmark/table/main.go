package main

import (
	"fmt"
	"time"

	"github.com/klauspost/cpuid/v2"
	"github.com/mlange-42/ark/benchmark"
)

const version = "v0.5.3-dev"
const goVersion = "1.25.1"

func main() {
	repetitions := 3

	fmt.Printf("Last run: %s  \n", time.Now().Format(time.RFC1123))
	fmt.Printf("Version: Ark %s  \n", version)
	fmt.Printf("Go version: %s  \n", goVersion)
	fmt.Printf("CPU: %s\n\n", cpuid.CPU.BrandName)

	benchmark.RunBenchmarks("Query", benchesQuery(), repetitions, benchmark.ToMarkdown)
	benchmark.RunBenchmarks("World access", benchesWorld(), repetitions, benchmark.ToMarkdown)
	benchmark.RunBenchmarks("Entities", benchesEntities(), repetitions, benchmark.ToMarkdown)
	benchmark.RunBenchmarks("Entities, batched", benchesEntitiesBatch(), repetitions, benchmark.ToMarkdown)
	benchmark.RunBenchmarks("Components", benchesComponents(), repetitions, benchmark.ToMarkdown)
	benchmark.RunBenchmarks("Components, batched", benchesComponentsBatch(), repetitions, benchmark.ToMarkdown)
	benchmark.RunBenchmarks("Other", benchesOther(), repetitions, benchmark.ToMarkdown)

	fmt.Print("\n\n")
}
