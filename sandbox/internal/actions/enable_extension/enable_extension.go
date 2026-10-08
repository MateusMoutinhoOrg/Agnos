package enable_extension

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

func EnableExtension(sandbox *api.Sandbox, props api.EnableExtensionProps) error {
	sandbox.Deps.StdDeps.Logf("enable-extension started with path %s extension %s \n", props.Path, props.Name)

	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := EnableExtensionInternal(sandbox, io, props.Name); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
