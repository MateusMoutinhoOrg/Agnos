package exec_test

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	execTestsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/exec_tests"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	exec_error := execTestsAction.ExecTest(sandbox, api.ExecTestProps{
		Path:   command.GetString("path"),
		Only:   command.GetString("only"),
		Update: command.GetBool("update"),
	})

	if exec_error != nil {
		sandbox.Deps.Std.Error("%s\n", exec_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
