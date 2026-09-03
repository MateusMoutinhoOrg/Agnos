package remove_flag

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	removeFlagAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/remove_flag"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	remove_error := removeFlagAction.RemoveFlag(deps, entries.Path, entries.Command, entries.Name)
	if remove_error != nil {
		deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
