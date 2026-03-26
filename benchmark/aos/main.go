package main

import (
	"flag"
	"fmt"
	"image/color"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

// Series data
type Series struct {
	Label string
	Color color.RGBA
	Data  plotter.XYs
}

func main() {
	testing.Init()
	flag.Parse()

	nValues := []int{100, 1_000, 10_000, 100_000, 1_000_000, 10_000_000}
	arkFunctions := []func(*testing.B, int){
		ark32Byte, ark64Byte, ark128Byte, ark256Byte,
	}
	aosFunctions := []func(*testing.B, int){
		aos32Byte, aos64Byte, aos128Byte, aos256Byte,
	}
	bytes := []int{32, 64, 128, 256}

	allSeries := []Series{}

	for i := range arkFunctions {
		arkSeries := Series{
			Label: fmt.Sprintf("Ark %3dB", bytes[i]),
			Color: color.RGBA{B: 255, A: 255},
		}
		aosSeries := Series{
			Label: fmt.Sprintf("AoS %3dB", bytes[i]),
			Color: color.RGBA{R: 255, A: 255},
		}
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

			arkSeries.Data = append(arkSeries.Data, plotter.XY{X: float64(n), Y: tArk})
			aosSeries.Data = append(aosSeries.Data, plotter.XY{X: float64(n), Y: tAos})
		}

		allSeries = append(allSeries, arkSeries, aosSeries)
	}

	p := plot.New()
	p.X.Label.Text = "entities"
	p.X.Scale = plot.LogScale{}
	p.X.Tick.Marker = plot.LogTicks{Prec: -1}

	p.Y.Label.Text = "time/entity [ns]"
	p.Y.Min = 0

	p.Legend = plot.NewLegend()
	p.Legend.TextStyle.Font.Variant = "Mono"
	p.Legend.Top = true
	p.Legend.Left = true

	for i := range allSeries {
		series := &allSeries[i]

		lines, err := plotter.NewLine(series.Data)
		if err != nil {
			panic(err)
		}
		lines.Color = series.Color
		p.Add(lines)
		p.Legend.Add(series.Label, lines)
	}

	err := p.Save(400, 300, "aos.png")
	if err != nil {
		panic(err)
	}
}
