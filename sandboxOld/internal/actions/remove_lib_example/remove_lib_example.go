package remove_lib_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RemoveLibExample deletes one example of examples/lib/ whole, then runs build
// so the example listing of the docs is rewritten without it.
func RemoveLibExample(sandbox *api.Sandbox, path string, name string) error {
	io := smartio.New(sandbox, path, config.ProjectName)
	if err := RemoveLibExampleInternal(sandbox, io, name); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}
