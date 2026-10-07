package backoffice_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	backofficeInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/backoffice_init"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	init_error := backofficeInitAction.BackofficeInit(sandbox, props.Path)

	if init_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", init_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
