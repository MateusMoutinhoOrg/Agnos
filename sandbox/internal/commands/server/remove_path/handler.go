package remove_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removePathAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_path"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	err := removePathAction.RemovePath(sandbox, api.RemovePathProps{Path: props.Path, Route: input.Route, Name: input.Name})
	if err != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", err.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
