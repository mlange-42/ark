module github.com/mlange-42/ark/benchmark

go 1.24.0

require (
	github.com/mlange-42/ark v0.2.0
	github.com/pkg/profile v1.7.0
)

replace github.com/mlange-42/ark v0.2.0 => ..

require (
	github.com/felixge/fgprof v0.9.3 // indirect
	github.com/google/pprof v0.0.0-20211214055906-6f57359322fd // indirect
)
