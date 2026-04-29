package gcscb

import (
	"net"
	"reflect"
	"sync"
	"testing"
	"time"

	sndr "github.com/cmd-stream/cmd-stream-go/sender"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	cdcjson "github.com/cmd-stream/codec-json-go"
	cdcmus "github.com/cmd-stream/codec-mus-stream-go"
	cdcproto "github.com/cmd-stream/codec-protobuf-go"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	cmdstream "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream"
	cmdstream_tj "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_json"
	cmdstream_tm "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus"
	cmdstream_tp "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_protobuf"
)

type ExchangeFn_CmdStream_JSON = func(cmd cmdstream_tj.EchoCmd,
	sender sndr.Sender[cmdstream_tj.Receiver], wg *sync.WaitGroup, b *testing.B)

type ExchangeFn_CmdStream_MUS = func(cmd cmdstream_tm.EchoCmd,
	sender sndr.Sender[cmdstream_tm.Receiver], wg *sync.WaitGroup, b *testing.B)

type ExchangeFn_CmdStream_Protobuf = func(cmd *cmdstream_tp.EchoCmd,
	sender sndr.Sender[cmdstream_tp.Receiver], wg *sync.WaitGroup, b *testing.B)

// -----------------------------------------------------------------------------

func ToCmdStreamTM_DataSet(dataSet [][]common.Data) [][]cmdstream_tm.EchoCmd {
	cmdDataSet := make([][]cmdstream_tm.EchoCmd, len(dataSet))
	for i := range len(dataSet) {
		cmdDataSet[i] = make([]cmdstream_tm.EchoCmd, len(dataSet[i]))
		for j := range len(dataSet[i]) {
			cmdDataSet[i][j] = cmdstream_tm.EchoCmd(dataSet[i][j])
		}
	}
	return cmdDataSet
}

func ToCmdStreamTP_DataSet(dataSet [][]common.Data) [][]*cmdstream_tp.EchoCmd {
	var (
		cmdDataSet   = make([][]*cmdstream_tp.EchoCmd, len(dataSet))
		protoDataSet = common.ToProtoData(dataSet)
	)
	for i := range len(protoDataSet) {
		cmdDataSet[i] = make([]*cmdstream_tp.EchoCmd, len(protoDataSet[i]))
		for j := range len(protoDataSet[i]) {
			cmdDataSet[i][j] = &cmdstream_tp.EchoCmd{ProtoData: protoDataSet[i][j]}
		}
	}
	return cmdDataSet
}

func ToCmdStreamTJ_DataSet(dataSet [][]common.Data) [][]cmdstream_tj.EchoCmd {
	cmdDataSet := make([][]cmdstream_tj.EchoCmd, len(dataSet))
	for i := range len(dataSet) {
		cmdDataSet[i] = make([]cmdstream_tj.EchoCmd, len(dataSet[i]))
		for j := range len(dataSet[i]) {
			cmdDataSet[i][j] = cmdstream_tj.EchoCmd(dataSet[i][j])
		}
	}
	return cmdDataSet
}

// -----------------------------------------------------------------------------
// cmd-stream/TCP,JSON
// -----------------------------------------------------------------------------

func benchmarkQPS_CmdStream_TCP_JSON(clientsCount int,
	dataSet [][]cmdstream_tj.EchoCmd,
	b *testing.B,
) {
	benchmark_CmdStream_TCP_JSON(clientsCount, 0, dataSet, cmdstream_tj.ExchangeQPS, b)
	b.ReportMetric(0, NsOpMetric)
	b.ReportMetric(float64(b.Elapsed()), NsMetric)
}

func benchmarkFixed_CmdStream_TCP_JSON(clientsCount, n int,
	dataSet [][]cmdstream_tj.EchoCmd,
	b *testing.B,
) {
	var (
		copsD      = make(chan time.Duration, n)
		exchangeFn = func(cmd cmdstream_tj.EchoCmd, sender sndr.Sender[cmdstream_tj.Receiver],
			wg *sync.WaitGroup,
			b *testing.B,
		) {
			cmdstream_tj.ExchangeFixed(cmd, sender, copsD, wg, b)
		}
		N = n / clientsCount
	)
	benchmark_CmdStream_TCP_JSON(clientsCount, N, dataSet, exchangeFn, b)
	b.ReportMetric(float64(N), NMetric)
	reportMetrics(copsD, b)
}

func benchmark_CmdStream_TCP_JSON(clientsCount, N int,
	dataSet [][]cmdstream_tj.EchoCmd,
	exchangFn ExchangeFn_CmdStream_JSON,
	b *testing.B,
) {
	var (
		invoker  = srv.NewInvoker(cmdstream_tj.Receiver{})
		cmdTypes = []reflect.Type{
			reflect.TypeFor[cmdstream_tj.EchoCmd](),
		}
		resultTypes = []reflect.Type{
			reflect.TypeFor[cmdstream_tj.EchoResult](),
		}
		serverCodec = cdcjson.NewServerCodec[cmdstream_tj.Receiver](cmdTypes,
			resultTypes)
		clientCodec = cdcjson.NewClientCodec[cmdstream_tj.Receiver](cmdTypes,
			resultTypes)
		wgS = &sync.WaitGroup{}
	)
	server, err := cmdstream.MakeServer(clientsCount, serverCodec, invoker)
	if err != nil {
		b.Fatal(err)
	}
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}
	wgS.Add(1)
	go func() {
		server.Serve(l.(*net.TCPListener))
		wgS.Done()
	}()
	time.Sleep(100 * time.Millisecond)

	sender, err := cmdstream.MakeSender(l.Addr().String(), clientsCount, clientCodec)
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
// cmd-stream/TCP,MUS
// -----------------------------------------------------------------------------

func benchmarkQPS_CmdStream_TCP_MUS(clientsCount int,
	dataSet [][]cmdstream_tm.EchoCmd,
	b *testing.B,
) {
	benchmark_CmdStream_TCP_MUS(clientsCount, 0, dataSet, cmdstream_tm.ExchangeQPS, b)
	b.ReportMetric(0, NsOpMetric)
	b.ReportMetric(float64(b.Elapsed()), NsMetric)
}

func benchmarkFixed_CmdStream_TCP_MUS(clientsCount, n int,
	dataSet [][]cmdstream_tm.EchoCmd,
	b *testing.B,
) {
	var (
		copsD      = make(chan time.Duration, n)
		exchangeFn = func(cmd cmdstream_tm.EchoCmd, sender sndr.Sender[cmdstream_tm.Receiver],
			wg *sync.WaitGroup,
			b *testing.B,
		) {
			cmdstream_tm.ExchangeFixed(cmd, sender, copsD, wg, b)
		}
		N = n / clientsCount
	)
	benchmark_CmdStream_TCP_MUS(clientsCount, N, dataSet, exchangeFn, b)
	b.ReportMetric(float64(N), NMetric)
	reportMetrics(copsD, b)
}

func benchmark_CmdStream_TCP_MUS(clientsCount, N int,
	dataSet [][]cmdstream_tm.EchoCmd,
	exchangFn ExchangeFn_CmdStream_MUS,
	b *testing.B,
) {
	var (
		invoker     = srv.NewInvoker(cmdstream_tm.Receiver{})
		serverCodec = cdcmus.NewServerCodec(cmdstream_tm.CmdMUS, cmdstream_tm.ResultMUS)
		clientCodec = cdcmus.NewClientCodec(cmdstream_tm.CmdMUS, cmdstream_tm.ResultMUS)
		wgS         = &sync.WaitGroup{}
	)
	server, err := cmdstream.MakeServer(clientsCount, serverCodec, invoker)
	if err != nil {
		b.Fatal(err)
	}
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}
	wgS.Add(1)
	go func() {
		server.Serve(l.(*net.TCPListener))
		wgS.Done()
	}()
	time.Sleep(100 * time.Millisecond)

	sender, err := cmdstream.MakeSender(l.Addr().String(), clientsCount, clientCodec)
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

func benchmarkQPS_CmdStream_TCP_Protobuf(clientsCount int,
	dataSet [][]*cmdstream_tp.EchoCmd,
	b *testing.B,
) {
	benchmark_CmdStream_TCP_Protobuf(clientsCount, 0, dataSet, cmdstream_tp.ExchangeQPS, b)
	b.ReportMetric(0, NsOpMetric)
	b.ReportMetric(float64(b.Elapsed()), NsMetric)
}

func benchmarkFixed_CmdStream_TCP_Protobuf(clientsCount, n int,
	dataSet [][]*cmdstream_tp.EchoCmd,
	b *testing.B,
) {
	var (
		copsD      = make(chan time.Duration, n)
		exchangeFn = func(cmd *cmdstream_tp.EchoCmd, sender sndr.Sender[cmdstream_tp.Receiver],
			wg *sync.WaitGroup,
			b *testing.B,
		) {
			cmdstream_tp.ExchangeFixed(cmd, sender, copsD, wg, b)
		}
		N = n / clientsCount
	)
	benchmark_CmdStream_TCP_Protobuf(clientsCount, N, dataSet, exchangeFn, b)
	b.ReportMetric(float64(N), NMetric)
	reportMetrics(copsD, b)
}

func benchmark_CmdStream_TCP_Protobuf(clientsCount, N int,
	dataSet [][]*cmdstream_tp.EchoCmd,
	exchangFn ExchangeFn_CmdStream_Protobuf,
	b *testing.B,
) {
	var (
		invoker  = srv.NewInvoker(cmdstream_tp.Receiver{})
		cmdTypes = []reflect.Type{
			reflect.TypeFor[*cmdstream_tp.EchoCmd](),
		}
		resultTypes = []reflect.Type{
			reflect.TypeFor[*cmdstream_tp.EchoResult](),
		}
		serverCodec = cdcproto.NewServerCodec[cmdstream_tp.Receiver](cmdTypes,
			resultTypes)
		clientCodec = cdcproto.NewClientCodec[cmdstream_tp.Receiver](cmdTypes,
			resultTypes)
		wgS = &sync.WaitGroup{}
	)
	server, err := cmdstream.MakeServer(clientsCount, serverCodec, invoker)
	if err != nil {
		b.Fatal(err)
	}
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}
	wgS.Add(1)
	go func() {
		server.Serve(l.(*net.TCPListener))
		wgS.Done()
	}()
	time.Sleep(100 * time.Millisecond)
	sender, err := cmdstream.MakeSender(l.Addr().String(), clientsCount, clientCodec)
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
