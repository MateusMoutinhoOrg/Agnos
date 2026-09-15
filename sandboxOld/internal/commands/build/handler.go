package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	verifyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/verify"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	if !command.GetBool("unsafe") {
		// The schema gate only: the toolchain runs after the render, on
		// what was rendered, through the --runtime flag below.
		if verify_error := verifyAction.Verify(sandbox, command.GetString("path")); verify_error != nil {
			sandbox.Deps.Std.Error("%s\n", verify_error.Error())
			return api.ExitFailure
		}
	}

	build_error := buildAction.Build(sandbox, api.BuildProps{Path: command.GetString("path"), Runtime: command.GetString("runtime")})

	if build_error != nil {
		sandbox.Deps.Std.Error("%s\n", build_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
