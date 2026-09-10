package list_deps

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	listDepsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_deps"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	deplist, list_error := listDepsAction.ListDeps(deps, entries.Path)

	if list_error != nil {
		deps.Std.Error("%s\n", list_error.Error())
		return api.ExitFailure
	}

	for _, dep := range deplist {
		deps.Std.Printf("%s\n", dep)
	}
	return api.ExitOk
}
