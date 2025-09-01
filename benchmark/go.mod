module github.com/mlange-42/ark/benchmark

go 1.24.0

require (
	github.com/klauspost/cpuid/v2 v2.3.0
	github.com/mlange-42/ark v0.4.6
	github.com/pkg/profile v1.7.0
)

replace github.com/mlange-42/ark v0.4.6 => ..

require (
	github.com/felixge/fgprof v0.9.3 // indirect
	github.com/google/pprof v0.0.0-20211214055906-6f57359322fd // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	golang.org/x/sys v0.35.0 // indirect
)
