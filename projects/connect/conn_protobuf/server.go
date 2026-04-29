package connectproto

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/connect/conn_protobuf/connectproto/connectprotoconnect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func StartServer(addr string, wg *sync.WaitGroup) (string, *http.Server) {
	mux := http.NewServeMux()
	path, handler := connectprotoconnect.NewEchoServiceHandler(&echoServer{})
	mux.Handle(path, handler)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		ReadHeaderTimeout: time.Second,
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	return l.Addr().String(), server
}

func StopServer(srv *http.Server, wg *sync.WaitGroup) error {
	err := srv.Shutdown(context.Background())
	wg.Wait()
	return err
}

type echoServer struct{}

func (s *echoServer) Echo(
	ctx context.Context,
	req *connect.Request[common.ProtoData],
) (*connect.Response[common.ProtoData], error) {
	time.Sleep(common.Delay)
	return connect.NewResponse(req.Msg), nil
}
