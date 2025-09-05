package cs

import (
	cmdstream "github.com/cmd-stream/cmd-stream-go"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	csrv "github.com/cmd-stream/core-go/server"
	"github.com/cmd-stream/handler-go"
	"github.com/cmd-stream/transport-go"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
)

func MakeServer[T any](workersCount int, codec srv.Codec[T],
	invoker handler.Invoker[T],
) *csrv.Server {
	return cmdstream.MakeServer(codec, invoker,
		srv.WithCore(
			csrv.WithWorkersCount(workersCount),
		),
		srv.WithTransport(
			transport.WithWriterBufSize(common.IOBufSize),
			transport.WithReaderBufSize(common.IOBufSize),
		),
	)
}
