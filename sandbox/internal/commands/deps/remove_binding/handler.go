package remove_binding

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeBindingAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_binding"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	remove_error := removeBindingAction.RemoveBinding(sandbox, api.RemoveBindingProps{Path: props.Path, Binding: input.Name})

	if remove_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", remove_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
