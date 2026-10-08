package run_examples

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	runExamplesAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/run_examples"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	exec_error := runExamplesAction.RunExamples(sandbox, api.RunExamplesProps{
		Path:   props.Path,
		Only:   input.Only,
		Update: input.Update,
	})

	if exec_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", exec_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
