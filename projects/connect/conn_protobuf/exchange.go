package connectproto

import (
	"context"
	"sync"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/connect/conn_protobuf/connectproto/connectprotoconnect"
)

func ExchangeQPS(data *common.ProtoData, client connectprotoconnect.EchoServiceClient,
	wg *sync.WaitGroup, b *testing.B,
) {
	defer wg.Done()
	_, err := client.Echo(context.Background(), connect.NewRequest(data))
	if err != nil {
		b.Error(err)
		return
	}
}

func ExchangeFixed(data *common.ProtoData, client connectprotoconnect.EchoServiceClient,
	copsD chan time.Duration, wg *sync.WaitGroup, b *testing.B,
) {
	defer wg.Done()
	start := time.Now()
	_, err := client.Echo(context.Background(), connect.NewRequest(data))
	if err != nil {
		b.Error(err)
		return
	}
	copsD <- time.Since(start)
}
