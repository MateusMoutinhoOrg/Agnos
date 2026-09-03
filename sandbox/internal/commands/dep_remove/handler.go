package dep_remove

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	depRemoveAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/dep_remove"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	remove_error := depRemoveAction.DepRemove(deps, entries.Path, entries.Dep)

	if remove_error != nil {
		deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
