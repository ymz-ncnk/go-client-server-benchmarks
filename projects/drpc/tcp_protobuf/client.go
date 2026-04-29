package tcpproto

import (
	"net"

	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	"storj.io/drpc/drpcconn"
	"storj.io/drpc/drpcmanager"
	"storj.io/drpc/drpcwire"
)

func MakeClients(addr string, count int) (clients []DRPCEchoServiceClient,
	err error) {
	clients = make([]DRPCEchoServiceClient, count)
	for i := range count {
		clients[i], _, err = makeClient(addr)
		if err != nil {
			return
		}
	}
	return
}

func makeClient(addr string) (client DRPCEchoServiceClient, conn *drpcconn.Conn,
	err error) {
	rawconn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}

	conn = drpcconn.NewWithOptions(rawconn, drpcconn.Options{
		Manager: drpcmanager.Options{
			WriterBufferSize: common.IOBufSize,
			Reader: drpcwire.ReaderOptions{
				MaximumBufferSize: common.IOBufSize,
			},
		},
	})

	client = NewDRPCEchoServiceClient(conn)
	return
}
