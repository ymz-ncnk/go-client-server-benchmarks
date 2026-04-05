package tcproto

import (
	"context"
	"sync"
	"testing"
	"time"

	sndr "github.com/cmd-stream/cmd-stream-go/sender"
	"google.golang.org/protobuf/proto"

	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
)

func ExchangeQPS(cmd *EchoCmd, sender sndr.Sender[Receiver],
	wg *sync.WaitGroup,
	b *testing.B,
) {
	defer wg.Done()
	result, err := sender.Send(context.Background(), cmd)
	if err != nil {
		b.Error(err)
		return
	}
	if !proto.Equal(cmd.ProtoData, result.(*EchoResult).ProtoData) {
		b.Error("unexpected result")
	}
}

func ExchangeFixed(cmd *EchoCmd, sender sndr.Sender[Receiver],
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

	if !proto.Equal(cmd.ProtoData, result.(*EchoResult).ProtoData) {
		b.Error("unexpected result")
	}
}
