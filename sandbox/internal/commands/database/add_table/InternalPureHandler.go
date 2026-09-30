package add_table

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addTableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_table"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	add_error := addTableAction.AddTable(
		sandbox,
		props.Path,
		entries.Database,
		entries.Name,
	)

	if add_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
