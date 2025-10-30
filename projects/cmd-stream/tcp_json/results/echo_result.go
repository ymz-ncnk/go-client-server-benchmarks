package results

import (
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
)

type EchoResult common.Data

func (r EchoResult) LastOne() bool {
	return true
}
