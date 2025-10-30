package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/core-go"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	rcvr "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_json/receiver"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_json/results"
)

type EchoCmd common.Data

func (c EchoCmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver rcvr.Receiver,
	proxy core.Proxy,
) (err error) {
	time.Sleep(common.Delay)
	_, err = proxy.Send(seq, results.EchoResult(c))
	return
}
