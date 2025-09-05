//go:generate go run ./gen/fixed2csv/ -i ./results/fixed/benchmarks.txt -d ./results/fixed
//go:generate go run ./gen/qps2csv/ -i ./results/qps/benchmarks.txt -d ./results/qps
package gcscb

import (
	"errors"
	"flag"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	srv "github.com/cmd-stream/cmd-stream-go/server"
	cdc "github.com/cmd-stream/codec-mus-stream-go"
	sndr "github.com/cmd-stream/sender-go"
	"github.com/montanaflynn/stats"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	cs "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream"
	cstm "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus"
	cstm_cmds "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus/cmds"
	cstm_rcvr "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus/receiver"
	cstm_rslts "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus/results"
	cstp "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_protobuf"
	cstp_cmds "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_protobuf/cmds"
	cstp_rcvr "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_protobuf/receiver"
	cstp_rslts "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_protobuf/results"
	ghp "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/grpc/http2_protobuf"
	kthp "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/ttheader_protobuf"
	kthp_echo "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/ttheader_protobuf/kitex_gen/echo"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/ttheader_protobuf/kitex_gen/echo/kitexechoservice"
)

const (
	GenSize    = 200000
	NsOpMetric = "ns/op"
	NsMetric   = "ns"
	NMetric    = "N"
)

type ExchangeFn_gRPC = func(data *common.ProtoData,
	client ghp.EchoServiceClient, wg *sync.WaitGroup, b *testing.B)

type ExchangeFn_Kitex = func(data *kthp_echo.KitexData,
	client kitexechoservice.Client, wg *sync.WaitGroup, b *testing.B)

type ExchangeFn_CmdStream = func(cmd cstm_cmds.EchoCmd,
	sender sndr.Sender[cstm_rcvr.Receiver], wg *sync.WaitGroup, b *testing.B)

type ExchangeFn_CmdStream_Protobuf = func(cmd cstp_cmds.EchoCmd,
	sender sndr.Sender[cstp_rcvr.Receiver], wg *sync.WaitGroup, b *testing.B)

func BenchmarkQPS(b *testing.B) {
	var (
		dataSet     = generateDataSet(16, genSize())
		ghpDataSet  = common.ToProtoData(dataSet)
		kthpDataSet = ToKthpDataSet(dataSet)
		cstmDataSet = ToCstmDataSet(dataSet)
		cstpDataSet = ToCstpDataSet(dataSet)
	)

	b.Run("1", func(b *testing.B) {
		clientsCount := 1

		b.Run("grpc_http2_protobuf", func(b *testing.B) {
			benchmarkQPS_gRPC_HTTP2_Protobuf(clientsCount, ghpDataSet, b)
		})
		b.Run("kitex_ttheader_protobuf", func(b *testing.B) {
			benchmarkQPS_Kitex_TTHeader_Protobuf(clientsCount, kthpDataSet, b)
		})
		b.Run("cmd-stream_tcp_mus", func(b *testing.B) {
			benchmarkQPS_CmdStream_TCP_MUS(clientsCount, cstmDataSet, b)
		})
		b.Run("cmd-stream_tcp_protobuf", func(b *testing.B) {
			benchmarkQPS_CmdStream_TCP_Protobuf(clientsCount, cstpDataSet, b)
		})
	})

	b.Run("2", func(b *testing.B) {
		clientsCount := 2

		b.Run("grpc_http2_protobuf", func(b *testing.B) {
			benchmarkQPS_gRPC_HTTP2_Protobuf(clientsCount, ghpDataSet, b)
		})
		b.Run("kitex_ttheader_protobuf", func(b *testing.B) {
			benchmarkQPS_Kitex_TTHeader_Protobuf(clientsCount, kthpDataSet, b)
		})
		b.Run("cmd-stream_tcp_mus", func(b *testing.B) {
			benchmarkQPS_CmdStream_TCP_MUS(clientsCount, cstmDataSet, b)
		})
		b.Run("cmd-stream_tcp_protobuf", func(b *testing.B) {
			benchmarkQPS_CmdStream_TCP_Protobuf(clientsCount, cstpDataSet, b)
		})
	})

	b.Run("4", func(b *testing.B) {
		clientsCount := 4

		b.Run("grpc_http2_protobuf", func(b *testing.B) {
			benchmarkQPS_gRPC_HTTP2_Protobuf(clientsCount, ghpDataSet, b)
		})
		b.Run("kitex_ttheader_protobuf", func(b *testing.B) {
			benchmarkQPS_Kitex_TTHeader_Protobuf(clientsCount, kthpDataSet, b)
		})
		b.Run("cmd-stream_tcp_mus", func(b *testing.B) {
			benchmarkQPS_CmdStream_TCP_MUS(clientsCount, cstmDataSet, b)
		})
		b.Run("cmd-stream_tcp_protobuf", func(b *testing.B) {
			benchmarkQPS_CmdStream_TCP_Protobuf(clientsCount, cstpDataSet, b)
		})
	})

	b.Run("8", func(b *testing.B) {
		clientsCount := 8

		b.Run("grpc_http2_protobuf", func(b *testing.B) {
			benchmarkQPS_gRPC_HTTP2_Protobuf(clientsCount, ghpDataSet, b)
		})
		b.Run("kitex_ttheader_protobuf", func(b *testing.B) {
			benchmarkQPS_Kitex_TTHeader_Protobuf(clientsCount, kthpDataSet, b)
		})
		b.Run("cmd-stream_tcp_mus", func(b *testing.B) {
			benchmarkQPS_CmdStream_TCP_MUS(clientsCount, cstmDataSet, b)
		})
		b.Run("cmd-stream_tcp_protobuf", func(b *testing.B) {
			benchmarkQPS_CmdStream_TCP_Protobuf(clientsCount, cstpDataSet, b)
		})
	})

	b.Run("16", func(b *testing.B) {
		clientsCount := 16

		b.Run("grpc_http2_protobuf", func(b *testing.B) {
			benchmarkQPS_gRPC_HTTP2_Protobuf(clientsCount, ghpDataSet, b)
		})
		b.Run("kitex_ttheader_protobuf", func(b *testing.B) {
			benchmarkQPS_Kitex_TTHeader_Protobuf(clientsCount, kthpDataSet, b)
		})
		b.Run("cmd-stream_tcp_mus", func(b *testing.B) {
			benchmarkQPS_CmdStream_TCP_MUS(clientsCount, cstmDataSet, b)
		})
		b.Run("cmd-stream_tcp_protobuf", func(b *testing.B) {
			benchmarkQPS_CmdStream_TCP_Protobuf(clientsCount, cstpDataSet, b)
		})
	})
}

// nethttp_json is commented out because it takes too long to run.
func BenchmarkFixed(b *testing.B) {
	n, err := n()
	if err != nil {
		b.Fatal(err)
	}

	var (
		dataSet     = generateDataSet(16, n)
		ghpDataSet  = common.ToProtoData(dataSet)
		kthpDataSet = ToKthpDataSet(dataSet)
		cstmDataSet = ToCstmDataSet(dataSet)
		cstpDataSet = ToCstpDataSet(dataSet)
	)

	b.Run("1", func(b *testing.B) {
		clientsCount := 1

		b.Run("grpc_http2_protobuf", func(b *testing.B) {
			benchmarkFixed_gRPC_HTTP2_Protobuf(clientsCount, n, ghpDataSet, b)
		})
		b.Run("kitex_ttheader_protobuf", func(b *testing.B) {
			benchmarkFixed_Kitex_TTHeader_Protobuf(clientsCount, n, kthpDataSet, b)
		})
		b.Run("cmd-stream_tcp_mus", func(b *testing.B) {
			benchmarkFixed_CmdStream_TCP_MUS(clientsCount, n, cstmDataSet, b)
		})
		b.Run("cmd-stream_tcp_protobuf", func(b *testing.B) {
			benchmarkFixed_CmdStream_TCP_Protobuf(clientsCount, n, cstpDataSet, b)
		})
	})

	b.Run("2", func(b *testing.B) {
		clientsCount := 2

		b.Run("grpc_http2_protobuf", func(b *testing.B) {
			benchmarkFixed_gRPC_HTTP2_Protobuf(clientsCount, n, ghpDataSet, b)
		})
		b.Run("kitex_ttheader_protobuf", func(b *testing.B) {
			benchmarkFixed_Kitex_TTHeader_Protobuf(clientsCount, n, kthpDataSet, b)
		})
		b.Run("cmd-stream_tcp_mus", func(b *testing.B) {
			benchmarkFixed_CmdStream_TCP_MUS(clientsCount, n, cstmDataSet, b)
		})
		b.Run("cmd-stream_tcp_protobuf", func(b *testing.B) {
			benchmarkFixed_CmdStream_TCP_Protobuf(clientsCount, n, cstpDataSet, b)
		})
	})

	b.Run("4", func(b *testing.B) {
		clientsCount := 4

		b.Run("grpc_http2_protobuf", func(b *testing.B) {
			benchmarkFixed_gRPC_HTTP2_Protobuf(clientsCount, n, ghpDataSet, b)
		})
		b.Run("kitex_ttheader_protobuf", func(b *testing.B) {
			benchmarkFixed_Kitex_TTHeader_Protobuf(clientsCount, n, kthpDataSet, b)
		})
		b.Run("cmd-stream_tcp_mus", func(b *testing.B) {
			benchmarkFixed_CmdStream_TCP_MUS(clientsCount, n, cstmDataSet, b)
		})
		b.Run("cmd-stream_tcp_protobuf", func(b *testing.B) {
			benchmarkFixed_CmdStream_TCP_Protobuf(clientsCount, n, cstpDataSet, b)
		})
	})

	b.Run("8", func(b *testing.B) {
		clientsCount := 8

		b.Run("grpc_http2_protobuf", func(b *testing.B) {
			benchmarkFixed_gRPC_HTTP2_Protobuf(clientsCount, n, ghpDataSet, b)
		})
		b.Run("kitex_ttheader_protobuf", func(b *testing.B) {
			benchmarkFixed_Kitex_TTHeader_Protobuf(clientsCount, n, kthpDataSet, b)
		})
		b.Run("cmd-stream_tcp_mus", func(b *testing.B) {
			benchmarkFixed_CmdStream_TCP_MUS(clientsCount, n, cstmDataSet, b)
		})
		b.Run("cmd-stream_tcp_protobuf", func(b *testing.B) {
			benchmarkFixed_CmdStream_TCP_Protobuf(clientsCount, n, cstpDataSet, b)
		})
	})

	b.Run("16", func(b *testing.B) {
		clientsCount := 16

		b.Run("grpc_http2_protobuf", func(b *testing.B) {
			benchmarkFixed_gRPC_HTTP2_Protobuf(clientsCount, n, ghpDataSet, b)
		})
		b.Run("kitex_ttheader_protobuf", func(b *testing.B) {
			benchmarkFixed_Kitex_TTHeader_Protobuf(clientsCount, n, kthpDataSet, b)
		})
		b.Run("cmd-stream_tcp_mus", func(b *testing.B) {
			benchmarkFixed_CmdStream_TCP_MUS(clientsCount, n, cstmDataSet, b)
		})
		b.Run("cmd-stream_tcp_protobuf", func(b *testing.B) {
			benchmarkFixed_CmdStream_TCP_Protobuf(clientsCount, n, cstpDataSet, b)
		})
	})
}

// -----------------------------------------------------------------------------
// gRPC/HTTP2,Protobuf
// -----------------------------------------------------------------------------

func benchmarkQPS_gRPC_HTTP2_Protobuf(clientsCount int,
	dataSet [][]*common.ProtoData, b *testing.B,
) {
	benchmark_gRPC_HTTP2_Protobuf(clientsCount, 0, dataSet, ghp.ExchangeQPS, b)
	b.ReportMetric(0, NsOpMetric)
	b.ReportMetric(float64(b.Elapsed()), NsMetric)
}

func benchmarkFixed_gRPC_HTTP2_Protobuf(clientsCount, n int,
	dataSet [][]*common.ProtoData, b *testing.B,
) {
	var (
		copsD      = make(chan time.Duration, n)
		exchangeFn = func(data *common.ProtoData, client ghp.EchoServiceClient,
			wg *sync.WaitGroup,
			b *testing.B,
		) {
			ghp.ExchangeFixed(data, client, copsD, wg, b)
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
		addr = "127.0.0.1:9001"
		l    net.Listener
		err  error
		wgS  = &sync.WaitGroup{}
	)
	l, err = ghp.StartServer(addr, wgS)
	if err != nil {
		b.Fatal(err)
	}
	clients, err := ghp.MakeClients(addr, clientsCount)
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
	if err = ghp.CloseServer(l, wgS); err != nil {
		b.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// Kitex/TTHeader,Protobuf
// -----------------------------------------------------------------------------

func benchmarkQPS_Kitex_TTHeader_Protobuf(clientsCount int,
	dataSet [][]*kthp_echo.KitexData,
	b *testing.B,
) {
	benchmark_Kitex_TTHeader_Protobuf(clientsCount, 0, dataSet, kthp.ExchangeQPS, b)
	b.ReportMetric(0, NsOpMetric)
	b.ReportMetric(float64(b.Elapsed()), NsMetric)
}

func benchmarkFixed_Kitex_TTHeader_Protobuf(clientsCount, n int,
	dataSet [][]*kthp_echo.KitexData,
	b *testing.B,
) {
	var (
		copsD                       = make(chan time.Duration, n)
		exchangeFn ExchangeFn_Kitex = func(data *kthp_echo.KitexData,
			client kitexechoservice.Client,
			wg *sync.WaitGroup,
			b *testing.B,
		) {
			kthp.ExchangeFixed(data, client, copsD, wg, b)
		}
		N = n / clientsCount
	)
	benchmark_Kitex_TTHeader_Protobuf(clientsCount, N, dataSet,
		exchangeFn, b)
	b.ReportMetric(float64(N), NMetric)
	reportMetrics(copsD, b)
}

func benchmark_Kitex_TTHeader_Protobuf(clientsCount, N int,
	dataSet [][]*kthp_echo.KitexData,
	exchangeFn ExchangeFn_Kitex,
	b *testing.B,
) {
	var (
		addr = "127.0.0.1:9002"
		wgS  = &sync.WaitGroup{}
	)
	server := kthp.StartServer(addr, wgS)
	clients, err := kthp.MakeClients(addr, clientsCount)
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
	if err = kthp.StopServer(server, wgS); err != nil {
		b.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// cmd-stream/TCP,MUS
// -----------------------------------------------------------------------------

func benchmarkQPS_CmdStream_TCP_MUS(clientsCount int,
	dataSet [][]cstm_cmds.EchoCmd,
	b *testing.B,
) {
	benchmark_CmdStream_TCP_MUS(clientsCount, 0, dataSet, cstm.ExchangeQPS, b)
	b.ReportMetric(0, NsOpMetric)
	b.ReportMetric(float64(b.Elapsed()), NsMetric)
}

func benchmarkFixed_CmdStream_TCP_MUS(clientsCount, n int,
	dataSet [][]cstm_cmds.EchoCmd,
	b *testing.B,
) {
	var (
		copsD      = make(chan time.Duration, n)
		exchangeFn = func(cmd cstm_cmds.EchoCmd, sender sndr.Sender[cstm_rcvr.Receiver],
			wg *sync.WaitGroup,
			b *testing.B,
		) {
			cstm.ExchangeFixed(cmd, sender, copsD, wg, b)
		}
		N = n / clientsCount
	)
	benchmark_CmdStream_TCP_MUS(clientsCount, N, dataSet, exchangeFn, b)
	b.ReportMetric(float64(N), NMetric)
	reportMetrics(copsD, b)
}

func benchmark_CmdStream_TCP_MUS(clientsCount, N int,
	dataSet [][]cstm_cmds.EchoCmd,
	exchangFn ExchangeFn_CmdStream,
	b *testing.B,
) {
	var (
		addr        = "127.0.0.1:9003"
		invoker     = srv.NewInvoker(cstm_rcvr.Receiver{})
		serverCodec = cdc.NewServerCodec(cstm_cmds.CmdMUS, cstm_rslts.ResultMUS)
		clientCodec = cdc.NewClientCodec(cstm_cmds.CmdMUS, cstm_rslts.ResultMUS)
		wgS         = &sync.WaitGroup{}
	)
	server := cs.MakeServer(clientsCount, serverCodec, invoker)
	wgS.Add(1)
	go func() {
		server.ListenAndServe(addr)
		wgS.Done()
	}()
	time.Sleep(100 * time.Millisecond)

	sender, err := cs.MakeSender(addr, clientsCount, clientCodec)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	wg := &sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		if N != 0 && i == N {
			break
		}
		wg.Add(clientsCount)
		for j := range clientsCount {
			go exchangFn(dataSet[j][i], sender, wg, b)
		}
	}
	wg.Wait()
	b.StopTimer()

	if err = sender.CloseAndWait(time.Second); err != nil {
		b.Fatal(err)
	}
	if err = server.Close(); err != nil {
		b.Fatal(err)
	}
	wgS.Wait()
}

// -----------------------------------------------------------------------------
// cmd-stream/TCP,Protobuf
// -----------------------------------------------------------------------------

// If you are looking for an example of using cmd-stream/Protobuf also check
// https://github.com/cmd-stream/cmd-stream-examples-go/tree/main/standard_protobuf.

func benchmarkQPS_CmdStream_TCP_Protobuf(clientsCount int,
	dataSet [][]cstp_cmds.EchoCmd,
	b *testing.B,
) {
	benchmark_CmdStream_TCP_Protobuf(clientsCount, 0, dataSet, cstp.ExchangeQPS, b)
	b.ReportMetric(0, NsOpMetric)
	b.ReportMetric(float64(b.Elapsed()), NsMetric)
}

func benchmarkFixed_CmdStream_TCP_Protobuf(clientsCount, n int,
	dataSet [][]cstp_cmds.EchoCmd,
	b *testing.B,
) {
	var (
		copsD      = make(chan time.Duration, n)
		exchangeFn = func(cmd cstp_cmds.EchoCmd, sender sndr.Sender[cstp_rcvr.Receiver],
			wg *sync.WaitGroup,
			b *testing.B,
		) {
			cstp.ExchangeFixed(cmd, sender, copsD, wg, b)
		}
		N = n / clientsCount
	)
	benchmark_CmdStream_TCP_Protobuf(clientsCount, N, dataSet, exchangeFn, b)
	b.ReportMetric(float64(N), NMetric)
	reportMetrics(copsD, b)
}

func benchmark_CmdStream_TCP_Protobuf(clientsCount, N int,
	dataSet [][]cstp_cmds.EchoCmd,
	exchangFn ExchangeFn_CmdStream_Protobuf,
	b *testing.B,
) {
	var (
		addr        = "127.0.0.1:9004"
		invoker     = srv.NewInvoker(cstp_rcvr.Receiver{})
		serverCodec = cdc.NewTypedServerCodec(cstp_cmds.EchoCmdMUS, cstp_rslts.EchoResultMUS)
		clientCodec = cdc.NewTypedClientCodec(cstp_cmds.EchoCmdMUS, cstp_rslts.EchoResultMUS)
		wgS         = &sync.WaitGroup{}
	)
	server := cs.MakeServer(clientsCount, serverCodec, invoker)
	wgS.Add(1)
	go func() {
		server.ListenAndServe(addr)
		wgS.Done()
	}()
	time.Sleep(100 * time.Millisecond)
	sender, err := cs.MakeSender(addr, clientsCount, clientCodec)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	wg := &sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		if N != 0 && i == N {
			break
		}
		wg.Add(clientsCount)
		for j := range clientsCount {
			go exchangFn(dataSet[j][i], sender, wg, b)
		}
	}
	wg.Wait()
	b.StopTimer()

	if err = sender.CloseAndWait(time.Second); err != nil {
		b.Fatal(err)
	}
	if err = server.Close(); err != nil {
		b.Fatal(err)
	}
	wgS.Wait()
}

func n() (n int, err error) {
	var (
		f         = flag.Lookup("test.benchtime")
		benchtime = f.Value.String()
	)
	if !strings.HasSuffix(benchtime, "x") {
		err = errors.New("you should specify -benchtime with x suffix")
		return
	}
	return strconv.Atoi(benchtime[:len(benchtime)-1])
}

func genSize() int {
	val := os.Getenv("GEN_SIZE")
	if val == "" {
		return GenSize
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return n
}

func generateDataSet(clientsCount, size int) (s [][]common.Data) {
	s = make([][]common.Data, clientsCount)
	for i := range clientsCount {
		s[i] = make([]common.Data, size)
		for j := range size {
			s[i][j] = common.NewRandomData()
		}
	}
	return
}

func reportMetrics(copsD chan time.Duration, b *testing.B) {
	copsDArr := makeCopsDArr(copsD)

	mean, _ := stats.Mean(copsDArr)
	med, _ := stats.Median(copsDArr)
	max, _ := stats.Max(copsDArr)
	min, _ := stats.Min(copsDArr)
	p99, _ := stats.Percentile(copsDArr, 99.9)

	b.ReportMetric(0, "ns/op")
	b.ReportMetric(float64(b.Elapsed()), "ns")
	b.ReportMetric(mean, "ns/cop")
	b.ReportMetric(med, "ns/med")
	b.ReportMetric(max, "ns/max")
	b.ReportMetric(min, "ns/min")
	b.ReportMetric(p99, "ns/p99")
}

func makeCopsDArr(copsD chan time.Duration) (copsDArr []float64) {
	close(copsD)
	copsDArr = []float64{}
	for spent := range copsD {
		copsDArr = append(copsDArr, float64(spent))
	}
	return
}
