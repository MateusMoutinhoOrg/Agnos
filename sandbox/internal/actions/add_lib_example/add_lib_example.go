package add_lib_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// AddLibExample scaffolds a new example under examples/lib/ — an example.go
// that already runs — then runs build as a follow-up step so the example
// listing of the docs names it. Nothing under examples/ is compiled, so the
// build renders only.
func AddLibExample(sandbox *api.Sandbox, props api.AddLibExampleProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := AddLibExampleInternal(sandbox, io, props.Name); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeNone})
}
