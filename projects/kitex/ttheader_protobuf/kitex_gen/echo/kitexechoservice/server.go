// Code generated by Kitex v0.11.3. DO NOT EDIT.
package kitexechoservice

import (
	server "github.com/cloudwego/kitex/server"
	echo "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/ttheader_protobuf/kitex_gen/echo"
)

// NewServer creates a server.Server with the given handler and options.
func NewServer(handler echo.KitexEchoService, opts ...server.Option) server.Server {
	var options []server.Option

	options = append(options, opts...)

	svr := server.NewServer(options...)
	if err := svr.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	return svr
}

func RegisterService(svr server.Server, handler echo.KitexEchoService, opts ...server.RegisterOption) error {
	return svr.RegisterService(serviceInfo(), handler, opts...)
}
