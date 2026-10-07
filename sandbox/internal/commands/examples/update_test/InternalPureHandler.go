package update_test

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	updateTestsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/update_tests"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	update_error := updateTestsAction.UpdateTest(sandbox, props.Path, entries.Name)

	if update_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", update_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
