package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// Build re-renders every generated file of the project at props.Path and then
// hands the result to props.Runtime, so a build only reports success when the
// toolchain accepts what was rendered.
func Build(sandbox *api.Sandbox, props api.BuildProps) error {
	io := smartio.New(sandbox, props.Path, config.ProjectName)
	err := BuildInternal(sandbox, io, props.Path)
	if err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return RunRuntime(sandbox, props.Path, props.Runtime)
}
