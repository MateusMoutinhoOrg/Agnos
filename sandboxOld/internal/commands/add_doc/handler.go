package add_doc

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDocAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_doc"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	add_error := addDocAction.AddDoc(sandbox, api.DocProps{
		Path:        command.GetString("path"),
		Name:        command.GetString("name"),
		Description: command.GetString("description"),
		Themes:      command.GetStrings("theme"),
	})
	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
