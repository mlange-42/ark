package main

import (
	"flag"
	"fmt"
	"image/color"
	"math"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Series data
type Series struct {
	Label string
	Color color.RGBA
	Width float64
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

	modes := [][]func(*testing.B, int){arkFunctions, aosFunctions}
	names := []string{"Ark", "AoS"}
	colors := []color.RGBA{{B: 255, A: 255}, {R: 255, A: 255}}

	bytes := []int{32, 64, 128, 256}

	allSeries := []Series{}

	for f := range 2 {
		for i := range modes[f] {
			series := Series{
				Label: fmt.Sprintf("%s %3dB", names[f], bytes[i]),
				Color: colors[f],
				Width: (math.Log2(float64(bytes[i])) - 3) / 2,
			}
			for _, n := range nValues {
				fn := func(b *testing.B) {
					modes[f][i](b, n)
				}
				res := testing.Benchmark(fn)
				t := float64(res.T.Nanoseconds()) / float64(n*res.N)

				fmt.Printf("%3dB: %s %0.2fns (%d entities)\n", bytes[i], names[f], t, n)

				series.Data = append(series.Data, plotter.XY{X: float64(n), Y: t})
			}

			allSeries = append(allSeries, series)
		}
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
		lines.Width = vg.Points(series.Width)

		p.Add(lines)
		p.Legend.Add(series.Label, lines)
	}

	err := p.Save(400, 300, "aos.png")
	if err != nil {
		panic(err)
	}
}
