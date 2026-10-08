package list_deps

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	listDepsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	catalog, list_error := listDepsAction.ListDeps(sandbox, api.ListDepsProps{Path: props.Path})

	if list_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", list_error.Error())
	}

	for _, dep := range catalog {
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
	return sandbox.Deps.StringsDeps.Join(dep.Adapters, ",")
}
