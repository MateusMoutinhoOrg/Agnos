package explain_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	explainCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/explain_command"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	lines, explain_error := explainCommandAction.ExplainCommand(sandbox, api.ExplainCommandProps{
		Path: props.Path,
		Argv: input.Argv,
	})
	if explain_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", explain_error.Error())
	}

	for _, line := range lines {
		response.Printf("%s\n", line)
	}
	response.SetStatus(api.ExitOk)
	return nil
}
