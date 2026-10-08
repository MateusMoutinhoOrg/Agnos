package list_adapters

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	listAdaptersAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_adapters"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	adapters, list_error := listAdaptersAction.ListAdapters(sandbox, api.ListAdaptersProps{Path: props.Path})

	if list_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", list_error.Error())
	}

	for _, adapter := range adapters {
		response.Printf("%-16s %-16s %-10s %-16s %s\n",
			adapter.Name, adapter.Dep, installedMark(adapter), boundBy(sandbox, adapter), adapter.Help)
	}
	response.SetStatus(api.ExitOk)
	return nil
}

// installedMark is the third column: whether the project has this package.
func installedMark(adapter api.AdapterInfo) string {
	if adapter.Installed {
		return "installed"
	}
	return "-"
}

// boundBy is the fourth column: the bindings that bind this adapter, which
// is what tells an installed adapter from a bound one.
func boundBy(sandbox *api.Sandbox, adapter api.AdapterInfo) string {
	if len(adapter.Bindings) == 0 {
		return "-"
	}
	return sandbox.Deps.StringsDeps.Join(adapter.Bindings, ",")
}
