package gcscb

import (
	"sync"
	"testing"
	"time"

	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	kitex_gp "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/grpc_protobuf"
	kitex_gp_echo "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/grpc_protobuf/kitex_gen/echo"
	kitex_gp_service "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/grpc_protobuf/kitex_gen/echo/kitexechoservice"
)

type ExchangeFn_Kitex_gRPC = func(data *kitex_gp_echo.KitexData,
	client kitex_gp_service.Client, wg *sync.WaitGroup, b *testing.B)

func ToKitexDataSet(dataSet [][]common.Data) [][]*kitex_gp_echo.KitexData {
	kitexDataSet := make([][]*kitex_gp_echo.KitexData, len(dataSet))
	for i := range len(dataSet) {
		kitexDataSet[i] = make([]*kitex_gp_echo.KitexData, len(dataSet[i]))
		for j := range len(dataSet[i]) {
			kitexDataSet[i][j] = &kitex_gp_echo.KitexData{
				Bool:    dataSet[i][j].Bool,
				Int64:   dataSet[i][j].Int64,
				String:  dataSet[i][j].String,
				Float64: dataSet[i][j].Float64,
			}
		}
	}
	return kitexDataSet
}

// -----------------------------------------------------------------------------
// Kitex/gRPC,Protobuf
// -----------------------------------------------------------------------------

func benchmarkQPS_Kitex_gRPC_Protobuf(clientsCount int,
	dataSet [][]*kitex_gp_echo.KitexData,
	b *testing.B,
) {
	benchmark_Kitex_gRPC_Protobuf(clientsCount, 0, dataSet, kitex_gp.ExchangeQPS, b)
	b.ReportMetric(0, NsOpMetric)
	b.ReportMetric(float64(b.Elapsed()), NsMetric)
}

func benchmarkFixed_Kitex_gRPC_Protobuf(clientsCount, n int,
	dataSet [][]*kitex_gp_echo.KitexData,
	b *testing.B,
) {
	var (
		copsD                            = make(chan time.Duration, n)
		exchangeFn ExchangeFn_Kitex_gRPC = func(data *kitex_gp_echo.KitexData,
			client kitex_gp_service.Client,
			wg *sync.WaitGroup,
			b *testing.B,
		) {
			kitex_gp.ExchangeFixed(data, client, copsD, wg, b)
		}
		N = n / clientsCount
	)
	benchmark_Kitex_gRPC_Protobuf(clientsCount, N, dataSet, exchangeFn, b)
	b.ReportMetric(float64(N), NMetric)
	reportMetrics(copsD, b)
}

func benchmark_Kitex_gRPC_Protobuf(clientsCount, N int,
	dataSet [][]*kitex_gp_echo.KitexData,
	exchangeFn ExchangeFn_Kitex_gRPC,
	b *testing.B,
) {
	var (
		wgS = &sync.WaitGroup{}
	)
	actualAddr, server := kitex_gp.StartServer("127.0.0.1:0", wgS)
	time.Sleep(500 * time.Millisecond)
	clients, err := kitex_gp.MakeClients(actualAddr.String(), clientsCount)
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
		for j := range len(clients) {
			go exchangeFn(dataSet[j][i], clients[j], wg, b)
		}
	}
	wg.Wait()
	b.StopTimer()
	if err = kitex_gp.StopServer(server, wgS); err != nil {
		b.Fatal(err)
	}
}
