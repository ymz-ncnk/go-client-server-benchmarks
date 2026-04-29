package grpcproto

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/grpc_protobuf/kitex_gen/echo"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/grpc_protobuf/kitex_gen/echo/kitexechoservice"
)

func ExchangeQPS(data *echo.KitexData, client kitexechoservice.Client,
	wg *sync.WaitGroup, b *testing.B,
) {
	defer wg.Done()
	_, err := client.Echo(context.Background(), data)
	if err != nil {
		b.Error(err)
		return
	}
}

func ExchangeFixed(data *echo.KitexData, client kitexechoservice.Client,
	copsD chan time.Duration, wg *sync.WaitGroup, b *testing.B,
) {
	defer wg.Done()
	start := time.Now()
	_, err := client.Echo(context.Background(), data)
	if err != nil {
		b.Error(err)
		return
	}
	copsD <- time.Since(start)
}
