package plotter

import "testing"

func TestBar(t *testing.T) {
	barChart := NewBarChart(nil)
	barChart.SetTitle("Example Bar Chart")
	barChart.SetYScale(0, 100, 20)
	barChart.SetXLabels([]string{"A", "B", "C", "D", "E"})
	barChart.AddSeries([]float64{30, 60, 80, 40, 90}, "blue", "Series 1")
	barChart.Close()
}
