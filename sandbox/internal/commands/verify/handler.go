package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	verifyAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/verify"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	verify_error := verifyAction.Verify(deps, entries.Path)

	// The schema check says the tree has the right shape; the runtime says
	// the Go toolchain accepts what is in it. `verify passed` means both.
	if verify_error == nil {
		verify_error = buildAction.RunRuntime(deps, entries.Path, entries.Runtime)
	}

	if verify_error != nil {
		deps.Std.Error("%s\n", verify_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
