package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	verifyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/verify"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	verify_error := verifyAction.Verify(sandbox, command.GetString("path"))

	// The schema check says the tree has the right shape; the runtime says
	// the Go toolchain accepts what is in it. `verify passed` means both.
	if verify_error == nil {
		verify_error = buildAction.RunRuntime(sandbox, command.GetString("path"), command.GetString("runtime"))
	}

	if verify_error != nil {
		sandbox.Deps.Std.Error("%s\n", verify_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
