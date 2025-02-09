package dataframe

import (
	"fmt"
	"io"
	"reflect"
	"slices"
	"sort"
	"strings"
	"time"
)

type Computed[T any] struct {
	Name string
	Func func(row map[string]any) T
}

type DataFrame struct {
	data     []IColumn
	headers  []string
	rowCount int
	index    map[string]int
}

func New() *DataFrame {
	return &DataFrame{
		index: make(map[string]int),
	}
}

func (df *DataFrame) Len() int {
	return df.rowCount
}

func (df *DataFrame) Headers() []string {
	return slices.Clone(df.headers)
}

func (df *DataFrame) NumColumns() int {
	return len(df.headers)
}

func (df *DataFrame) AddColumn(name string, col IColumn) {
	length := col.Len()

	if length > df.rowCount {
		for _, c := range df.data {
			c.Extend(length)
		}
		df.rowCount = length
	} else if length < df.rowCount {
		col.Extend(df.rowCount)
	}

	df.headers = append(df.headers, name)
	df.data = append(df.data, col)
	df.index[name] = len(df.headers) - 1
}

func (df *DataFrame) DeleteColumn(name string) IColumn {
	if len(df.headers) == 0 {
		return nil
	}

	idx, ok := df.index[name]
	if !ok {
		return nil
	}

	col := df.data[idx]
	df.data = append(df.data[:idx], df.data[idx+1:]...)
	return col
}

func (df *DataFrame) DeleteRow(index int) {
	if index < 0 || index >= df.rowCount {
		panic("Index out of range")
	}

	for _, col := range df.data {
		col.DeleteRow(index)
	}

	df.rowCount--
}

func (df *DataFrame) Format(wr io.Writer) {
	if df.rowCount == 0 || len(df.headers) == 0 {
		fmt.Fprintln(wr, "empty dataframe")
		return
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("> DataFrame [%dx%d]\n", len(df.headers), df.rowCount))

	colWidths := make([]int, len(df.headers))
	rowNumWidth := len(fmt.Sprintf("%d", df.rowCount)) + 1 // Width for row numbers (e.g., "10:")

	for i, header := range df.headers {
		colWidths[i] = len(header)
	}

	for i := 0; i < df.rowCount; i++ {
		for j, col := range df.data {
			fmt.Println(col, i)
			valStr := formatValue(col.Index(i))
			if len(valStr) > colWidths[j] {
				colWidths[j] = len(valStr)
			}
		}
	}

	builder.WriteString(fmt.Sprintf("%-*s ", rowNumWidth, ""))
	for i, header := range df.headers {
		builder.WriteString(fmt.Sprintf("%-*s ", colWidths[i], header))
	}
	builder.WriteString("\n")

	for i := 0; i < df.rowCount; i++ {
		builder.WriteString(fmt.Sprintf("%-*d: ", rowNumWidth-1, i+1))
		for j, col := range df.data {
			valStr := formatValue(col.Index(i))
			builder.WriteString(fmt.Sprintf("%-*s ", colWidths[j], valStr))
		}
		builder.WriteString("\n")
	}

	wr.Write([]byte(builder.String()))
}
func formatValue(value any) string {
	switch v := value.(type) {
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%.3f", v)
	case string:
		return v
	case bool:
		return fmt.Sprintf("%t", v)
	case time.Time:
		return v.Format("2006-01-02 15:04:05")
	default:
		return "?"
	}
}

func (df *DataFrame) FilterFunc(predicate func(row []any) bool) *DataFrame {
	newDF := New()

	newDF.headers = append([]string{}, df.headers...)
	newDF.index = make(map[string]int)
	for k, v := range df.index {
		newDF.index[k] = v
	}

	newDF.data = make([]IColumn, len(df.data))
	for i, col := range df.data {
		newDF.data[i] = col.New()
	}

	for i := 0; i < df.rowCount; i++ {
		row := make([]any, len(df.data))
		for j, col := range df.data {
			row[j] = col.Index(i)
		}

		if predicate(row) {
			for j, col := range newDF.data {
				switch c := col.(type) {
				case *Int:
					c.data = append(c.data, row[j].(int64))
				case *Float:
					c.data = append(c.data, row[j].(float64))
				case *String:
					c.data = append(c.data, row[j].(string))
				case *Bool:
					c.data = append(c.data, row[j].(bool))
				case *Time:
					c.data = append(c.data, row[j].(time.Time))
				}
			}
			newDF.rowCount++
		}
	}

	return newDF
}

type filter interface {
	check(df *DataFrame, i int) bool
}

func (df *DataFrame) Filtered(f filter) *DataFrame {
	newDF := New()

	newDF.headers = append([]string{}, df.headers...)
	newDF.index = make(map[string]int)
	for k, v := range df.index {
		newDF.index[k] = v
	}

	newDF.data = make([]IColumn, len(df.data))
	for i, col := range df.data {
		newDF.data[i] = col.New()
	}

outer:
	for i := 0; i < df.rowCount; i++ {
		row := make([]any, len(df.data))
		for j, col := range df.data {
			row[j] = col.Index(i)
		}

		if !f.check(df, i) {
			continue outer
		}

		for j, col := range newDF.data {
			switch c := col.(type) {
			case *Int:
				c.data = append(c.data, row[j].(int64))
			case *Float:
				c.data = append(c.data, row[j].(float64))
			case *String:
				c.data = append(c.data, row[j].(string))
			case *Bool:
				c.data = append(c.data, row[j].(bool))
			case *Time:
				c.data = append(c.data, row[j].(time.Time))
			}

			newDF.rowCount++
		}
	}

	return newDF
}
func (df *DataFrame) SortBy(columnName string, ascending bool) {
	colIndex, exists := df.index[columnName]
	if !exists {
		panic("column not found")
	}

	sortOrder := make([]int, df.rowCount)
	for i := range sortOrder {
		sortOrder[i] = i
	}

	switch c := df.data[colIndex].(type) {
	case *Int:
		sort.SliceStable(sortOrder, func(i, j int) bool {
			if ascending {
				return c.data[sortOrder[i]] < c.data[sortOrder[j]]
			}
			return c.data[sortOrder[i]] > c.data[sortOrder[j]]
		})
	case *Float:
		sort.SliceStable(sortOrder, func(i, j int) bool {
			if ascending {
				return c.data[sortOrder[i]] < c.data[sortOrder[j]]
			}
			return c.data[sortOrder[i]] > c.data[sortOrder[j]]
		})
	case *String:
		sort.SliceStable(sortOrder, func(i, j int) bool {
			if ascending {
				return c.data[sortOrder[i]] < c.data[sortOrder[j]]
			}
			return c.data[sortOrder[i]] > c.data[sortOrder[j]]
		})
	case *Time:
		sort.SliceStable(sortOrder, func(i, j int) bool {
			if ascending {
				return c.data[sortOrder[i]].Before(c.data[sortOrder[j]])
			}
			return c.data[sortOrder[i]].After(c.data[sortOrder[j]])
		})
	case *Bool:
		sort.SliceStable(sortOrder, func(i, j int) bool {
			if ascending {
				return !c.data[sortOrder[i]] && (c.data[sortOrder[j]])
			}
			return c.data[sortOrder[i]] && !(c.data[sortOrder[j]])
		})
	default:
	}

	for i, col := range df.data {
		newCol := col.Clone()
		for j, idx := range sortOrder {
			newCol.Set(j, col.Index(idx))
		}
		df.data[i] = newCol
	}
}

func (df *DataFrame) SliceColumns(columns ...string) *DataFrame {
	frame := New()
	for _, c := range columns {
		colIndex, exists := df.index[c]
		if !exists {
			panic("column not found")
		}

		frame.AddColumn(c, df.data[colIndex].Clone())
	}

	return frame
}

func (df *DataFrame) IndexColumn(index int) IColumn {
	return df.data[index]
}

func (df *DataFrame) Column(column string) IColumn {
	colIndex, exists := df.index[column]
	if !exists {
		panic("column not found")
	}

	return df.data[colIndex]
}

func (df *DataFrame) Row(index int) []any {
	if index > df.rowCount {
		return nil
	}

	row := []any{}

	for i := range df.headers {
		switch c := df.data[i].(type) {
		case *Int:
			row = append(row, c.data[index])
		case *Float:
			row = append(row, c.data[index])
		case *String:
			row = append(row, c.data[index])
		case *Bool:
			row = append(row, c.data[index])
		case *Time:
			row = append(row, c.data[index])
		}
	}

	return row
}

func (df *DataFrame) AppendRow(row ...any) {
	if len(row) > len(df.headers) {
		row = row[:len(df.headers)]
	}

	row = convert(row)

	for i := range df.headers {
		if row[i] == nil {
			continue
		}

		switch c := df.data[i].(type) {
		case *Int:
			v := row[i].(int64)
			c.Append(v)
		case *Float:
			v := row[i].(float64)
			c.Append(v)
		case *Bool:
			v := row[i].(bool)
			c.Append(v)
		case *String:
			v := row[i].(string)
			c.Append(v)
		case *Time:
			v := row[i].(time.Time)
			c.Append(v)
		}
	}

	df.rowCount++
	for i := range df.headers {
		df.data[i].Extend(df.rowCount)
	}
}

func (df *DataFrame) Rename(oldName, newName string) {
	idx, ok := df.index[oldName]
	if !ok {
		return
	}

	df.index[newName] = idx
	df.headers[idx] = newName
	delete(df.index, oldName)
}

func (df *DataFrame) Computed(
	compute any,
	newCol ...IColumn,
) {
	cols := map[string]IColumn{}
	for i := range df.headers {
		cols[df.headers[i]] = df.data[df.index[df.headers[i]]]
	}

	var col IColumn
	var name string

	switch c := compute.(type) {
	case Computed[int64]:
		name = c.Name
		r := NewInt()
		computeHelper(df, cols, &c, func(v int64) {
			r.Append(v)
		})
		col = r
	case Computed[float64]:
		name = c.Name
		r := NewFloat()
		computeHelper(df, cols, &c, func(v float64) {
			r.Append(v)
		})
		col = r
	case Computed[bool]:
		name = c.Name
		r := NewBool()
		computeHelper(df, cols, &c, func(v bool) {
			r.Append(v)
		})
		col = r
	case Computed[string]:
		name = c.Name
		r := NewString()
		computeHelper(df, cols, &c, func(v string) {
			r.Append(v)
		})
		col = r
	case Computed[time.Time]:
		name = c.Name
		r := NewTime()
		computeHelper(df, cols, &c, func(v time.Time) {
			r.Append(v)
		})
		col = r
	default:
		panic(fmt.Errorf("unknown column - %v", reflect.TypeOf(c)))
	}

	col.Extend(df.rowCount)
	df.index[name] = len(df.data)
	df.data = append(df.data, col)
	df.headers = append(df.headers, name)
}

func computeHelper[T any](
	df *DataFrame,
	cols map[string]IColumn,
	compute *Computed[T],
	appendFunc func(T),
) {
	if compute.Func == nil {
		return
	}

	for i := range df.Len() {
		row := map[string]any{}
		for k, v := range cols {
			row[k] = v.Index(i)
		}
		appendFunc(compute.Func(row))
	}
}

func convert(row []any) []any {
	res := make([]any, 0, len(row))

	for _, v := range row {
		switch c := v.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			res = append(res, int64(toInt64(c)))
		case float32, float64:
			switch r := c.(type) {
			case float32:
				res = append(res, float64(r))
			case float64:
				res = append(res, r)
			}
		case bool, string, time.Time:
			res = append(res, c)
		}
	}

	return res
}
func toInt64(v any) int64 {
	switch c := v.(type) {
	case int:
		return int64(c)
	case int8:
		return int64(c)
	case int16:
		return int64(c)
	case int32:
		return int64(c)
	case int64:
		return c
	case uint:
		return int64(c)
	case uint8:
		return int64(c)
	case uint16:
		return int64(c)
	case uint32:
		return int64(c)
	case uint64:
		return int64(c)
	default:
		return 0
	}
}
