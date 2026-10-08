package remove_database

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeDatabaseAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_database"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	remove_error := removeDatabaseAction.RemoveDatabase(sandbox, props.Path, entries.Name)

	if remove_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", remove_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
