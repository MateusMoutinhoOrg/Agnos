package remove_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_command"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	remove_error := removeCommandAction.RemoveCommand(sandbox, props.Path, entries.Name)
	if remove_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", remove_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
