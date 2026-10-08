package remove_parameter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeParameterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_parameter"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	err := removeParameterAction.RemoveParameter(sandbox, api.RemoveParameterProps{Path: props.Path, Route: input.Route, Name: input.Name})
	if err != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", err.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
