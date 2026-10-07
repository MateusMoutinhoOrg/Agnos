package add_table

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addTableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_table"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	add_error := addTableAction.AddTable(
		sandbox,
		props.Path,
		entries.Database,
		entries.Name,
	)

	if add_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
