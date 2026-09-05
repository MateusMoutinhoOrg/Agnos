package add_cli_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AddCliExample scaffolds a new example under examples/cli/ — an example.sh
// that already runs — then runs build as a follow-up step so the example
// listing of the docs names it. Nothing under examples/ is compiled, so the
// build renders only.
func AddCliExample(deps *deps.Deps, path string, name string) error {
	io := smartio.New(deps, path, config.ProjectName)
	if err := AddCliExampleInternal(deps, io, name); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}
