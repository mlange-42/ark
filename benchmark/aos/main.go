package main

import (
	"fmt"
	"testing"
)

func main() {
	nValues := []int{100, 1_000, 10_000, 100_000, 1_000_000, 10_000_000}
	arkFunctions := []func(*testing.B, int){
		ark16Byte, ark32Byte, ark64Byte,
	}
	aosFunctions := []func(*testing.B, int){
		aos16Byte, aos32Byte, aos64Byte,
	}
	bytes := []int{16, 32, 64}

	for i := range arkFunctions {
		for _, n := range nValues {
			fn := func(b *testing.B) {
				arkFunctions[i](b, n)
			}
			res := testing.Benchmark(fn)
			tArk := float64(res.T.Nanoseconds()) / float64(n*res.N)

			fn = func(b *testing.B) {
				aosFunctions[i](b, n)
			}
			res = testing.Benchmark(fn)
			tAos := float64(res.T.Nanoseconds()) / float64(n*res.N)

			fmt.Printf("%dB: Ark %0.2fns | Aos %0.2fns (%d entities)\n", bytes[i], tArk, tAos, n)
		}
	}
}
