package exec_test

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	execTestsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/exec_tests"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	exec_error := execTestsAction.ExecTest(deps, api.ExecTestProps{
		Path:   entries.Path,
		Only:   entries.Only,
		Update: entries.Update,
	})

	if exec_error != nil {
		deps.Std.Error("%s\n", exec_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
