package rebalance_commands

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	rebalanceCommandsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/rebalance_commands"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	rebalance_error := rebalanceCommandsAction.RebalanceCommands(sandbox, api.RebalanceCommandsProps{
		Path: props.Path,
		Step: entries.Step,
	})
	if rebalance_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", rebalance_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
