package gcscb

import (
	"sync"
	"testing"
	"time"

	connect_cp "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/connect/conn_protobuf"

	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	connect_cp_service "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/connect/conn_protobuf/connectproto/connectprotoconnect"
)

type ExchangeFn_Connect = func(data *common.ProtoData,
	client connect_cp_service.EchoServiceClient, wg *sync.WaitGroup, b *testing.B)

// -----------------------------------------------------------------------------
// Connect/Conn,Protobuf
// -----------------------------------------------------------------------------

func benchmarkQPS_Connect_Conn_Protobuf(clientsCount int,
	dataSet [][]*common.ProtoData, b *testing.B,
) {
	benchmark_Connect_Conn_Protobuf(clientsCount, 0, dataSet, connect_cp.ExchangeQPS, b)
	b.ReportMetric(0, NsOpMetric)
	b.ReportMetric(float64(b.Elapsed()), NsMetric)
}

func benchmark_Connect_Conn_Protobuf(clientsCount, N int,
	dataSet [][]*common.ProtoData,
	exchangeFn ExchangeFn_Connect,
	b *testing.B,
) {
	var (
		wgS = &sync.WaitGroup{}
	)
	actualAddr, server := connect_cp.StartServer("127.0.0.1:0", wgS)
	time.Sleep(100 * time.Millisecond)
	clients, err := connect_cp.MakeClients(actualAddr, clientsCount)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	wg := &sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		if N != 0 && i == N {
			break
		}
		wg.Add(len(clients))
		for j := 0; j < len(clients); j++ {
			go exchangeFn(dataSet[j][i], clients[j], wg, b)
		}
	}
	wg.Wait()
	b.StopTimer()
	if err = connect_cp.StopServer(server, wgS); err != nil {
		b.Fatal(err)
	}
}
