package deps_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

func DepsPurge(sandbox *api.Sandbox, props api.DepsPurgeProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := DepsPurgeInternal(sandbox, io, props.Path); err != nil {
		return err
	}
	if err := buildAction.BuildInternal(sandbox, io, props.Path); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.RunRuntime(sandbox, props.Path, api.RuntimeNone)
}
