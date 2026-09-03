package dep_list

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	depListAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/dep_list"
)

func CommandHander(deps *deps.Deps, entries *Entries) int {
	deplist, list_error := depListAction.DepList(deps, entries.Path)

	if !entries.Quiet && list_error != nil {
		deps.Std.Error(list_error.Error())
	}
	if list_error != nil {
		return api.ExitFailure
	}

	for _, dep := range deplist {
		deps.Std.Printf("%s\n", dep)
	}
	return api.ExitOk
}
