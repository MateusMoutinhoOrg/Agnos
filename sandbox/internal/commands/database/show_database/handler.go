package show_database

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	showDatabaseAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/show_database"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	lines, show_error := showDatabaseAction.ShowDatabase(sandbox, api.ShowDatabaseProps{Path: props.Path, Name: input.Name})

	if show_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", show_error.Error())
	}

	for _, line := range lines {
		response.Printf("%s\n", line)
	}
	response.SetStatus(api.ExitOk)
	return nil
}
