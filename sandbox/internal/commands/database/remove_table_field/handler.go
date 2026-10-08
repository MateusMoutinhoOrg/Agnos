package remove_table_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeTableFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_table_field"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	remove_error := removeTableFieldAction.RemoveTableField(sandbox, api.RemoveTableFieldProps{
		Path:     props.Path,
		Database: input.Database,
		Table:    input.Table,
		Parent:   input.Parent,
		Name:     input.Name,
	})

	if remove_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", remove_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
