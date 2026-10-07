package remove_lib_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeLibExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_lib_example"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	remove_error := removeLibExampleAction.RemoveLibExample(sandbox, props.Path, entries.Name)
	if remove_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", remove_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
