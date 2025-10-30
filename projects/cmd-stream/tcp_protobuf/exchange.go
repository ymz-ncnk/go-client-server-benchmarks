package tcpproto

import (
	"context"
	"sync"
	"testing"
	"time"

	sndr "github.com/cmd-stream/sender-go"
	"google.golang.org/protobuf/proto"

	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_protobuf/cmds"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_protobuf/receiver"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_protobuf/results"
)

func ExchangeQPS(cmd *cmds.EchoCmd, sender sndr.Sender[receiver.Receiver],
	wg *sync.WaitGroup,
	b *testing.B,
) {
	defer wg.Done()
	result, err := sender.Send(context.Background(), cmd)
	if err != nil {
		b.Error(err)
		return
	}
	if !proto.Equal(cmd.ProtoData, result.(*results.EchoResult).ProtoData) {
		b.Error("unexpected result")
	}
}

func ExchangeFixed(cmd *cmds.EchoCmd, sender sndr.Sender[receiver.Receiver],
	copsD chan<- time.Duration,
	wg *sync.WaitGroup,
	b *testing.B,
) {
	defer wg.Done()
	start := time.Now()
	result, err := sender.Send(context.Background(), cmd)
	if err != nil {
		b.Error(err)
		return
	}
	common.QueueCopD(copsD, time.Since(start))

	if !proto.Equal(cmd.ProtoData, result.(*results.EchoResult).ProtoData) {
		b.Error("unexpected result")
	}
}
