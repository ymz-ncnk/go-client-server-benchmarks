package ttheaderproto

import (
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/ttheader_protobuf/kitex_gen/echo"
)

func EqualData(d1, d2 *echo.KitexData) bool {
	return d1.Bool == d2.Bool && d1.Int64 == d2.Int64 &&
		d1.String_ == d2.String_ &&
		d1.Float64 == d2.Float64
}
