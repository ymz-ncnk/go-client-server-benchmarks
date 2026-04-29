package tcpproto

import (
	context "context"
	"sync"
	"testing"
	"time"

	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	"google.golang.org/protobuf/proto"
)

func ExchangeFixed(data *common.ProtoData, client DRPCEchoServiceClient,
	copsD chan<- time.Duration,
	wg *sync.WaitGroup,
	b *testing.B,
) {
	defer wg.Done()
	start := time.Now()
	resultData, err := exchange(data, client)
	if err != nil {
		b.Error(err)
		return
	}
	common.QueueCopD(copsD, time.Since(start))
	if !proto.Equal(data, resultData) {
		b.Error("unexpected result")
	}
}

func ExchangeQPS(data *common.ProtoData, client DRPCEchoServiceClient,
	wg *sync.WaitGroup,
	b *testing.B,
) {
	defer wg.Done()
	resultData, err := exchange(data, client)
	if err != nil {
		b.Error(err)
		return
	}
	if !proto.Equal(data, resultData) {
		b.Error("unexpected result")
	}
}

func exchange(data *common.ProtoData, client DRPCEchoServiceClient) (
	resultData *common.ProtoData, err error,
) {
	return client.Echo(context.Background(), data)
}
