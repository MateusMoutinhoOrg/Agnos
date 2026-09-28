package add_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_command"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	add_error := addCommandAction.AddCommand(sandbox, props.Path, entries.Name, entries.Help, entries.Category)

	if add_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
