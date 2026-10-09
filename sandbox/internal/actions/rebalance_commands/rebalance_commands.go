package rebalance_commands

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// RebalanceCommands lays the chain down again Step rungs apart, then runs build
// so every generated.new.go carries its new priority.
func RebalanceCommands(sandbox *api.Sandbox, props api.RebalanceCommandsProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := RebalanceCommandsInternal(sandbox, io, props); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
