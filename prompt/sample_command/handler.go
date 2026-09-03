package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	verifyAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/verify"
)

func CommandHander(deps *deps.Deps, entries *Entries) int {
	if !entries.Unsafe {
		if verify_error := verifyAction.Verify(deps, entries.Path); verify_error != nil {
			if !entries.Quiet {
				deps.Std.Error(verify_error.Error())
			}
			return api.ExitFailure
		}
	}

	build_error := buildAction.Build(deps, entries.Path)

	if !entries.Quiet && build_error != nil {
		deps.Std.Error(build_error.Error())
	}
	if build_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
