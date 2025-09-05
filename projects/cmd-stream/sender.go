package cs

import (
	cln "github.com/cmd-stream/cmd-stream-go/client"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	sndr "github.com/cmd-stream/sender-go"
	"github.com/cmd-stream/transport-go"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
)

func MakeSender[T any](addr string, clientsCount int, codec cln.Codec[T]) (
	sender sndr.Sender[T], err error,
) {
	return sndr.Make(addr, codec,
		sndr.WithGroup(
			grp.WithClient[T](
				cln.WithTransport(
					transport.WithWriterBufSize(common.IOBufSize),
					transport.WithReaderBufSize(common.IOBufSize),
				),
			),
		),
		sndr.WithClientsCount[T](clientsCount),
	)
}
