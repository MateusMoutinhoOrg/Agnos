package add_segment

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addSegmentAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_segment"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	add_error := addSegmentAction.AddSegment(sandbox, api.RouteFieldProps{
		Path:        command.GetString("path"),
		Route:       command.GetString("route"),
		Name:        command.GetString("name"),
		Identifier:  command.GetString("identifier"),
		Description: command.GetString("description"),
		Examples:    command.GetStrings("example"),
		Type:        command.GetString("type"),
		Array:       command.GetBool("array"),
		Min:         command.GetString("min"),
		Max:         command.GetString("max"),
		Position:    command.GetInt("position"),
	})
	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
