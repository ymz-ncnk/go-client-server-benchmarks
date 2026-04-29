package grpcproto

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/remote/codec/protobuf"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	kitex_utils "github.com/cloudwego/kitex/pkg/utils"
	srv "github.com/cloudwego/kitex/server"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/grpc_protobuf/kitex_gen/echo"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/grpc_protobuf/kitex_gen/echo/kitexechoservice"
)

func init() {
	klog.SetOutput(io.Discard)
}

func StartServer(addr string, wg *sync.WaitGroup) (net.Addr, srv.Server) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	var opts []srv.Option
	opts = append(opts, srv.WithListener(l))
	opts = append(opts, srv.WithPayloadCodec(protobuf.NewProtobufCodec()))
	opts = append(opts, withServerIOBufferSize())
	server := kitexechoservice.NewServer(new(echoImpl), opts...)
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Run()
	}()
	return l.Addr(), server
}

func StopServer(server srv.Server, wg *sync.WaitGroup) (err error) {
	err = server.Stop()
	wg.Wait()
	return
}

type echoImpl struct{}

func (s *echoImpl) Echo(ctx context.Context, req *echo.KitexData) (resp *echo.KitexData, err error) {
	time.Sleep(common.Delay)
	return req, nil
}

func withServerIOBufferSize() srv.Option {
	return srv.Option{F: func(o *srv.Options, di *kitex_utils.Slice) {
		err := o.Configs.(rpcinfo.MutableRPCConfig).SetIOBufferSize(common.IOBufSize)
		if err != nil {
			panic(err)
		}
	}}
}
