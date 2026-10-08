package enable_extension

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	enableExtensionAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/enable_extension"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	enable_error := enableExtensionAction.EnableExtension(sandbox, api.EnableExtensionProps{Path: props.Path, Name: input.Name})

	if enable_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", enable_error.Error())
	}

	response.Printf("%s is on\n", input.Name)
	response.SetStatus(api.ExitOk)
	return nil
}
