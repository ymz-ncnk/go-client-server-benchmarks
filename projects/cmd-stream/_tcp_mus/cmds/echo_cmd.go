package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	rcvr "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus/receiver"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus/results"
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

// func (c EchoCmd) MarshalTypedMUS(w mus.Writer) (n int, err error) {
// 	return EchoCmdDTS.Marshal(c, w)
// }

// func (c EchoCmd) SizeTypedMUS() (size int) {
// 	return EchoCmdDTS.Size(c)
// }
