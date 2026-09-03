package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	verifyAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/verify"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	if !entries.Unsafe {
		// The schema gate only: the toolchain runs after the render, on
		// what was rendered, through entries.Runtime below.
		if verify_error := verifyAction.Verify(deps, entries.Path); verify_error != nil {
			deps.Std.Error("%s\n", verify_error.Error())
			return api.ExitFailure
		}
	}

	build_error := buildAction.Build(deps, api.BuildProps{Path: entries.Path, Runtime: entries.Runtime})

	if build_error != nil {
		deps.Std.Error("%s\n", build_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
