package list_deps

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	listDepsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_deps"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	deplist, list_error := listDepsAction.ListDeps(sandbox, entries.Path)

	if list_error != nil {
		sandbox.Deps.Std.Error("%s\n", list_error.Error())
		return api.ExitFailure
	}

	for _, dep := range deplist {
		sandbox.Deps.Std.Printf("%-16s %-10s %-16s %s\n", dep.Name, installedMark(dep), adapters(sandbox, dep), dep.Help)
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
func adapters(sandbox *api.Sandbox, dep api.DepInfo) string {
	if len(dep.Adapters) == 0 {
		return dep.DefaultAdapter
	}
	return sandbox.Deps.Stringsdeps.Join(dep.Adapters, ",")
}
