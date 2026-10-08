package disable_extension

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	disableExtensionAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/disable_extension"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	disable_error := disableExtensionAction.DisableExtension(sandbox, api.DisableExtensionProps{Path: props.Path, Name: input.Name})

	if disable_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", disable_error.Error())
	}

	response.Printf("%s is off; what it wrote is yours now\n", input.Name)
	response.SetStatus(api.ExitOk)
	return nil
}
