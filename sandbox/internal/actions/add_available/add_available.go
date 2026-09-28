package add_available

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

func AddAvailable(sandbox *api.Sandbox, path string, available string) error {
	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	if err := AddAvailableInternal(sandbox, io, path, available); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: path, Runtime: api.RuntimeGo})
}
