package list_commands

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	listCommandsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_commands"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	lines, read_error := listCommandsAction.ListCommands(sandbox, props.Path)
	if read_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", read_error.Error())
	}

	for _, line := range lines {
		response.Printf("%s\n", line)
	}
	response.SetStatus(api.ExitOk)
	return nil
}
