package add_table_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addTableFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_table_field"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	add_error := addTableFieldAction.AddTableField(sandbox, api.AddTableFieldProps{
		Path:     props.Path,
		Database: input.Database,
		Table:    input.Table,
		Parent:   input.Parent,
		Name:     input.Name,
		Type:     input.Type,
		Required: input.Required,
		Target:   input.Target,
	})

	if add_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
