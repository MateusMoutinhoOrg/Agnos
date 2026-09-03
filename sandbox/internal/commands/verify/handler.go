package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	verifyAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/verify"
)

func CommandHander(deps *deps.Deps, entries *Entries) int {
	verify_error := verifyAction.Verify(deps, entries.Path)

	if !entries.Quiet && verify_error != nil {
		deps.Std.Error(verify_error.Error())
	}
	if verify_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
