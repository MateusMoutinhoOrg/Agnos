package add_command

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	addCommandAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/add_command"
)

func CommandHander(deps *deps.Deps, entries *Entries) int {
	add_error := addCommandAction.AddCommand(deps, entries.Path, entries.Name)

	if !entries.Quiet && add_error != nil {
		deps.Std.Error(add_error.Error())
	}
	if add_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
