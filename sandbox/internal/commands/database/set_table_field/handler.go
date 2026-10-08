package set_table_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setTableFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_table_field"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	set_error := setTableFieldAction.SetTableField(sandbox, api.DatabaseFieldEditProps{
		Path:     props.Path,
		Database: entries.Database,
		Table:    entries.Table,
		Parent:   entries.Parent,
		Name:     entries.Name,
		Rename:   entries.Rename,
		Type:     entries.Type,
		Required: entries.Required,
		Target:   entries.Target,
		Clear:    entries.Clear,
	})

	if set_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", set_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
