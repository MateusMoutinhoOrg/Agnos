package backoffice_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	backofficeInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/backoffice_init"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	init_error := backofficeInitAction.BackofficeInit(sandbox, api.BackofficeInitProps{Path: props.Path})

	if init_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", init_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
