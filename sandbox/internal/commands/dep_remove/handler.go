package dep_remove

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	depRemoveAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/dep_remove"
)

func CommandHander(deps *deps.Deps, entries *Entries) int {
	remove_error := depRemoveAction.DepRemove(deps, entries.Path, entries.Dep)

	if !entries.Quiet && remove_error != nil {
		deps.Std.Error(remove_error.Error())
	}
	if remove_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
