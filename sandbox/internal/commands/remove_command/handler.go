package remove_command

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	removeCommandAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/remove_command"
)

func CommandHander(deps *deps.Deps, entries *Entries) int {
	remove_error := removeCommandAction.RemoveCommand(deps, entries.Path, entries.Name)
	if !entries.Quiet && remove_error != nil {
		deps.Std.Error(remove_error.Error())
	}
	if remove_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
