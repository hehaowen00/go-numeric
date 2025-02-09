package dataframe

import (
	"slices"
	"sort"
	"time"
)

type Time struct {
	data []time.Time
}

func NewTime(data ...time.Time) *Time {
	return &Time{
		data: data,
	}
}

func (col *Time) New() IColumn {
	return &Time{}
}

func (col *Time) Len() int {
	return len(col.data)
}

func (col *Time) Data() []time.Time {
	return slices.Clone(col.data)
}

func (col *Time) Extend(length int) {
	diff := length - len(col.data)
	if diff > 0 {
		col.data = append(col.data, make([]time.Time, diff)...)
	}
}

func (col *Time) Index(idx int) any {
	return col.data[idx]
}

func (col *Time) Clone() IColumn {
	newData := make([]time.Time, len(col.data))
	copy(newData, col.data)
	return &Time{data: newData}
}

func (col *Time) DeleteRow(index int) {
	col.data = append(col.data[:index], col.data[index+1:]...)
}

func (col *Time) Set(index int, value any) {
	col.data[index] = value.(time.Time)
}

func (col *Time) Append(value time.Time) {
	col.data = append(col.data, value)
}

func (col *Time) Head() (time.Time, bool) {
	if len(col.data) == 0 {
		return time.Time{}, false
	}
	return col.data[0], true
}

func (col *Time) Tail() []time.Time {
	if len(col.data) <= 1 {
		return nil
	}

	return col.data[1:]
}

func (col *Time) Last() (time.Time, bool) {
	if len(col.data) == 0 {
		return time.Time{}, false
	}

	return col.data[len(col.data)-1], true
}

func (col *Time) SortBy(asc bool) {
	sort.Slice(col.data, func(i, j int) bool {
		if asc {
			return col.data[i].Before(col.data[j])
		}
		return col.data[i].After(col.data[j])
	})
}
