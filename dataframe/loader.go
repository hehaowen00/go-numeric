package dataframe

import "io"

func LoadStruct(data ...any) *DataFrame {
	df := New()
	return df
}

func LoadCSV(rdr io.Reader) *DataFrame {
	df := New()
	return df
}

func LoadJSON(rdr io.Reader) *DataFrame {
	df := New()
	return df
}

func (df *DataFrame) ToCSV() {
}

func (df *DataFrame) ToJSON() {

}
