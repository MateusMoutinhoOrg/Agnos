package exec_test

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	execTestsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/exec_tests"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	exec_error := execTestsAction.ExecTest(sandbox, api.ExecTestProps{
		Path:   props.Path,
		Only:   entries.Only,
		Update: entries.Update,
	})

	if exec_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", exec_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
