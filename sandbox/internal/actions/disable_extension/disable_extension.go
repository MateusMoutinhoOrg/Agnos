package disable_extension

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

func DisableExtension(sandbox *api.Sandbox, path string, name string) error {
	sandbox.Deps.Std.Log("disable-extension started with path %s extension %s \n", path, name)

	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	if err := DisableExtensionInternal(sandbox, io, name); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}
