package list_deps

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	listDepsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	deplist, list_error := listDepsAction.ListDeps(sandbox, props.Path)

	if list_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", list_error.Error())
	}

	for _, dep := range deplist {
		response.Printf("%-16s %-10s %-16s %s\n", dep.Name, installedMark(dep), adapters(sandbox, dep), dep.Help)
	}
	response.SetStatus(api.ExitOk)
	return nil
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
