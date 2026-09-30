package remove_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_arg"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	remove_error := removeArgAction.RemoveArg(sandbox, props.Path, entries.Target, entries.Name)
	if remove_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", remove_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
