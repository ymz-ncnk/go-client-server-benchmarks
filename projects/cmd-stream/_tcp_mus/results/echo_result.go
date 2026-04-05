package results

import (
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
)

type EchoResult common.Data

func (r EchoResult) LastOne() bool {
	return true
}

// func (c EchoResult) MarshalTypedMUS(w muss.Writer) (n int, err error) {
// 	return EchoResultDTS.Marshal(c, w)
// }

// func (c EchoResult) SizeTypedMUS() (size int) {
// 	return EchoResultDTS.Size(c)
// }
