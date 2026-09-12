package remove_segment

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeSegmentAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_segment"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	remove_error := removeSegmentAction.RemoveSegment(sandbox, entries.Path, entries.Route, entries.Name)

	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
