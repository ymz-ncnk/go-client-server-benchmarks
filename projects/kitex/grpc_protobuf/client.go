package grpcproto

import (
	"io"

	"github.com/cloudwego/kitex/client"
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
	var opts []client.Option
	opts = append(opts, client.WithHostPorts(addr))
	opts = append(opts, client.WithTransportProtocol(transport.GRPC))
	opts = append(opts, withClientIOBufferSize())
	return kitexechoservice.NewClient("KitexEchoService", opts...)
}

func withClientIOBufferSize() client.Option {
	return client.Option{F: func(o *client.Options, di *kitex_utils.Slice) {
		err := o.Configs.(rpcinfo.MutableRPCConfig).SetIOBufferSize(common.IOBufSize)
		if err != nil {
			panic(err)
		}
	}}
}
