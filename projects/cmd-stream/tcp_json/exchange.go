package tcpmus

import (
	"context"
	"sync"
	"testing"
	"time"

	sndr "github.com/cmd-stream/sender-go"

	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_json/cmds"
	rcvr "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_json/receiver"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_json/results"
)

func ExchangeQPS(cmd cmds.EchoCmd, sender sndr.Sender[rcvr.Receiver],
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

func ExchangeFixed(cmd cmds.EchoCmd, sender sndr.Sender[rcvr.Receiver],
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
