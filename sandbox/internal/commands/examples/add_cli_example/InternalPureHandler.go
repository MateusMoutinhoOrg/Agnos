package add_cli_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addCliExampleAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_cli_example"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	add_error := addCliExampleAction.AddCliExample(sandbox, props.Path, entries.Name)
	if add_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
