package main

import (
	"os"
	"reflect"

	"github.com/cmd-stream/cmd-stream-go/core"
	musgen "github.com/mus-format/mus-gen-go/mus"
	genops "github.com/mus-format/mus-gen-go/options/gen"
	introps "github.com/mus-format/mus-gen-go/options/interface"
	assert "github.com/ymz-ncnk/assert/panic"
	tcpmus "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus"
)

func init() {
	assert.On = true
}

func main() {
	g, err := musgen.NewGenerator(
		genops.WithPkgPath("github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus"),
		genops.WithPackage("tcpmus"),
		genops.WithStream(),
	)
	assert.EqualError(err, nil)

	err = g.RegisterInterface(reflect.TypeFor[core.Cmd[tcpmus.Receiver]](),
		introps.WithStructImpl(reflect.TypeFor[tcpmus.EchoCmd]()),
		introps.WithRegisterMarshaller(),
	)
	assert.EqualError(err, nil)

	err = g.RegisterInterface(reflect.TypeFor[core.Result](),
		introps.WithStructImpl(reflect.TypeFor[tcpmus.EchoResult]()),
		introps.WithRegisterMarshaller(),
	)
	assert.EqualError(err, nil)

	bs, err := g.Generate()
	assert.EqualError(err, nil)
	err = os.WriteFile("./mus.gen.go", bs, 0755)
	assert.EqualError(err, nil)
}
