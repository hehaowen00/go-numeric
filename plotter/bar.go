package plotter

import (
	"fmt"
	"io"
	"os"

	"github.com/ajstarks/svgo"
)

type BarChart struct {
	wr      io.WriteCloser
	title   string
	xLabels []string
	yStart  float64
	yEnd    float64
	yInc    float64
	series  []barSeries
	width   int
	height  int
	padding int
}

type barSeries struct {
	values []float64
	color  string
	label  string
}

func NewBarChart(wr io.Writer) *BarChart {
	return &BarChart{
		wr:      nil,
		width:   500,
		height:  500,
		padding: 50,
	}
}

func (b *BarChart) SetTitle(title string) {
	b.title = title
}

func (b *BarChart) SetYScale(start, end, increments float64) {
	b.yStart, b.yEnd, b.yInc = start, end, increments
}

func (b *BarChart) SetXLabels(labels []string) {
	b.xLabels = labels
}

func (b *BarChart) AddSeries(values []float64, color, label string) {
	b.series = append(b.series, barSeries{values, color, label})
}

func (b *BarChart) Close() {
	f, err := os.OpenFile("bar.svg", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	canvas := svg.New(f)
	canvas.Start(b.width, b.height)
	b.drawAxes(canvas)
	b.drawBars(canvas)
	canvas.End()
}

func (b *BarChart) drawAxes(canvas *svg.SVG) {
	x0 := b.padding
	y0 := b.height - b.padding

	canvas.Line(x0, y0, x0, b.padding, "stroke:black;stroke-width:2")

	for y := b.yStart; y <= b.yEnd; y += b.yInc {
		py := b.scaleY(y)
		canvas.Line(x0-5, py, x0+5, py, "stroke:black")
		canvas.Text(x0-10, py+4, fmt.Sprintf("%.1f", y), "text-anchor:end;font-size:12px;fill:black")
	}

	canvas.Text(b.width/2, b.padding/2, b.title, "text-anchor:middle;font-size:16px;fill:black")
}

func (b *BarChart) drawBars(canvas *svg.SVG) {
	barWidth := (b.width - 2*b.padding) / len(b.xLabels)
	offset := barWidth / len(b.series)

	for i, label := range b.xLabels {
		x := b.padding + i*barWidth
		canvas.Text(x+barWidth/2, b.height-b.padding+20, label, "text-anchor:middle;font-size:12px;fill:black")

		for j, series := range b.series {
			val := series.values[i]
			y := b.scaleY(val)
			canvas.Rect(x+j*offset, y, offset-5, b.height-b.padding-y, fmt.Sprintf("fill:%s", series.color))
		}
	}
}

func (b *BarChart) scaleY(y float64) int {
	return b.height - b.padding - int(((y-b.yStart)/(b.yEnd-b.yStart))*float64(b.height-2*b.padding))
}
