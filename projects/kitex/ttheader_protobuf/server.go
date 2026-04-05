package ttheaderproto

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	kitex_utils "github.com/cloudwego/kitex/pkg/utils"
	srv "github.com/cloudwego/kitex/server"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/ttheader_protobuf/kitex_gen/echo"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/ttheader_protobuf/kitex_gen/echo/kitexechoservice"

	"github.com/cloudwego/kitex/pkg/remote/codec/protobuf"
)

func init() {
	klog.SetLevel(klog.LevelError)
}

func StartServer(addr string, wg *sync.WaitGroup) (server srv.Server) {
	wg.Add(1)
	a, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}
	var opts []srv.Option
	opts = append(opts, srv.WithServiceAddr(a))
	opts = append(opts, srv.WithMetaHandler(transmeta.ServerTTHeaderHandler))
	opts = append(opts, srv.WithPayloadCodec(protobuf.NewProtobufCodec()))
	opts = append(opts, srv.WithMuxTransport())
	opts = append(opts, withServerIOBufferSize())
	server = kitexechoservice.NewServer(new(echoImpl), opts...)
	go func() {
		defer wg.Done()
		server.Run()
	}()
	return
}

func StopServer(server srv.Server, wg *sync.WaitGroup) (err error) {
	if err = server.Stop(); err != nil {
		return
	}
	wg.Wait()
	return
}

type echoImpl struct{}

func (s *echoImpl) Echo(ctx context.Context, req *echo.KitexData) (resp *echo.KitexData, err error) {
	time.Sleep(common.Delay)
	return &echo.KitexData{
		Bool:    req.Bool,
		Int64:   req.Int64,
		String:  req.String,
		Float64: req.Float64,
	}, nil
}

func withServerIOBufferSize() srv.Option {
	return srv.Option{F: func(o *srv.Options, di *kitex_utils.Slice) {
		err := o.Configs.(rpcinfo.MutableRPCConfig).SetIOBufferSize(common.IOBufSize)
		if err != nil {
			panic(err)
		}
	}}
}
