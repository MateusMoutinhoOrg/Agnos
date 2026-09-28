package remove_table

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeTableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_table"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	remove_error := removeTableAction.RemoveTable(
		sandbox,
		props.Path,
		entries.Database,
		entries.Name,
	)

	if remove_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", remove_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
