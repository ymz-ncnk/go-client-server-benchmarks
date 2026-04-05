package tcpmus

import (
	"context"
	"sync"
	"testing"
	"time"

	sndr "github.com/cmd-stream/cmd-stream-go/sender"

	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus/cmds"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus/receiver"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus/results"
)

func ExchangeQPS(cmd cmds.EchoCmd, sender sndr.Sender[receiver.Receiver],
	wg *sync.WaitGroup,
	b *testing.B,
) {
	defer wg.Done()
	result, err := sender.Send(context.Background(), cmd)
	if err != nil {
		b.Error(err)
		return
	}
	if !common.EqualData(common.Data(cmd), common.Data(result.(results.EchoResult))) {
		b.Error("unexpected result")
	}
}

func ExchangeFixed(cmd cmds.EchoCmd, sender sndr.Sender[receiver.Receiver],
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
	if !common.EqualData(common.Data(cmd), common.Data(result.(results.EchoResult))) {
		b.Error("unexpected result")
	}
}
