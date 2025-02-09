package dataframe

import (
	"slices"
	"sort"
)

type Bool struct {
	data []bool
}

func NewBool(data ...bool) *Bool {
	return &Bool{
		data: data,
	}
}

func (col *Bool) New() IColumn {
	return &Bool{}
}

func (col *Bool) Len() int {
	return len(col.data)
}

func (col *Bool) Data() []bool {
	return slices.Clone(col.data)
}

func (col *Bool) Extend(length int) {
	diff := length - len(col.data)
	if diff > 0 {
		col.data = append(col.data, make([]bool, diff)...)
	}
}

func (col *Bool) Index(idx int) any {
	return col.data[idx]
}

func (col *Bool) Clone() IColumn {
	newData := make([]bool, len(col.data))
	copy(newData, col.data)
	return &Bool{data: newData}
}

func (col *Bool) DeleteRow(index int) {
	col.data = append(col.data[:index], col.data[index+1:]...)
}

func (col *Bool) Set(index int, value any) {
	col.data[index] = value.(bool)
}

func (col *Bool) Append(value bool) {
	col.data = append(col.data, value)
}

func (col *Bool) Head() (bool, bool) {
	if len(col.data) == 0 {
		return false, false
	}
	return col.data[0], true
}

func (col *Bool) Tail() []bool {
	if len(col.data) <= 1 {
		return nil
	}
	return col.data[1:]
}

func (col *Bool) Last() (bool, bool) {
	if len(col.data) == 0 {
		return false, false
	}

	return col.data[len(col.data)-1], true
}

func (col *Bool) Slice(start, end int) []bool {
	return slices.Clone(col.data[start:end])
}

func (col *Bool) SortBy(asc bool) {
	sort.Slice(col.data, func(i, j int) bool {
		if asc {
			return !col.data[i] && col.data[j]
		}
		return col.data[i] && !col.data[j]
	})
}
