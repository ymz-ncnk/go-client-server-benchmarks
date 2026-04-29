package grpcproto

import (
	"io"
	"time"

	cln "github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	kitex_utils "github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/transport"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/grpc_protobuf/kitex_gen/echo/kitexechoservice"
)

func init() {
	klog.SetOutput(io.Discard)
}

func MakeClients(addr string, count int) (clients []kitexechoservice.Client,
	err error) {
	clients = make([]kitexechoservice.Client, count)
	for i := range count {
		clients[i], err = makeClient(addr)
		if err != nil {
			return
		}
	}
	return
}

func makeClient(addr string) (kitexechoservice.Client, error) {
	var opts []cln.Option
	opts = append(opts, cln.WithHostPorts(addr))
	opts = append(opts, cln.WithTransportProtocol(transport.GRPC))
	opts = append(opts, cln.WithConnectTimeout(5*time.Second))
	opts = append(opts, withClientIOBufferSize())
	return kitexechoservice.NewClient("KitexEchoService", opts...)
}

func withClientIOBufferSize() cln.Option {
	return cln.Option{F: func(o *cln.Options, di *kitex_utils.Slice) {
		err := o.Configs.(rpcinfo.MutableRPCConfig).SetIOBufferSize(common.IOBufSize)
		if err != nil {
			panic(err)
		}
	}}
}
