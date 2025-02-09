package plotter

import "testing"

func TestHistogram(t *testing.T) {
	plot := NewHistogram(nil)
	plot.SetTitle("Histogram of Random Data")
	plot.SetXLabel("Value Bins")
	plot.SetYLabel("Frequency")

	data := []float64{1, 2, 2.5, 3, 3.5, 3.8, 4, 4.1, 4.5, 4.8, 5, 5.2, 5.5, 6, 6.2, 6.5, 7, 7.3, 7.5, 8, 8.1, 8.3, 9, 9.5}
	plot.SetData(data, 5)
	plot.Close()
}
