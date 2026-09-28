package show_database

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	showDatabaseAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/show_database"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	lines, show_error := showDatabaseAction.ShowDatabase(sandbox, props.Path, entries.Database)

	if show_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", show_error.Error())
	}

	for _, line := range lines {
		response.Printf("%s\n", line)
	}
	response.SetStatus(api.ExitOk)
	return nil
}
