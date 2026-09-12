package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	verifyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/verify"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	verify_error := verifyAction.Verify(sandbox, entries.Path)

	// The schema check says the tree has the right shape; the runtime says
	// the Go toolchain accepts what is in it. `verify passed` means both.
	if verify_error == nil {
		verify_error = buildAction.RunRuntime(sandbox, entries.Path, entries.Runtime)
	}

	if verify_error != nil {
		sandbox.Deps.Std.Error("%s\n", verify_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
