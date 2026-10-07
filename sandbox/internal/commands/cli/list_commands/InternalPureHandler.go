package list_commands

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	listCommandsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_commands"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	lines, read_error := listCommandsAction.ListCommands(sandbox, props.Path)
	if read_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", read_error.Error())
	}

	for _, line := range lines {
		response.Printf("%s\n", line)
	}
	response.SetStatus(api.ExitOk)
	return nil
}
