package main

import (
	"os"
	"reflect"

	"github.com/cmd-stream/cmd-stream-go/core"
	musgen "github.com/mus-format/mus-gen-go/mus"
	genops "github.com/mus-format/mus-gen-go/options/gen"
	introps "github.com/mus-format/mus-gen-go/options/interface"
	assert "github.com/ymz-ncnk/assert/panic"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus/cmds"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus/receiver"
)

func init() {
	assert.On = true
}

func main() {
	g, err := musgen.NewGenerator(
		genops.WithPkgPath("github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus/cmds"),
		genops.WithStream(),
	)
	assert.EqualError(err, nil)

	echoCmdType := reflect.TypeFor[cmds.EchoCmd]()
	err = g.AddStruct(echoCmdType)
	assert.EqualError(err, nil)

	err = g.AddTyped(echoCmdType)
	assert.EqualError(err, nil)

	err = g.AddInterface(reflect.TypeFor[core.Cmd[receiver.Receiver]](),
		introps.WithImpl(echoCmdType),
		introps.WithMarshaller(),
	)
	assert.EqualError(err, nil)

	// Generate
	bs, err := g.Generate()
	assert.EqualError(err, nil)
	err = os.WriteFile("./mus.gen.go", bs, 0755)
	assert.EqualError(err, nil)
}
