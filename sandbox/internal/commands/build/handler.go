package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	verifyAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/verify"
)

func CommandHander(deps *deps.Deps, entries api.CliEntrys) int {
	quietFlag := entries.GetFlagById("quiet")
	unsafeFlag := entries.GetFlagById("unsafe")
	pathFlag := entries.GetFlagById("path")
	path := pathFlag.Values[0].String()

	if !unsafeFlag.Exist {
		if verify_error := verifyAction.Verify(deps, path); verify_error != nil {
			if !quietFlag.Exist {
				deps.Std.Error(verify_error.Error())
			}
			return api.ExitFailure
		}
	}

	build_error := buildAction.Build(deps, path)

	if !quietFlag.Exist && build_error != nil {
		deps.Std.Error(build_error.Error())
	}
	if build_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
