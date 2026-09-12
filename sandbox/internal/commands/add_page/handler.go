package add_page

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addPageAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_page"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	add_error := addPageAction.AddPage(sandbox, api.PageProps{
		Path:    entries.Path,
		Name:    entries.Name,
		Trigger: entries.Trigger,
		Title:   entries.Title,
		Help:    entries.Help,
	})

	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
