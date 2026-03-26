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
	Name  string
	Label string
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
	colorsLight := map[string]color.RGBA{
		"Ark": {R: 30, G: 30, B: 30, A: 255},
		"AoS": {R: 0, G: 107, B: 230, A: 255},
	}
	colorsDark := map[string]color.RGBA{
		"Ark": {R: 250, G: 250, B: 250, A: 255},
		"AoS": {R: 51, G: 173, B: 255, A: 255},
	}

	bytes := []int{32, 64, 128, 256}

	allSeries := []Series{}

	for f := range 2 {
		for i := range modes[f] {
			series := Series{
				Name:  names[f],
				Label: fmt.Sprintf("%s %3dB", names[f], bytes[i]),
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

	plotResults(allSeries, color.RGBA{R: 250, G: 250, B: 250, A: 255}, color.RGBA{A: 255}, colorsLight, "aos_light.svg")
	plotResults(allSeries, color.RGBA{R: 30, G: 30, B: 30, A: 255}, color.RGBA{R: 243, G: 244, B: 246, A: 255}, colorsDark, "aos_dark.svg")
}

func plotResults(data []Series, bg color.RGBA, fg color.RGBA, colors map[string]color.RGBA, file string) {
	plot.DefaultFont.Variant = "Mono"

	p := plot.New()
	p.BackgroundColor = bg
	p.X.Padding = 0
	p.Y.Padding = 0

	p.X.Label.Text = "entities"
	p.X.Label.TextStyle.Color = fg
	p.X.Tick.Color = fg
	p.X.Tick.Label.Color = fg
	p.X.Color = fg
	p.X.Scale = plot.LogScale{}
	p.X.Tick.Marker = plot.LogTicks{Prec: -1}

	p.Y.Label.Text = "time per entity [ns]"
	p.Y.Label.TextStyle.Color = fg
	p.Y.Tick.Color = fg
	p.Y.Tick.Label.Color = fg
	p.Y.Color = fg
	p.Y.Min = 0

	p.Legend = plot.NewLegend()
	p.Legend.TextStyle.Color = fg
	p.Legend.Top = true
	p.Legend.Left = true

	for i := range data {
		series := &data[i]

		lines, err := plotter.NewLine(series.Data)
		if err != nil {
			panic(err)
		}
		lines.Color = colors[series.Name]
		lines.Width = vg.Points(series.Width)

		p.Add(lines)
		p.Legend.Add(series.Label, lines)
	}

	err := p.Save(460, 300, file)
	if err != nil {
		panic(err)
	}
}
