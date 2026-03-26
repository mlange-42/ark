package main

import (
	"fmt"
	"testing"
)

func main() {
	nValues := []int{100, 1_000, 10_000, 100_000, 1_000_000, 10_000_000}
	arkFunctions := []func(*testing.B, int){
		ark32Byte, ark64Byte, ark128Byte, ark256Byte,
	}
	aosFunctions := []func(*testing.B, int){
		aos32Byte, aos64Byte, aos128Byte, aos256Byte,
	}
	bytes := []int{32, 64, 128, 256}

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

			fmt.Printf("%3dB: Ark %0.2fns | Aos %0.2fns (%d entities)\n", bytes[i], tArk, tAos, n)
		}
	}
}
