package rename_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RenameRoute moves one route package to a new name, then runs build so the
// generated new.go, entries.go and the server's list of routes follow it.
func RenameRoute(sandbox *api.Sandbox, props api.RenameRouteProps) error {
	io := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := RenameRouteInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
