package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	verifyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/verify"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	// The runtime is checked before verify walks the whole project: a typo in
	// it is a usage error, not something to find out at the end.
	runtime := sandbox.Deps.StringsDeps.ToLower(sandbox.Deps.StringsDeps.TrimSpace(input.Runtime))
	if runtime != api.RuntimeGo && runtime != api.RuntimeNone {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "runtime",
			sandbox.Deps.StdDeps.Sprintf("unknown runtime %q (use %q or %q)", input.Runtime, api.RuntimeGo, api.RuntimeNone))
	}

	if !input.Unsafe {
		// The schema gate only: the toolchain runs after the render, on
		// what was rendered, through the --runtime flag below.
		if verify_error := verifyAction.Verify(sandbox, api.VerifyProps{Path: props.Path}); verify_error != nil {
			return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", verify_error.Error())
		}
	}

	build_error := buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: runtime})

	if build_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", build_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
