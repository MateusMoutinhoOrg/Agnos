package explain_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	explainCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/explain_command"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	lines, explain_error := explainCommandAction.ExplainCommand(sandbox, api.ExplainCommandProps{
		Path: props.Path,
		Argv: entries.Argv,
	})
	if explain_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", explain_error.Error())
	}

	for _, line := range lines {
		response.Printf("%s\n", line)
	}
	response.SetStatus(api.ExitOk)
	return nil
}
