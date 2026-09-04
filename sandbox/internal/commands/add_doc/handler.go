package add_doc

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	addDocAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_doc"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	add_error := addDocAction.AddDoc(deps, api.DocProps{
		Path:        entries.Path,
		Name:        entries.Name,
		Description: entries.Description,
		Themes:      entries.Theme,
	})
	if add_error != nil {
		deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
