package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	verifyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/verify"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	verify_error := verifyAction.Verify(sandbox, props.Path)

	// The schema check says the tree has the right shape; the runtime says
	// the Go toolchain accepts what is in it. `verify passed` means both.
	if verify_error == nil {
		verify_error = buildAction.RunRuntime(sandbox, props.Path, entries.Runtime)
	}

	if verify_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", verify_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
