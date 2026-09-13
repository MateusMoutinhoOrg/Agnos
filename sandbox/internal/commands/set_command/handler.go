package set_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_command"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	set_error := setCommandAction.SetCommand(sandbox, api.CommandProps{
		Path:            command.GetString("path"),
		Command:         command.GetString("name"),
		Help:            command.GetString("help"),
		Category:        command.GetString("category"),
		LongDescription: command.GetString("long-description"),
		Hidden:          command.GetBool("hidden"),
		Visible:         command.GetBool("visible"),
		Identifiers:     command.GetStrings("identifier"),
		Examples:        command.GetStrings("example"),
	})
	if set_error != nil {
		sandbox.Deps.Std.Error("%s\n", set_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
