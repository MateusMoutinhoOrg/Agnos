package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// Build re-renders every generated file of the project at props.Path and then
// hands the result to props.Runtime, so a build only reports success when the
// toolchain accepts what was rendered.
func Build(deps *deps.Deps, props api.BuildProps) error {
	io := smartio.New(deps, props.Path, config.ProjectName)
	err := BuildInternal(deps, io, props.Path)
	if err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return RunRuntime(deps, props.Path, props.Runtime)
}
