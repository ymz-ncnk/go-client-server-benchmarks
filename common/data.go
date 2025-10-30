package common

import (
	"math/rand"

	"github.com/brianvoe/gofakeit"
)

const (
	Charset         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	MaxStringLength = 1007
)

type Data struct {
	Bool    bool
	Int64   int64
	String  string
	Float64 float64
}

func NewRandomData() Data {
	return Data{
		Bool:    gofakeit.Bool(),
		Int64:   gofakeit.Int64(),
		String:  RandomString(),
		Float64: gofakeit.Float64(),
	}
}

func RandomString() (s string) {
	b := make([]byte, rand.Intn(MaxStringLength))
	for i := range b {
		b[i] = Charset[rand.Intn(len(Charset))]
	}
	return string(b)
}

func ToProtoData(dataSet [][]Data) [][]*ProtoData {
	s := make([][]*ProtoData, len(dataSet))
	for i := range len(dataSet) {
		s[i] = make([]*ProtoData, len(dataSet[i]))
		for j := range len(dataSet[i]) {
			s[i][j] = &ProtoData{
				Bool:    dataSet[i][j].Bool,
				Int64:   dataSet[i][j].Int64,
				String_: dataSet[i][j].String,
				Float64: dataSet[i][j].Float64,
			}
		}
	}
	return s
}

func EqualData(d1, d2 Data) bool {
	return d1.Bool == d2.Bool && d1.Int64 == d2.Int64 &&
		d1.String == d2.String &&
		d1.Float64 == d2.Float64
}

// func EqualProtoData(d1, d2 *ProtoData) bool {
// 	return d1.Bool == d2.Bool && d1.Int64 == d2.Int64 &&
// 		d1.String_ == d2.String_ &&
// 		d1.Float64 == d2.Float64
// }
