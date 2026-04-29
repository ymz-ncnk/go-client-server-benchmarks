package gcscb

import (
	"net"
	"sync"
	"testing"

	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	drpc_tp "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/drpc/tcp_protobuf"
)

type ExchangeFn_DRPC = func(data *common.ProtoData,
	client drpc_tp.DRPCEchoServiceClient, wg *sync.WaitGroup, b *testing.B)

// -----------------------------------------------------------------------------
// DRPC/TCP,Protobuf
// -----------------------------------------------------------------------------

func benchmarkQPS_DRPC_TCP_Protobuf(clientsCount int,
	dataSet [][]*common.ProtoData, b *testing.B,
) {
	benchmark_DRPC_TCP_Protobuf(clientsCount, 0, dataSet, drpc_tp.ExchangeQPS, b)
	b.ReportMetric(0, NsOpMetric)
	b.ReportMetric(float64(b.Elapsed()), NsMetric)
}

func benchmark_DRPC_TCP_Protobuf(clientsCount, N int,
	dataSet [][]*common.ProtoData,
	exchangeFn ExchangeFn_DRPC,
	b *testing.B,
) {
	var (
		l   net.Listener
		err error
		wgS = &sync.WaitGroup{}
	)
	l, err = drpc_tp.StartServer("127.0.0.1:0", wgS)
	if err != nil {
		b.Fatal(err)
	}
	clients, err := drpc_tp.MakeClients(l.Addr().String(), clientsCount)
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
	if err = drpc_tp.CloseServer(l, wgS); err != nil {
		b.Fatal(err)
	}
}
