package add_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_arg"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	add_error := addArgAction.AddArg(sandbox, api.FieldProps{
		Path:        command.GetString("path"),
		Command:     command.GetString("command"),
		Name:        command.GetString("name"),
		Identifiers: nil,
		Description: command.GetString("description"),
		Examples:    command.GetStrings("example"),
		Type:        command.GetString("type"),
		Default:     command.GetString("default"),
		Required:    command.GetBool("required"),
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
