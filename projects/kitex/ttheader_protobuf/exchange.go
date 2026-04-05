package ttheaderproto

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/ttheader_protobuf/kitex_gen/echo"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/ttheader_protobuf/kitex_gen/echo/kitexechoservice"
)

func ExchangeFixed(data *echo.KitexData, client kitexechoservice.Client,
	copsD chan<- time.Duration,
	wg *sync.WaitGroup,
	b *testing.B,
) {
	defer wg.Done()
	start := time.Now()
	r, err := exchange(data, client)
	if err != nil {
		b.Error(err)
		return
	}
	common.QueueCopD(copsD, time.Since(start))
	if !EqualData(data, r) {
		b.Error("unexpected result")
	}
}

func ExchangeQPS(data *echo.KitexData, client kitexechoservice.Client,
	wg *sync.WaitGroup,
	b *testing.B,
) {
	defer wg.Done()
	r, err := exchange(data, client)
	if err != nil {
		b.Error(err)
		return
	}
	if !EqualData(data, r) {
		b.Error("unexpected result")
	}
}

func exchange(data *echo.KitexData, client kitexechoservice.Client,
) (r *echo.KitexData, err error) {
	return client.Echo(context.Background(), data)
}
