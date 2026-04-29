package tcpproto

import (
	context "context"
	"net"
	"sync"
	"time"

	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	"storj.io/drpc/drpcmanager"
	"storj.io/drpc/drpcmux"
	"storj.io/drpc/drpcserver"
	"storj.io/drpc/drpcwire"
)

func StartServer(addr string, wg *sync.WaitGroup) (listener net.Listener, err error) {
	listener, err = net.Listen("tcp", addr)
	if err != nil {
		return
	}
	wg.Add(1)
	server := makeServer()
	go func() {
		defer wg.Done()
		server.Serve(context.Background(), listener)
	}()
	return
}

func CloseServer(l net.Listener, wg *sync.WaitGroup) (err error) {
	if err = l.Close(); err != nil {
		return
	}
	wg.Wait()
	return
}

func makeServer() *drpcserver.Server {
	m := drpcmux.New()

	err := DRPCRegisterEchoService(m, &echoServer{})
	if err != nil {
		panic(err)
	}

	return drpcserver.NewWithOptions(m, drpcserver.Options{
		Manager: drpcmanager.Options{
			WriterBufferSize: common.IOBufSize,
			Reader: drpcwire.ReaderOptions{
				MaximumBufferSize: common.IOBufSize,
			},
		},
	})
}

type echoServer struct {
	DRPCEchoServiceServer
}

func (server *echoServer) Echo(ctx context.Context, data *common.ProtoData) (
	*common.ProtoData, error) {
	time.Sleep(common.Delay)
	return &common.ProtoData{
		Bool:    data.Bool,
		Int64:   data.Int64,
		String_: data.String_,
		Float64: data.Float64,
	}, nil
}
