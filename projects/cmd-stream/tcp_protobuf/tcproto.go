package tcproto

import (
	"context"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
)

func (c *EchoCmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver Receiver, proxy core.Proxy,
) (err error) {
	time.Sleep(common.Delay)
	_, err = proxy.Send(seq, &EchoResult{ProtoData: c.ProtoData})
	return
}

func (r *EchoResult) LastOne() bool {
	return true
}
