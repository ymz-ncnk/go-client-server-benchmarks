package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	rcvr "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_protobuf/receiver"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_protobuf/results"
)

func (c *EchoCmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver rcvr.Receiver,
	proxy core.Proxy,
) (err error) {
	time.Sleep(common.Delay)
	_, err = proxy.Send(seq, &results.EchoResult{ProtoData: c.ProtoData})
	return
}
