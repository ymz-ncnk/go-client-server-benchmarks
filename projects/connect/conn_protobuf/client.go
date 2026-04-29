package connectproto

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"

	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/connect/conn_protobuf/connectproto/connectprotoconnect"
	"golang.org/x/net/http2"
)

func MakeClients(addr string, count int) (clients []connectprotoconnect.EchoServiceClient, err error) {
	clients = make([]connectprotoconnect.EchoServiceClient, count)
	for i := range count {
		httpClient := &http.Client{
			Transport: &http2.Transport{
				AllowHTTP: true,
				DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
					return net.Dial(network, addr)
				},
				ReadIdleTimeout: common.Delay * 2,
			},
		}
		clients[i] = connectprotoconnect.NewEchoServiceClient(httpClient,
			"http://"+addr,
		)
	}
	return clients, nil
}
