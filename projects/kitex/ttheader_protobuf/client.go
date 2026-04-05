package ttheaderproto

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	kitex_utils "github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/transport"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/ttheader_protobuf/kitex_gen/echo/kitexechoservice"

	"github.com/cloudwego/kitex/pkg/remote/codec/protobuf"
)

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
	opts = append(opts, client.WithMetaHandler(transmeta.ClientTTHeaderHandler))
	opts = append(opts, client.WithPayloadCodec(protobuf.NewProtobufCodec()))
	opts = append(opts, client.WithTransportProtocol(transport.TTHeader))
	opts = append(opts, client.WithMuxConnection(1))
	opts = append(opts, withClientIOBufferSize())
	return kitexechoservice.NewClient("echo", opts...)
}

func withClientIOBufferSize() client.Option {
	return client.Option{F: func(o *client.Options, di *kitex_utils.Slice) {
		err := o.Configs.(rpcinfo.MutableRPCConfig).SetIOBufferSize(common.IOBufSize)
		if err != nil {
			panic(err)
		}
	}}
}
