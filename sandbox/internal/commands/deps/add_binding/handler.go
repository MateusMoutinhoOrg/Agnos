package add_binding

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addBindingAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_binding"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	add_error := addBindingAction.AddBinding(sandbox, api.AddBindingProps{Path: props.Path, Binding: input.Name})

	if add_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
