package dataframe

import "time"

type IApplicable interface {
	Apply(df *DataFrame)
}

func (df *DataFrame) Apply(xs ...IApplicable) {
	for _, el := range xs {
		el.Apply(df)
	}
}

type Or struct {
	filters []filter
}

func (or Or) check(df *DataFrame, i int) bool {
	for _, f := range or.filters {
		if f.check(df, i) {
			return true
		}
	}

	return false
}

func OR(filters ...filter) filter {
	return &Or{
		filters: filters,
	}
}

type And struct {
	filters []filter
}

func (and And) check(df *DataFrame, i int) bool {
	for _, f := range and.filters {
		if !f.check(df, i) {
			return false
		}
	}

	return true
}

func AND(filters ...filter) filter {
	return &And{
		filters: filters,
	}
}

type EQ struct {
	Column string
	Op     int
	Value  any
}

func (eq *EQ) check(df *DataFrame, i int) bool {
	idx := df.index[eq.Column]
	col := df.data[idx]

	switch col.(type) {
	case *Int:
		x := col.Index(i).(int64)
		y := eq.Value.(int64)
		return x == y
	case *Float:
		x := col.Index(i).(float64)
		y := eq.Value.(float64)
		return x == y
	case *Bool:
		x := col.Index(i).(bool)
		y := eq.Value.(bool)
		return x == y
	case *String:
		x := col.Index(i).(string)
		y := eq.Value.(string)
		return x == y
	case *Time:
		x := col.Index(i).(time.Time)
		y := eq.Value.(time.Time)
		return x == y
	default:
		return false
	}
}

type NEQ struct {
	Column string
	Op     int
	Value  any
}

func (neq *NEQ) check(df *DataFrame, i int) bool {
	idx := df.index[neq.Column]
	col := df.data[idx]

	switch col.(type) {
	case *Int:
		x := col.Index(i).(int64)
		y := neq.Value.(int64)
		return x != y
	case *Float:
		x := col.Index(i).(float64)
		y := neq.Value.(float64)
		return x != y
	case *Bool:
		x := col.Index(i).(bool)
		y := neq.Value.(bool)
		return x != y
	case *String:
		x := col.Index(i).(string)
		y := neq.Value.(string)
		return x != y
	case *Time:
		x := col.Index(i).(time.Time)
		y := neq.Value.(time.Time)
		return x != y
	default:
		return false
	}
}

type LT struct {
	Column string
	Op     int
	Value  any
}

func (lt *LT) check(df *DataFrame, i int) bool {
	idx := df.index[lt.Column]
	col := df.data[idx]

	switch col.(type) {
	case *Int:
		x := col.Index(i).(int64)
		y := lt.Value.(int64)
		return x < y
	case *Float:
		x := col.Index(i).(float64)
		y := lt.Value.(float64)
		return x < y
	case *Bool:
		x := col.Index(i).(bool)
		y := lt.Value.(bool)
		return !x && y
	case *String:
		x := col.Index(i).(string)
		y := lt.Value.(string)
		return x < y
	case *Time:
		x := col.Index(i).(time.Time)
		y := lt.Value.(time.Time)
		return x.Before(y)
	default:
		return false
	}
}

type GT struct {
	Column string
	Op     int
	Value  any
}

func (gt *GT) check(df *DataFrame, i int) bool {
	idx := df.index[gt.Column]
	col := df.data[idx]

	switch col.(type) {
	case *Int:
		x := col.Index(i).(int64)
		y := gt.Value.(int64)
		return x > y
	case *Float:
		x := col.Index(i).(float64)
		y := gt.Value.(float64)
		return x > y
	case *Bool:
		x := col.Index(i).(bool)
		y := gt.Value.(bool)
		return x && !y
	case *String:
		x := col.Index(i).(string)
		y := gt.Value.(string)
		return x > y
	case *Time:
		x := col.Index(i).(time.Time)
		y := gt.Value.(time.Time)
		return x.After(y)
	default:
		return false
	}
}

type LTE struct {
	Column string
	Op     int
	Value  any
}

func (lte *LTE) check(df *DataFrame, i int) bool {
	idx := df.index[lte.Column]
	col := df.data[idx]

	switch col.(type) {
	case *Int:
		x := col.Index(i).(int64)
		y := lte.Value.(int64)
		return x <= y
	case *Float:
		x := col.Index(i).(float64)
		y := lte.Value.(float64)
		return x <= y
	case *Bool:
		return true
	case *String:
		x := col.Index(i).(string)
		y := lte.Value.(string)
		return x <= y
	case *Time:
		x := col.Index(i).(time.Time)
		y := lte.Value.(time.Time)
		return x.Before(y) || x == y
	default:
		return false
	}
}

type GTE struct {
	Column string
	Op     int
	Value  any
}

func (gte *GTE) check(df *DataFrame, i int) bool {
	idx := df.index[gte.Column]
	col := df.data[idx]

	switch col.(type) {
	case *Int:
		x := col.Index(i).(int64)
		y := gte.Value.(int64)
		return x >= y
	case *Float:
		x := col.Index(i).(float64)
		y := gte.Value.(float64)
		return x >= y
	case *Bool:
		return true
	case *String:
		x := col.Index(i).(string)
		y := gte.Value.(string)
		return x >= y
	case *Time:
		x := col.Index(i).(time.Time)
		y := gte.Value.(time.Time)
		return x.After(y) || x == y
	default:
		return false
	}
}

type IN[T any] struct {
	Column string
	Values []T
}
