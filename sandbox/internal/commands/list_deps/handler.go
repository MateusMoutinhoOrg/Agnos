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
		deps.Std.Printf("%-16s %-10s %-16s %s\n", dep.Name, installedMark(dep), adapters(deps, dep), dep.Help)
	}
	return api.ExitOk
}

// installedMark is the middle column: whether the project has this contract.
func installedMark(dep api.DepInfo) string {
	if dep.Installed {
		return "installed"
	}
	return "-"
}

// adapters is the third column: the installed adapters filling this dep, or
// the catalog's default one when nothing is installed yet.
func adapters(deps *deps.Deps, dep api.DepInfo) string {
	if len(dep.Adapters) == 0 {
		return dep.DefaultAdapter
	}
	return deps.Stringsdeps.Join(dep.Adapters, ",")
}
