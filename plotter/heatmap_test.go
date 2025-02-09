package plotter_test

import (
	"go-numeric/plotter"
	"testing"
)

func TestHeatmap(t *testing.T) {
	heatmap := plotter.NewHeatmap(nil)
	heatmap.SetTitle("Example Heatmap")

	data := [][]float64{
		{0.1, 0.3, 0.5, 0.7, 1.0},
		{0.2, 0.4, 0.6, 0.8, 0.9},
		{0.3, 0.5, 0.7, 0.9, 0.6},
		{0.4, 0.6, 0.8, 0.5, 0.3},
	}

	xLabels := []string{"A", "B", "C", "D", "E"}
	yLabels := []string{"W", "X", "Y", "Z"}

	heatmap.SetData(data)
	heatmap.SetXLabels(xLabels)
	heatmap.SetYLabels(yLabels)
	heatmap.Close()
}
