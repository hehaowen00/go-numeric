package dataframe

import (
	"slices"
	"sort"
)

type String struct {
	data []string
}

func NewString(data ...string) *String {
	return &String{
		data: data,
	}
}

func (col *String) New() IColumn {
	return &String{}
}

func (col *String) Len() int {
	return len(col.data)
}

func (col *String) Data() []string {
	return slices.Clone(col.data)
}

func (col *String) Extend(length int) {
	diff := length - len(col.data)
	if diff > 0 {
		col.data = append(col.data, make([]string, diff)...)
	}
}

func (col *String) Index(idx int) any {
	return col.data[idx]
}

func (col *String) Clone() IColumn {
	newData := make([]string, len(col.data))
	copy(newData, col.data)
	return &String{data: newData}
}

func (col *String) DeleteRow(index int) {
	col.data = append(col.data[:index], col.data[index+1:]...)
}

func (col *String) Set(index int, value any) {
	col.data[index] = value.(string)
}

func (col *String) Append(value string) {
	col.data = append(col.data, value)
}

func (col *String) Head() (string, bool) {
	if len(col.data) == 0 {
		return "", false
	}
	return col.data[0], true
}

func (col *String) Tail() []string {
	if len(col.data) <= 1 {
		return nil
	}
	return col.data[1:]
}

func (col *String) Last() (string, bool) {
	if len(col.data) == 0 {
		return "", false
	}
	return col.data[len(col.data)-1], true
}

func (col *String) SortBy(asc bool) {
	sort.Slice(col.data, func(i, j int) bool {
		if asc {
			return col.data[i] < col.data[j]
		}
		return col.data[i] > col.data[j]
	})
}
