package exec_test

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	execTestsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/exec_tests"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	exec_error := execTestsAction.ExecTest(sandbox, api.ExecTestProps{
		Path:   entries.Path,
		Only:   entries.Only,
		Update: entries.Update,
	})

	if exec_error != nil {
		sandbox.Deps.Std.Error("%s\n", exec_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
