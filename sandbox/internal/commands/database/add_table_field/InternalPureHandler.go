package add_table_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addTableFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_table_field"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	add_error := addTableFieldAction.AddTableField(sandbox, api.DatabaseFieldProps{
		Path:     props.Path,
		Database: entries.Database,
		Table:    entries.Table,
		Parent:   entries.Parent,
		Name:     entries.Name,
		Type:     entries.Type,
		Required: entries.Required,
		Target:   entries.Target,
	})

	if add_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
