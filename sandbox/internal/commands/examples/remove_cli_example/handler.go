package remove_cli_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeCliExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_cli_example"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	remove_error := removeCliExampleAction.RemoveCliExample(sandbox, api.RemoveCliExampleProps{Path: props.Path, Name: input.Name})
	if remove_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", remove_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
