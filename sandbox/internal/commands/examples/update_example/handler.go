package update_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	updateExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/update_example"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	update_error := updateExampleAction.UpdateExample(sandbox, api.UpdateExampleProps{Path: props.Path, Name: input.Name})

	if update_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", update_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
