package cs

import (
	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	sndr "github.com/cmd-stream/cmd-stream-go/sender"
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
)

func MakeSender[T any](addr string, clientsCount int, codec cln.Codec[T]) (
	sender sndr.Sender[T], err error,
) {
	return cmdstream.NewSender(addr, codec,
		sndr.WithGroup(
			grp.WithClient[T](
				cln.WithTransport(
					tspt.WithWriterBufSize(common.IOBufSize),
					tspt.WithReaderBufSize(common.IOBufSize),
				),
			),
		),
		sndr.WithClientsCount[T](clientsCount),
	)
}
