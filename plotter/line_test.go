package plotter_test

import (
	"go-numeric/plotter"
	"testing"
)

func TestLine(t *testing.T) {
	plot := plotter.NewLinePlot(nil)
	plot.SetTitle("Example Line Plot")
	plot.SetXScale(0, 10, 2)   // X-axis from 0 to 10 with increments of 2
	plot.SetYScale(0, 100, 20) // Y-axis from 0 to 100 with increments of 20
	plot.AddSeries(
		[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		[]float64{10, 30, 50, 40, 80, 60, 90, 70, 100, 95},
		"circle", "Series 1", "blue", "solid",
	)
	plot.Close()
}
