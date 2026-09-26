package add_page

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addPageAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_page"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	add_error := addPageAction.AddPage(sandbox, api.PageProps{
		Path:  command.GetString("path"),
		Name:  command.GetString("name"),
		Title: command.GetString("title"),
	})

	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
