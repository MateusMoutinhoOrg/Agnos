package remove_table

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeTableAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_table"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	remove_error := removeTableAction.RemoveTable(
		sandbox,
		props.Path,
		entries.Database,
		entries.Name,
	)

	if remove_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", remove_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
