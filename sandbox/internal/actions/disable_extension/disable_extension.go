package disable_extension

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

func DisableExtension(sandbox *api.Sandbox, props api.DisableExtensionProps) error {
	sandbox.Deps.StdDeps.Logf("disable-extension started with path %s extension %s \n", props.Path, props.Name)

	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := DisableExtensionInternal(sandbox, io, props.Name); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeNone})
}
