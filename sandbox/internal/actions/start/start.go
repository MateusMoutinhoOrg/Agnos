package start

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

func Start(sandbox *api.Sandbox, props api.StartProps) error {
	io := smartio.New(sandbox, props.Path, props.ProjectName)
	if err := StartInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := buildAction.BuildInternal(sandbox, io, props.Path); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.RunRuntime(sandbox, props.Path, api.RuntimeGo)
}
