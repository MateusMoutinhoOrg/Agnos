package add_page

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	addPageAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_page"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	add_error := addPageAction.AddPage(deps, api.PageProps{
		Path:    entries.Path,
		Name:    entries.Name,
		Trigger: entries.Trigger,
		Title:   entries.Title,
		Help:    entries.Help,
	})

	if add_error != nil {
		deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
