package dataframe

import (
	"sort"
	"time"
)

func (col *Float) Min() float64 {
	if len(col.data) == 0 {
		panic("empty column")
	}
	min := col.data[0]
	for _, v := range col.data {
		if v < min {
			min = v
		}
	}
	return min
}

func (col *Float) Max() float64 {
	if len(col.data) == 0 {
		panic("empty column")
	}
	max := col.data[0]
	for _, v := range col.data {
		if v > max {
			max = v
		}
	}
	return max
}

func (col *Float) Mean() float64 {
	if len(col.data) == 0 {
		panic("empty column")
	}
	sum := float64(0)
	for _, v := range col.data {
		sum += v
	}
	return float64(sum) / float64(len(col.data))
}

func (col *Float) Median() float64 {
	if len(col.data) == 0 {
		panic("empty column")
	}
	sorted := append([]float64{}, col.data...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})
	n := len(sorted)
	if n%2 == 0 {
		return float64(sorted[n/2-1]+sorted[n/2]) / 2
	}
	return float64(sorted[n/2])
}

func (col *Float) Sum() float64 {
	if len(col.data) == 0 {
		return 0
	}

	var sum float64
	for _, v := range col.data {
		sum += v
	}

	return sum
}

// Analytical

func (col *Int) Min() int64 {
	if len(col.data) == 0 {
		panic("empty column")
	}
	min := col.data[0]
	for _, v := range col.data {
		if v < min {
			min = v
		}
	}
	return min
}

func (col *Int) Max() int64 {
	if len(col.data) == 0 {
		panic("empty column")
	}
	max := col.data[0]
	for _, v := range col.data {
		if v > max {
			max = v
		}
	}
	return max
}

func (col *Int) Mean() float64 {
	if len(col.data) == 0 {
		return 0.0
	}

	sum := int64(0)
	for _, v := range col.data {
		sum += v
	}
	return float64(sum) / float64(len(col.data))
}

func (col *Int) Median() float64 {
	if len(col.data) == 0 {
		panic("empty column")
	}

	sorted := append([]int64{}, col.data...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	n := len(sorted)
	if n%2 == 0 {
		return float64(sorted[n/2-1]+sorted[n/2]) / 2
	}

	return float64(sorted[n/2])
}

func (col *Int) Sum() int64 {
	if len(col.data) == 0 {
		return 0
	}

	var sum int64
	for _, v := range col.data {
		sum += v
	}

	return sum
}

// Unique

func (col *Int) Unique() IColumn {
	set := make(map[int64]struct{})
	uniqueData := []int64{}

	for _, v := range col.data {
		if _, exists := set[v]; !exists {
			set[v] = struct{}{}
			uniqueData = append(uniqueData, v)
		}
	}
	return &Int{data: uniqueData}
}

func (col *Float) Unique() IColumn {
	set := make(map[float64]struct{})
	uniqueData := []float64{}

	for _, v := range col.data {
		if _, exists := set[v]; !exists {
			set[v] = struct{}{}
			uniqueData = append(uniqueData, v)
		}
	}

	return &Float{data: uniqueData}
}

func (col *Bool) Unique() *Bool {
	set := map[bool]struct{}{}
	uniqueData := []bool{}

	for _, v := range col.data {
		if _, exists := set[v]; !exists {
			set[v] = struct{}{}
			uniqueData = append(uniqueData, v)
		}
	}

	return &Bool{data: uniqueData}
}

func (col *String) Unique() IColumn {
	set := make(map[string]struct{})
	uniqueData := []string{}

	for _, v := range col.data {
		if _, exists := set[v]; !exists {
			set[v] = struct{}{}
			uniqueData = append(uniqueData, v)
		}
	}
	return &String{data: uniqueData}
}

func (col *Time) Unique() IColumn {
	set := make(map[time.Time]struct{})
	uniqueData := []time.Time{}

	for _, v := range col.data {
		if _, exists := set[v]; !exists {
			set[v] = struct{}{}
			uniqueData = append(uniqueData, v)
		}
	}
	return &Time{data: uniqueData}
}
