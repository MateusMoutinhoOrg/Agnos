package list_adapters

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	listAdaptersAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_adapters"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	adapters, list_error := listAdaptersAction.ListAdapters(sandbox, entries.Path)

	if list_error != nil {
		sandbox.Deps.Std.Error("%s\n", list_error.Error())
		return api.ExitFailure
	}

	for _, adapter := range adapters {
		sandbox.Deps.Std.Printf("%-16s %-16s %-10s %-16s %s\n",
			adapter.Name, adapter.Dep, installedMark(adapter), boundBy(sandbox, adapter), adapter.Help)
	}
	return api.ExitOk
}

// installedMark is the third column: whether the project has this package.
func installedMark(adapter api.AdapterInfo) string {
	if adapter.Installed {
		return "installed"
	}
	return "-"
}

// boundBy is the fourth column: the availables that bind this adapter, which
// is what tells an installed adapter from a bound one.
func boundBy(sandbox *api.Sandbox, adapter api.AdapterInfo) string {
	if len(adapter.Availables) == 0 {
		return "-"
	}
	return sandbox.Deps.Stringsdeps.Join(adapter.Availables, ",")
}
