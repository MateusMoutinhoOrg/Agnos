package rebalance_routes

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	rebalanceRoutesAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/rebalance_routes"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	rebalance_error := rebalanceRoutesAction.RebalanceRoutes(sandbox, api.RebalanceRoutesProps{
		Path: command.GetString("path"),
		Step: command.GetInt("step"),
	})
	if rebalance_error != nil {
		sandbox.Deps.Std.Error("%s\n", rebalance_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
