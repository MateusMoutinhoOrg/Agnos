package update_test

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	updateTestsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/update_tests"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	update_error := updateTestsAction.UpdateTest(sandbox, command.GetString("path"), command.GetString("name"))

	if update_error != nil {
		sandbox.Deps.Std.Error("%s\n", update_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
