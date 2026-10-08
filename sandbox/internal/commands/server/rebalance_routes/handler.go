package rebalance_routes

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	rebalanceRoutesAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/rebalance_routes"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	rebalance_error := rebalanceRoutesAction.RebalanceRoutes(sandbox, api.RebalanceRoutesProps{
		Path: props.Path,
		Step: input.Step,
	})
	if rebalance_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", rebalance_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
