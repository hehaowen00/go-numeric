package plotter

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/ajstarks/svgo"
)

type Histogram struct {
	wr       io.WriteCloser
	title    string
	width    int
	height   int
	padding  int
	data     []float64
	binCount int
	xLabel   string
	yLabel   string
	barColor string
}

func NewHistogram(wr io.Writer) *Histogram {
	return &Histogram{
		width:    800,
		height:   500,
		padding:  50,
		binCount: 12,
		barColor: "blue",
	}
}

func (h *Histogram) SetTitle(title string) {
	h.title = title
}

func (h *Histogram) SetXLabel(label string) {
	h.xLabel = label
}

func (h *Histogram) SetYLabel(label string) {
	h.yLabel = label
}

func (h *Histogram) SetData(data []float64, bins int) {
	h.data = data
	h.binCount = bins
}

func (h *Histogram) Close() {
	f, err := os.OpenFile("histogram.svg", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	canvas := svg.New(f)
	canvas.Start(h.width, h.height)
	h.drawHistogram(canvas)
	canvas.End()
}

func (h *Histogram) drawHistogram(canvas *svg.SVG) {
	if len(h.data) == 0 {
		return
	}

	plotX := h.padding
	plotY := h.padding
	plotWidth := h.width - 2*h.padding
	plotHeight := h.height - 2*h.padding

	minVal, maxVal := h.data[0], h.data[0]
	for _, v := range h.data {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	binWidth := (maxVal - minVal) / float64(h.binCount)
	bins := make([]int, h.binCount)

	for _, v := range h.data {
		binIndex := int(math.Floor((v - minVal) / binWidth))
		if binIndex >= h.binCount {
			binIndex = h.binCount - 1
		}
		bins[binIndex]++
	}

	maxFreq := 0
	for _, count := range bins {
		if count > maxFreq {
			maxFreq = count
		}
	}

	canvas.Text(h.width/2, 20, h.title, "text-anchor:middle;font-size:16px;fill:black")

	canvas.Line(plotX, plotY+plotHeight, plotX+plotWidth, plotY+plotHeight, "stroke:black;stroke-width:2")
	canvas.Line(plotX, plotY, plotX, plotY+plotHeight, "stroke:black;stroke-width:2")

	barSpacing := plotWidth / h.binCount
	scaleY := func(freq int) int {
		return plotY + plotHeight - int(float64(freq)/float64(maxFreq)*float64(plotHeight))
	}

	for i, count := range bins {
		x := plotX + i*barSpacing
		y := scaleY(count)
		height := plotY + plotHeight - y

		canvas.Rect(x, y, barSpacing-2, height, fmt.Sprintf("fill:%s;", h.barColor))

		binLabel := fmt.Sprintf("%.1f", minVal+float64(i)*binWidth)
		canvas.Text(x+barSpacing/2, plotY+plotHeight+15, binLabel, "text-anchor:middle;font-size:12px;fill:black")
	}

	for i := 0; i <= maxFreq; i += int(math.Max(1, float64(maxFreq)/5)) {
		y := scaleY(i)
		canvas.Text(plotX-10, y, fmt.Sprintf("%d", i), "text-anchor:end;font-size:12px;fill:black")
	}

	canvas.Text(h.width/2, h.height-10, h.xLabel, "text-anchor:middle;font-size:14px;fill:black")
	canvas.Text(10, h.height/2, h.yLabel, "text-anchor:middle;font-size:14px;fill:black;transform:rotate(-90deg)")
}
