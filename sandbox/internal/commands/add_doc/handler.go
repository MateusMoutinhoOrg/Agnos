package add_doc

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDocAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_doc"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	add_error := addDocAction.AddDoc(sandbox, api.DocProps{
		Path:        entries.Path,
		Name:        entries.Name,
		Description: entries.Description,
		Themes:      entries.Theme,
	})
	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
