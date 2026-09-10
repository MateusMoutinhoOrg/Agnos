package add_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	install_error := addDepAction.AddDep(deps, api.AddDepProps{
		Path:    entries.Path,
		Dep:     entries.Dep,
		Adapter: entries.Adapter,
	})

	if install_error != nil {
		deps.Std.Error("%s\n", install_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
