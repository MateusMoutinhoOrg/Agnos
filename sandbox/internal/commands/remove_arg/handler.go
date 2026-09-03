package remove_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	removeArgAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/remove_arg"
)

func CommandHander(deps *deps.Deps, entries *Entries) int {
	remove_error := removeArgAction.RemoveArg(deps, entries.Path, entries.Command, entries.Name)
	if !entries.Quiet && remove_error != nil {
		deps.Std.Error(remove_error.Error())
	}
	if remove_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
