package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	verifyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/verify"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	if !entries.Unsafe {
		// The schema gate only: the toolchain runs after the render, on
		// what was rendered, through the --runtime flag below.
		if verify_error := verifyAction.Verify(sandbox, props.Path); verify_error != nil {
			return cliio.Fail(sandbox, api.ExitFailure, "", verify_error.Error())
		}
	}

	build_error := buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: entries.Runtime})

	if build_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", build_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
