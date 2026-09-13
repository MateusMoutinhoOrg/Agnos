package remove_segment

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeSegmentAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_segment"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	remove_error := removeSegmentAction.RemoveSegment(sandbox, command.GetString("path"), command.GetString("route"), command.GetString("name"))

	if remove_error != nil {
		sandbox.Deps.Std.Error("%s\n", remove_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
