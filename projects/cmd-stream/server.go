package cs

import (
	cmdstream "github.com/cmd-stream/cmd-stream-go"
	csrv "github.com/cmd-stream/cmd-stream-go/core/srv"
	hdlr "github.com/cmd-stream/cmd-stream-go/handler"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/cmd-stream-go/transport"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
)

func MakeServer[T any](workersCount int, codec srv.Codec[T],
	invoker hdlr.Invoker[T],
) (*csrv.Server, error) {
	return cmdstream.NewServerWithInvoker(invoker, codec,
		srv.WithCore(
			csrv.WithWorkersCount(workersCount),
		),
		srv.WithTransport(
			transport.WithWriterBufSize(common.IOBufSize),
			transport.WithReaderBufSize(common.IOBufSize),
		),
	)
}
