package remove_segment

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	removeSegmentAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_segment"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	remove_error := removeSegmentAction.RemoveSegment(deps, entries.Path, entries.Route, entries.Name)

	if remove_error != nil {
		deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
