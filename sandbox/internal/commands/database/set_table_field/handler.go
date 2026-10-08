package set_table_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setTableFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_table_field"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	set_error := setTableFieldAction.SetTableField(sandbox, api.SetTableFieldProps{
		Path:     props.Path,
		Database: input.Database,
		Table:    input.Table,
		Parent:   input.Parent,
		Name:     input.Name,
		Rename:   input.Rename,
		Type:     input.Type,
		Required: input.Required,
		Target:   input.Target,
		Clear:    input.Clear,
	})

	if set_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", set_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
