package dataframe_test

import (
	"fmt"
	"go-numeric/dataframe"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestBool(t *testing.T) {
	col := dataframe.NewBool(true, false, true, false)
	col.SortBy(false)
	fmt.Println(col)
}

func TestDataFrame1(t *testing.T) {
	df1 := dataframe.New()
	df1.AddColumn("strings", dataframe.NewString("a", "b", "c", "d"))
	df1.AddColumn("ints", dataframe.NewInt(1, 2, 3, 4))

	df1.Format(os.Stdout)
	df1.SortBy("ints", false)
	df1.Format(os.Stdout)

	res := df1.FilterFunc(func(row []any) bool {
		x := row[1].(int64)
		return x%2 == 0
	})

	df1.AddColumn("ints2", dataframe.NewInt(1, 2, 3, 4))
	df1.AddColumn("ints3", dataframe.NewInt(1, 2, 3, 4))

	df1.SortBy("ints", true)

	res.Format(os.Stdout)
	df1.Format(os.Stdout)

	df2 := df1.SliceColumns("ints", "ints3")
	df2.Format(os.Stdout)

	df2.AppendRow(int64(9), int64(9))
	df2.Format(os.Stdout)
	fmt.Println(df1.Row(2))

	df1.AppendRow("really long string", 0, 0, 0)
	df1.Format(os.Stdout)

	df1.Computed(dataframe.Computed[float64]{
		"sum",
		func(row map[string]any) float64 {
			a := row["ints"].(int64)
			b := row["ints2"].(int64)
			c := row["ints3"].(int64)
			return float64(a)/float64(a+b+c) + math.Pi
		}})

	df1.Computed(dataframe.Computed[time.Time]{
		"times",
		func(row map[string]any) time.Time {
			return time.Now().Add(time.Second * time.Duration(rand.Intn(300)))
		}})

	df1.Computed(dataframe.Computed[float64]{
		"zero",
		nil,
	})

	df1.Rename("ints", "random")

	fmt.Println(df1.Len())
	df1.Format(os.Stdout)

	sumValue := df1.Column("random").(*dataframe.Int).Sum()
	fmt.Println("sum value", sumValue)

	minValue := df1.Column("random").(*dataframe.Int).Min()
	fmt.Println("min value", minValue)

	maxValue := df1.Column("random").(*dataframe.Int).Max()
	fmt.Println("max value", maxValue)

	meanValue := df1.Column("random").(*dataframe.Int).Mean()
	fmt.Println("mean value", meanValue)

	tail := df1.Column("random").(*dataframe.Int).Tail()
	fmt.Println("tail values", tail)

	results := df1.Filtered(
		dataframe.AND(
			&dataframe.GT{"random", 0, int64(1)},
			&dataframe.GT{"random", 0, int64(1)},
		),
	)

	fmt.Println(results)
	results.Format(os.Stdout)
}
