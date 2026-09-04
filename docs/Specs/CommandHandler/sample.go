//go:build ignore

package dep_list

// This file is an illustrative sample, not part of the build.

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	depListAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/dep_list"
)

// CommandHandler backs the `dep-list` verb. Entries.Path is the --path flag,
// always populated by its default "."; Entries.Quiet has already silenced the
// progress channel by the time this runs.
func CommandHandler(deps *deps.Deps, entries *Entries) int {
	names, err := depListAction.DepList(deps, entries.Path)
	if err != nil {
		deps.Std.Error("%v\n", err) // failure: stderr, never silenced
		return api.ExitFailure
	}
	for _, name := range names {
		deps.Std.Printf("%s\n", name) // the result: stdout, never silenced
	}
	return api.ExitOk
}