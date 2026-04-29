package gcscb

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	grpc_hp "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/grpc/http2_protobuf"
)

type ExchangeFn_gRPC = func(data *common.ProtoData,
	client grpc_hp.EchoServiceClient, wg *sync.WaitGroup, b *testing.B)

// -----------------------------------------------------------------------------
// gRPC/HTTP2,Protobuf
// -----------------------------------------------------------------------------

func benchmarkQPS_gRPC_HTTP2_Protobuf(clientsCount int,
	dataSet [][]*common.ProtoData, b *testing.B,
) {
	benchmark_gRPC_HTTP2_Protobuf(clientsCount, 0, dataSet, grpc_hp.ExchangeQPS, b)
	b.ReportMetric(0, NsOpMetric)
	b.ReportMetric(float64(b.Elapsed()), NsMetric)
}

func benchmarkFixed_gRPC_HTTP2_Protobuf(clientsCount, n int,
	dataSet [][]*common.ProtoData, b *testing.B,
) {
	var (
		copsD      = make(chan time.Duration, n)
		exchangeFn = func(data *common.ProtoData, client grpc_hp.EchoServiceClient,
			wg *sync.WaitGroup,
			b *testing.B,
		) {
			grpc_hp.ExchangeFixed(data, client, copsD, wg, b)
		}
		N = n / clientsCount
	)
	benchmark_gRPC_HTTP2_Protobuf(clientsCount, N, dataSet, exchangeFn, b)
	b.ReportMetric(float64(N), NMetric)
	reportMetrics(copsD, b)
}

func benchmark_gRPC_HTTP2_Protobuf(clientsCount, N int,
	dataSet [][]*common.ProtoData,
	exchangeFn ExchangeFn_gRPC,
	b *testing.B,
) {
	var (
		l   net.Listener
		err error
		wgS = &sync.WaitGroup{}
	)
	l, err = grpc_hp.StartServer("127.0.0.1:0", wgS)
	if err != nil {
		b.Fatal(err)
	}
	clients, err := grpc_hp.MakeClients(l.Addr().String(), clientsCount)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	wg := &sync.WaitGroup{}
	var i int
	for i = 0; i < b.N; i++ {
		if N != 0 && i == N {
			break
		}
		wg.Add(clientsCount)
		for j := range clientsCount {
			go exchangeFn(dataSet[j][i], clients[j], wg, b)
		}
	}
	wg.Wait()
	b.StopTimer()
	if err = grpc_hp.CloseServer(l, wgS); err != nil {
		b.Fatal(err)
	}
}
