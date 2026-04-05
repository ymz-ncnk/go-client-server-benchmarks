package main

import (
	"os"
	"reflect"

	"github.com/cmd-stream/cmd-stream-go/core"
	musgen "github.com/mus-format/mus-gen-go/mus"
	genops "github.com/mus-format/mus-gen-go/options/gen"
	introps "github.com/mus-format/mus-gen-go/options/interface"
	assert "github.com/ymz-ncnk/assert/panic"
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus/results"
)

func init() {
	assert.On = true
}

func main() {
	g, err := musgen.NewGenerator(
		genops.WithPkgPath("github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus/results"),
		genops.WithStream(),
	)
	assert.EqualError(err, nil)

	echoResultType := reflect.TypeFor[results.EchoResult]()
	err = g.AddStruct(echoResultType)
	assert.EqualError(err, nil)

	err = g.AddTyped(echoResultType)
	assert.EqualError(err, nil)

	err = g.AddInterface(reflect.TypeFor[core.Result](),
		introps.WithImpl(echoResultType),
		introps.WithMarshaller(),
	)
	assert.EqualError(err, nil)

	bs, err := g.Generate()
	assert.EqualError(err, nil)
	err = os.WriteFile("./mus.gen.go", bs, 0755)
	assert.EqualError(err, nil)
}
