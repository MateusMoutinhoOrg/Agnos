package remove_parameter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeParameterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_parameter"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	err := removeParameterAction.RemoveParameter(sandbox, props.Path, entries.Route, entries.Name)
	if err != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", err.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
