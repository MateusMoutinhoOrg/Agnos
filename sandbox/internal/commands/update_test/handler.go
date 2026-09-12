package update_test

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	updateTestsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/update_tests"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	update_error := updateTestsAction.UpdateTest(sandbox, entries.Path, entries.Name)

	if update_error != nil {
		sandbox.Deps.Std.Error("%s\n", update_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
