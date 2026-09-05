package remove_cli_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RemoveCliExample deletes one example of examples/cli/ whole, then runs build
// so the example listing of the docs is rewritten without it.
func RemoveCliExample(deps *deps.Deps, path string, name string) error {
	io := smartio.New(deps, path, config.ProjectName)
	if err := RemoveCliExampleInternal(deps, io, name); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}
