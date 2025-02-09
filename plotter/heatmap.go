package plotter

import (
	"fmt"
	"io"
	"os"

	"github.com/ajstarks/svgo"
)

type Heatmap struct {
	wr         io.WriteCloser
	title      string
	data       [][]float64
	xLabels    []string
	yLabels    []string
	width      int
	height     int
	padding    int
	labelSpace int
}

func NewHeatmap(wr io.Writer) *Heatmap {
	return &Heatmap{
		wr:         nil,
		width:      500,
		height:     500,
		padding:    50,
		labelSpace: 30,
	}
}

func (h *Heatmap) SetTitle(title string) {
	h.title = title
}

func (h *Heatmap) SetData(data [][]float64) {
	h.data = data
}

func (h *Heatmap) SetXLabels(labels []string) {
	h.xLabels = labels
}

func (h *Heatmap) SetYLabels(labels []string) {
	h.yLabels = labels
}

func (h *Heatmap) Close() {
	f, err := os.OpenFile("heatmap.svg", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	canvas := svg.New(f)
	canvas.Start(h.width, h.height)
	h.drawHeatmap(canvas)
	canvas.End()
}

func (h *Heatmap) drawHeatmap(canvas *svg.SVG) {
	rows := len(h.data)
	cols := len(h.data[0])

	cellWidth := (h.width - 2*h.padding - h.labelSpace) / cols
	cellHeight := (h.height - 2*h.padding - h.labelSpace) / rows

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			val := h.data[i][j]
			color := h.colorScale(val)
			x := h.padding + h.labelSpace + j*cellWidth
			y := h.padding + i*cellHeight

			canvas.Rect(x, y, cellWidth, cellHeight, fmt.Sprintf("fill:%s;stroke:black;stroke-width:1", color))

			textColor := "black"
			if val > 0.5 {
				textColor = "white"
			}
			canvas.Text(x+cellWidth/2, y+cellHeight/2+5, fmt.Sprintf("%.2f", val),
				fmt.Sprintf("text-anchor:middle;font-size:12px;fill:%s", textColor))
		}
	}

	if len(h.xLabels) == cols {
		for j, label := range h.xLabels {
			x := h.padding + h.labelSpace + j*cellWidth + cellWidth/2
			y := h.padding + rows*cellHeight + 20
			canvas.Text(x, y, label, "text-anchor:middle;font-size:14px;fill:black")
		}
	}

	if len(h.yLabels) == rows {
		for i, label := range h.yLabels {
			x := h.padding + h.labelSpace - 10
			y := h.padding + i*cellHeight + cellHeight/2 + 5
			canvas.Text(x, y, label, "text-anchor:end;font-size:14px;fill:black")
		}
	}

	canvas.Text(h.width/2, h.padding/2, h.title, "text-anchor:middle;font-size:16px;fill:black")
}

func (h *Heatmap) colorScale(value float64) string {
	red := int(255 * value)
	blue := 255 - red
	return fmt.Sprintf("rgb(%d,0,%d)", red, blue)
}
