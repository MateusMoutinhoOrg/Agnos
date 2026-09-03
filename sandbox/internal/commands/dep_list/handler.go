package dep_list

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	depListAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/dep_list"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	deplist, list_error := depListAction.DepList(deps, entries.Path)

	if list_error != nil {
		deps.Std.Error("%s\n", list_error.Error())
		return api.ExitFailure
	}

	for _, dep := range deplist {
		deps.Std.Printf("%s\n", dep)
	}
	return api.ExitOk
}
