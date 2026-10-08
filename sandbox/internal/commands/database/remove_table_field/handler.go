package remove_table_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeTableFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_table_field"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	remove_error := removeTableFieldAction.RemoveTableField(sandbox, api.DatabaseFieldProps{
		Path:     props.Path,
		Database: entries.Database,
		Table:    entries.Table,
		Parent:   entries.Parent,
		Name:     entries.Name,
	})

	if remove_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", remove_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
