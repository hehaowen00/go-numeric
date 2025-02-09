package plotter

import (
	"fmt"
	"io"
	"os"

	"github.com/ajstarks/svgo"
)

type ScatterPlot struct {
	wr      io.WriteCloser
	title   string
	xStart  float64
	xEnd    float64
	xInc    float64
	yStart  float64
	yEnd    float64
	yInc    float64
	series  []plotSeries
	width   int
	height  int
	padding int
}

func NewScatterPlot(wr io.Writer) *ScatterPlot {
	return &ScatterPlot{
		wr:      nil,
		width:   500,
		height:  500,
		padding: 50,
	}
}

func (p *ScatterPlot) SetTitle(title string) {
	p.title = title
}

func (p *ScatterPlot) SetXScale(start, end, increments float64) {
	p.xStart, p.xEnd, p.xInc = start, end, increments
}

func (p *ScatterPlot) SetYScale(start, end, increments float64) {
	p.yStart, p.yEnd, p.yInc = start, end, increments
}

func (p *ScatterPlot) AddSeries(
	x []float64,
	y []float64,
	marker string,
	label string,
	color string,
	style string,
) {
	p.series = append(p.series, plotSeries{x, y, marker, label, color, style})
}

func (p *ScatterPlot) Close() {
	f, err := os.OpenFile("scatter.svg", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	canvas := svg.New(f)
	canvas.Start(p.width, p.height)
	p.drawAxes(canvas)
	p.drawSeries(canvas)
	canvas.End()
}

func (p *ScatterPlot) drawAxes(canvas *svg.SVG) {
	x0 := p.padding
	y0 := p.height - p.padding
	x1 := p.width - p.padding
	y1 := p.padding

	canvas.Line(x0, y0, x1, y0, "stroke:black;stroke-width:2")
	canvas.Line(x0, y0, x0, y1, "stroke:black;stroke-width:2")

	for x := p.xStart; x <= p.xEnd; x += p.xInc {
		px := p.scaleX(x)
		canvas.Line(px, y0-5, px, y0+5, "stroke:black")
		canvas.Text(px, y0+20, fmt.Sprintf("%.1f", x), "text-anchor:middle;font-size:12px;fill:black")
	}

	for y := p.yStart; y <= p.yEnd; y += p.yInc {
		py := p.scaleY(y)
		canvas.Line(x0-5, py, x0+5, py, "stroke:black")
		canvas.Text(x0-10, py+4, fmt.Sprintf("%.1f", y), "text-anchor:end;font-size:12px;fill:black")
	}

	canvas.Text(p.width/2, p.padding/2, p.title, "text-anchor:middle;font-size:16px;fill:black")
}

func (p *ScatterPlot) drawSeries(canvas *svg.SVG) {
	for _, s := range p.series {
		for i := 1; i < len(s.x); i++ {
			x1 := p.scaleX(s.x[i-1])
			y1 := p.scaleY(s.y[i-1])
			x2 := p.scaleX(s.x[i])
			y2 := p.scaleY(s.y[i])

			canvas.Circle(x1, y1, 3, fmt.Sprintf("fill:%s", s.color))
			canvas.Circle(x2, y2, 3, fmt.Sprintf("fill:%s", s.color))
		}
	}
}

func (p *ScatterPlot) scaleX(x float64) int {
	return int(((x-p.xStart)/(p.xEnd-p.xStart))*float64(p.width-2*p.padding)) + p.padding
}

func (p *ScatterPlot) scaleY(y float64) int {
	return p.height - p.padding - int(((y-p.yStart)/(p.yEnd-p.yStart))*float64(p.height-2*p.padding))
}
