package list_adapters

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	listAdaptersAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_adapters"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	adapters, list_error := listAdaptersAction.ListAdapters(deps, entries.Path)

	if list_error != nil {
		deps.Std.Error("%s\n", list_error.Error())
		return api.ExitFailure
	}

	for _, adapter := range adapters {
		deps.Std.Printf("%-16s %-16s %-10s %-16s %s\n",
			adapter.Name, adapter.Dep, installedMark(adapter), boundBy(deps, adapter), adapter.Help)
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
func boundBy(deps *deps.Deps, adapter api.AdapterInfo) string {
	if len(adapter.Availables) == 0 {
		return "-"
	}
	return deps.Stringsdeps.Join(adapter.Availables, ",")
}
