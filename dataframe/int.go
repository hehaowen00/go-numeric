package dataframe

import (
	"slices"
	"sort"
)

type IColumn interface {
	Len() int
	Extend(int)
	Index(int) any
	New() IColumn
	Clone() IColumn
	DeleteRow(index int)
	Set(index int, value any)
}

type Int struct {
	data []int64
}

func NewInt(data ...int64) *Int {
	return &Int{
		data: data,
	}
}

func (col *Int) New() IColumn {
	return &Int{}
}

func (col *Int) Len() int {
	return len(col.data)
}

func (col *Int) Data() []int64 {
	return slices.Clone(col.data)
}

func (col *Int) Extend(length int) {
	diff := length - len(col.data)
	if diff > 0 {
		col.data = append(col.data, make([]int64, diff)...)
	}
}

func (col *Int) Index(idx int) any {
	return col.data[idx]
}

func (col *Int) Clone() IColumn {
	newData := make([]int64, len(col.data))
	copy(newData, col.data)
	return &Int{data: newData}
}

func (col *Int) DeleteRow(index int) {
	col.data = append(col.data[:index], col.data[index+1:]...)
}

func (col *Int) Set(index int, value any) {
	col.data[index] = value.(int64)
}

func (col *Int) Append(value int64) {
	col.data = append(col.data, value)
}

func (col *Int) Head() (int64, bool) {
	if len(col.data) == 0 {
		return 0, false
	}
	return col.data[0], true
}

func (col *Int) Tail() []int64 {
	if len(col.data) <= 1 {
		return nil
	}
	return col.data[1:]
}

func (col *Int) Last() (int64, bool) {
	if len(col.data) == 0 {
		return 0, false
	}
	return col.data[len(col.data)-1], true
}

func (col *Int) Slice(start, end int) []int64 {
	return slices.Clone(col.data[start:end])
}

func (col *Int) SortBy(asc bool) {
	sort.Slice(col.data, func(i, j int) bool {
		if asc {
			return col.data[i] < col.data[j]
		}
		return col.data[i] > col.data[j]
	})
}

func (col *Int) Add(other *Int) {
	for i := range col.data {
		col.data[i] += other.data[i]
	}
}

func (col *Int) Sub(other *Int) {
	for i := range col.data {
		col.data[i] -= other.data[i]
	}
}

func (col *Int) Mul(other *Int) {
	for i := range col.data {
		col.data[i] *= other.data[i]
	}
}

func (col *Int) Div(other *Int) {
	for i := range col.data {
		if other.data[i] == 0 {
			panic("division by zero")
		}
		col.data[i] /= other.data[i]
	}
}
