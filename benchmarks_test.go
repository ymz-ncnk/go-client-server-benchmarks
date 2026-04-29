//go:generate go run ./gen/fixed2csv/ -i ./results/fixed/benchmarks.txt -d ./results/fixed
//go:generate go run ./gen/qps2csv/ -i ./results/qps/benchmarks.txt -d ./results/qps
package gcscb

import (
	"errors"
	"flag"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
)

const (
	GenSize    = 200000
	NsOpMetric = "ns/op"
	NsMetric   = "ns"
	NMetric    = "N"
)

func BenchmarkQPS(b *testing.B) {
	var (
		size = genSize()
	)

	for _, clientsCount := range []int{1, 2, 4, 8, 16} {
		b.Run(strconv.Itoa(clientsCount), func(b *testing.B) {
			dataSet := genDataSet(clientsCount, size)

			b.Run("cmd-stream_tcp_json", func(b *testing.B) {
				ds := ToCmdStreamTJ_DataSet(dataSet)
				benchmarkQPS_CmdStream_TCP_JSON(clientsCount, ds, b)
				ds = nil
				runtime.GC()
			})
			b.Run("cmd-stream_tcp_mus", func(b *testing.B) {
				ds := ToCmdStreamTM_DataSet(dataSet)
				benchmarkQPS_CmdStream_TCP_MUS(clientsCount, ds, b)
				ds = nil
				runtime.GC()
			})
			b.Run("cmd-stream_tcp_protobuf", func(b *testing.B) {
				ds := ToCmdStreamTP_DataSet(dataSet)
				benchmarkQPS_CmdStream_TCP_Protobuf(clientsCount, ds, b)
				ds = nil
				runtime.GC()
			})
			b.Run("connect_conn_protobuf", func(b *testing.B) {
				ds := common.ToProtoData(dataSet)
				benchmarkQPS_Connect_Conn_Protobuf(clientsCount, ds, b)
				ds = nil
				runtime.GC()
			})
			b.Run("drpc_tcp_protobuf", func(b *testing.B) {
				ds := common.ToProtoData(dataSet)
				benchmarkQPS_DRPC_TCP_Protobuf(clientsCount, ds, b)
				ds = nil
				runtime.GC()
			})
			b.Run("grpc_http2_protobuf", func(b *testing.B) {
				ds := common.ToProtoData(dataSet)
				benchmarkQPS_gRPC_HTTP2_Protobuf(clientsCount, ds, b)
				ds = nil
				runtime.GC()
			})
			b.Run("kitex_grpc_protobuf", func(b *testing.B) {
				ds := ToKitexDataSet(dataSet)
				benchmarkQPS_Kitex_gRPC_Protobuf(clientsCount, ds, b)
				ds = nil
				runtime.GC()
			})
		})
	}
}

func BenchmarkFixed(b *testing.B) {
	n, err := n()
	if err != nil {
		b.Fatal(err)
	}

	for _, clientsCount := range []int{1, 2, 4, 8, 16} {
		b.Run(strconv.Itoa(clientsCount), func(b *testing.B) {
			dataSet := genDataSet(clientsCount, n)

			b.Run("cmd-stream_tcp_json", func(b *testing.B) {
				ds := ToCmdStreamTJ_DataSet(dataSet)
				benchmarkFixed_CmdStream_TCP_JSON(clientsCount, n, ds, b)
				ds = nil
				runtime.GC()
			})
			b.Run("cmd-stream_tcp_mus", func(b *testing.B) {
				ds := ToCmdStreamTM_DataSet(dataSet)
				benchmarkFixed_CmdStream_TCP_MUS(clientsCount, n, ds, b)
				ds = nil
				runtime.GC()
			})
			b.Run("cmd-stream_tcp_protobuf", func(b *testing.B) {
				ds := ToCmdStreamTP_DataSet(dataSet)
				benchmarkFixed_CmdStream_TCP_Protobuf(clientsCount, n, ds, b)
				ds = nil
				runtime.GC()
			})
			b.Run("grpc_http2_protobuf", func(b *testing.B) {
				ds := common.ToProtoData(dataSet)
				benchmarkFixed_gRPC_HTTP2_Protobuf(clientsCount, n, ds, b)
				ds = nil
				runtime.GC()
			})
			b.Run("kitex_grpc_protobuf", func(b *testing.B) {
				ds := ToKitexDataSet(dataSet)
				benchmarkFixed_Kitex_gRPC_Protobuf(clientsCount, n, ds, b)
				ds = nil
				runtime.GC()
			})
		})
	}
}

// -----------------------------------------------------------------------------

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

func genDataSet(clientsCount, size int) (s [][]common.Data) {
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
