package set_segment

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setSegmentAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_segment"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	set_error := setSegmentAction.SetSegment(sandbox, api.RouteFieldEditProps{
		Path:        command.GetString("path"),
		Route:       command.GetString("route"),
		Name:        command.GetString("name"),
		Rename:      command.GetString("rename"),
		Identifier:  command.GetString("identifier"),
		Description: command.GetString("description"),
		Examples:    command.GetStrings("example"),
		Type:        command.GetString("type"),
		Required:    command.GetBool("required"),
		Array:       command.GetBool("array"),
		Min:         command.GetString("min"),
		Max:         command.GetString("max"),
		Clear:       command.GetStrings("clear"),
	})
	if set_error != nil {
		sandbox.Deps.Std.Error("%s\n", set_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
