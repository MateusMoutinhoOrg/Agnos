package rebalance_routes

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RebalanceRoutes lays the chain down again with room between its rungs, then
// runs build so every generated new.go carries its new priority.
func RebalanceRoutes(sandbox *api.Sandbox, props api.RebalanceRoutesProps) error {
	io := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := RebalanceRoutesInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
