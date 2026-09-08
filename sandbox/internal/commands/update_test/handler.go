package update_test

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	updateTestsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/update_tests"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	update_error := updateTestsAction.UpdateTest(deps, entries.Path, entries.Name)

	if update_error != nil {
		deps.Std.Error("%s\n", update_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
