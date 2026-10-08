package show_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	showCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/show_command"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	lines, read_error := showCommandAction.ShowCommand(sandbox, props.Path, entries.Target)
	if read_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", read_error.Error())
	}

	for _, line := range lines {
		response.Printf("%s\n", line)
	}
	response.SetStatus(api.ExitOk)
	return nil
}
