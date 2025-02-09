package dataframe

import (
	"slices"
	"sort"
)

type Float struct {
	data []float64
}

func NewFloat(data ...float64) *Float {
	return &Float{
		data: data,
	}
}

func (col *Float) New() IColumn {
	return &Float{}
}

func (col *Float) Len() int {
	return len(col.data)
}

func (col *Float) Data() []float64 {
	return slices.Clone(col.data)
}

func (col *Float) Extend(length int) {
	diff := length - len(col.data)
	if diff > 0 {
		col.data = append(col.data, make([]float64, diff)...)
	}
}

func (col *Float) Index(idx int) any {
	return col.data[idx]
}

func (col *Float) Clone() IColumn {
	newData := make([]float64, len(col.data))
	copy(newData, col.data)
	return &Float{data: newData}
}

func (col *Float) DeleteRow(index int) {
	col.data = append(col.data[:index], col.data[index+1:]...)
}

func (col *Float) Set(index int, value any) {
	col.data[index] = value.(float64)
}

func (col *Float) Append(value float64) {
	col.data = append(col.data, value)
}

func (col *Float) Head() (float64, bool) {
	if len(col.data) == 0 {
		return 0, false
	}

	return col.data[0], true
}

func (col *Float) Tail() []float64 {
	if len(col.data) <= 1 {
		return nil
	}
	return col.data[1:]
}

func (col *Float) Last() (float64, bool) {
	if len(col.data) == 0 {
		return 0, false
	}
	return col.data[len(col.data)-1], true
}

func (col *Float) Slice(start, end int) []float64 {
	return slices.Clone(col.data[start:end])
}

func (col *Float) SortBy(asc bool) {
	sort.Slice(col.data, func(i, j int) bool {
		if asc {
			return col.data[i] < col.data[j]
		}
		return col.data[i] > col.data[j]
	})
}

func (col *Float) Add(other *Float) {
	for i := range col.data {
		col.data[i] += other.data[i]
	}
}

func (col *Float) Sub(other *Float) {
	for i := range col.data {
		col.data[i] -= other.data[i]
	}
}

func (col *Float) Mul(other *Float) {
	for i := range col.data {
		col.data[i] *= other.data[i]
	}
}

func (col *Float) Div(other *Float) {
	for i := range col.data {
		if other.data[i] == 0 {
			panic("division by zero")
		}
		col.data[i] /= other.data[i]
	}
}
